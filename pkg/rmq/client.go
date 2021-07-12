package rmq

import (
	"github.com/streadway/amqp"
)

const (
	defaultContentType = "application/json"
)

// client wraps connection and channel of amqp.
type client struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger Logger
}

// newClient creates new instance of client with its own connection and channel.
func newClient(url string, logger Logger) (*client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Fatal("failed to connect to RMQ")
		return nil, err
	}

	// TODO: use 1 channel per request.
	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal("failed to open a channel")
		return nil, err
	}

	return &client{
		conn:   conn,
		ch:     ch,
		logger: logger,
	}, nil
}

func (c *client) declareQueue(name string, durable bool, exclusive bool) (amqp.Queue, error) {
	q, err := c.ch.QueueDeclare(
		name,
		durable,
		false, // delete when unused
		exclusive,
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		c.logger.Error("failed to declare queue")
		return q, err
	}

	return q, nil
}

func (c *client) bindQueue(meta *Meta) error {
	if err := c.ch.QueueBind(
		meta.QName,
		meta.Key, // routing key
		meta.XName,
		false,
		nil,
	); err != nil {
		c.logger.Error("failed to bind queue")
		return err
	}

	return nil
}

func (c *client) declareExchange(name string, xtype string, durable bool) error {
	if err := c.ch.ExchangeDeclare(
		name,
		xtype,
		durable,
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	); err != nil {
		c.logger.Error("failed to declare exchange")
		return err
	}

	return nil
}

func (c *client) setQos(prefetchCount int, prefetchSize int) error {
	if err := c.ch.Qos(
		prefetchCount,
		prefetchSize,
		false, // global
	); err != nil {
		c.logger.Error("failed to set Qos")
		return err
	}

	return nil
}
