package rmq

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"runtime/debug"
	"time"

	amqpv1 "github.com/gxravel/bus-routes/pkg/rmq/v1"

	"github.com/streadway/amqp"
)

const (
	defaultTimeout = time.Second * 5
)

// internalErrBody is amqpv1.Response of internal server error to use in recover block.
var internalErrBody []byte

func init() {
	var jerr error
	internalErrBody, jerr = json.Marshal(&amqpv1.Response{
		Error: &amqpv1.APIError{
			Code: http.StatusInternalServerError,
			Reason: &amqpv1.APIReason{
				Err: "Internal server error",
			},
		},
	})
	if jerr != nil {
		panic("failed to marshal internal server error body")
	}
}

// handlerFunc is function signature for an amqp handler.
type handlerFunc func(context.Context, *amqp.Delivery) (interface{}, error)

// wrapHandler wraps handler with meta, message and processing message function.
func (p *Publisher) wrapHandler(
	meta *Meta,
	delivery <-chan amqp.Delivery,
	handler handlerFunc,
) func(context.Context) {
	return func(ctx context.Context) {
		for {
			select {
			case message := <-delivery:
				go p.processMessage(ctx, meta, &message, handler)

			case <-ctx.Done():
				return
			}
		}
	}
}

// handlerResult describes result of an amqp handler.
type handlerResult struct {
	Data interface{}
	Err  error
}

// processMessage processes incoming message, adding recoverer and timeout,
// and listening for result to produce data or error.
func (p *Publisher) processMessage(ctx context.Context, meta *Meta, message *amqp.Delivery, handler handlerFunc) {
	p.UseFreeChannel()
	defer p.FreeChannel()

	defer p.recover(meta)

	meta.CorrID = message.CorrelationId
	if message.ReplyTo != "" {
		meta.Key = message.ReplyTo
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	ch := make(chan *handlerResult)

	go func() {
		data, err := handler(ctx, message)
		ch <- &handlerResult{Data: data, Err: err}
	}()

	select {
	case <-ctx.Done():
		p.ProduceError(ctx, meta, errors.New("request canceled"))
		return

	case result := <-ch:
		if result.Err != nil {
			p.ProduceError(ctx, meta, result.Err)
			return
		}
		p.ProduceData(ctx, meta, result.Data)
	}
}

func (p *Publisher) recover(meta *Meta) {
	if err := recover(); err != nil {
		p.logger.Errorf("panic: %v", err)
		debug.PrintStack()

		if perr := p.produce(meta, internalErrBody); perr != nil {
			p.logger.Error("failed to produce internal error body while recovering")
		}
	}
}
