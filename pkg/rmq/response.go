package rmq

import (
	"context"
	"encoding/json"

	amqpv1 "github.com/gxravel/bus-routes/pkg/rmq/v1"
)

// produceJSON produces a message in format of JSON.
func (c *Client) produceJSON(ctx context.Context, meta *Meta, data interface{}) {
	c.logger.Debugf("sending data: %v", data)
	body, err := c.ConvertToMessage(data)
	if err != nil {
		return
	}

	if err := c.produce(meta, body); err != nil {
		c.logger.Error("failed to publish message")
	}
}

// ProduceError resolves status code and produces an APIError.
func (c *Client) ProduceError(ctx context.Context, meta *Meta, err error) {
	reason := ConvertToReason(err)
	code := ResolveStatusCode(reason.Err)

	c.produceJSON(ctx, meta, &amqpv1.Response{
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
func (c *Client) ProduceData(ctx context.Context, meta *Meta, data interface{}) {
	c.produceJSON(ctx, meta, &amqpv1.Response{
		Data: data,
	})
}

// TranslateMessage translate a message body to JSON.
func (c *Client) TranslateMessage(message []byte, data interface{}) error {
	if err := json.Unmarshal(message, data); err != nil {
		c.logger.Error("failed to translate message")
		return err
	}

	return nil
}

// ConvertToMessage converts a JSON-encoded data to message body.
func (c *Client) ConvertToMessage(data interface{}) ([]byte, error) {
	message, err := json.Marshal(data)
	if err != nil {
		c.logger.Error("failed to convert to message")
		return nil, err
	}

	return message, nil
}
