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
	shouldAcknowledge bool,
	handler handlerFunc,
) func(context.Context) {
	return func(ctx context.Context) {
		for {
			select {
			case message := <-delivery:
				go p.processMessage(ctx, meta, &message, shouldAcknowledge, handler)

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
// and listening for result to produce data or an error.
func (p *Publisher) processMessage(ctx context.Context, meta *Meta, message *amqp.Delivery, shouldAcknowledge bool, handler handlerFunc) {
	pub, free := p.WithFreeChannel()
	defer free()

	defer pub.recover(meta)

	meta.CorrID = message.CorrelationId
	if message.ReplyTo != "" {
		meta.Key = message.ReplyTo
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	rch := make(chan *handlerResult)

	go func() {
		data, err := handler(ctx, message)
		rch <- &handlerResult{Data: data, Err: err}
	}()

	select {
	case <-ctx.Done():
		pub.ProduceError(ctx, meta, errors.New("request canceled"))
		return

	case result := <-rch:
		if result.Err != nil {
			pub.ProduceError(ctx, meta, result.Err)
		} else {
			pub.ProduceData(ctx, meta, result.Data)
		}

		if shouldAcknowledge {
			if err := message.Ack(false); err != nil {
				pub.logger.Fatal("could not acknowledge a delivery")
			}
		}
	}
}

func (p *Publisher) recover(meta *Meta) {
	if err := recover(); err != nil {
		p.logger.Errorf("panic: %v", err)
		debug.PrintStack()

		if perr := p.Produce(meta, internalErrBody); perr != nil {
			p.logger.Error("failed to produce internal error body while recovering")
		}
	}
}
