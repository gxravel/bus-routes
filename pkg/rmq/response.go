package rmq

import (
	"context"
	"encoding/json"

	amqpv1 "github.com/gxravel/bus-routes/pkg/rmq/v1"

	"github.com/pkg/errors"
)

// produceJSON produces a message in format of JSON.
func (p *Publisher) produceJSON(ctx context.Context, meta *Meta, data interface{}) {
	p.logger.Debugf("send data: %v", data)

	body, err := ConvertToMessage(data)
	if err != nil {
		return
	}

	if err := p.Produce(meta, body); err != nil {
		p.logger.Fatalf(err.Error(), "failed to publish message")
	}
}

// ProduceError resolves status code and produces an APIError.
func (p *Publisher) ProduceError(ctx context.Context, meta *Meta, err error) {
	reason := ConvertToReason(err)
	code := ResolveStatusCode(reason.Err)

	p.produceJSON(ctx, meta, &amqpv1.Response{
		Error: &amqpv1.APIError{
			Code: code,
			Reason: &amqpv1.APIReason{
				Err:     reason.Error(),
				Message: reason.Message,
			},
		},
	})
}

// ProduceData produces a message with the amqpv1.Response.
func (p *Publisher) ProduceData(ctx context.Context, meta *Meta, data interface{}) {
	p.produceJSON(ctx, meta, &amqpv1.Response{
		Data: data,
	})
}

// TranslateMessage translate a message body to JSON.
func TranslateMessage(message []byte, data interface{}) error {
	if err := json.Unmarshal(message, data); err != nil {
		return errors.Wrap(err, "failed to translate message")
	}

	return nil
}

// ConvertToMessage converts a JSON-encoded data to message body.
func ConvertToMessage(data interface{}) ([]byte, error) {
	message, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert to message")
	}

	return message, nil
}
