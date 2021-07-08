package rmq

import (
	"github.com/streadway/amqp"
)

type MessageBroker interface {
	Produce(x string, key string, body []byte) error
	Consume(qname string) (<-chan amqp.Delivery, error)
	DeclareQueue(name string) (amqp.Queue, error)
	DeclareExchange(name string, xType string) error
}

type Client struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger Logger
}

func NewClient(cfg Config, logger Logger) (MessageBroker, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		logger.WithErr(err).Fatal("failed to connect to RMQ")
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.WithErr(err).Fatal("failed to open a channel")
		return nil, err
	}

	return &Client{
		conn:   conn,
		ch:     ch,
		logger: logger,
	}, nil
}

func (c *Client) Produce(x string, key string, body []byte) error {
	if err := c.ch.Publish(
		x,     // exchange
		key,   // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	); err != nil {
		c.logger.WithErr(err).Error("failed to publish a message")
		return err
	}

	return nil
}
func (c *Client) Consume(qname string) (<-chan amqp.Delivery, error) {
	msgs, err := c.ch.Consume(
		qname, // queue
		"",    // consumer
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)
	if err != nil {
		c.logger.WithErr(err).Error("failed to subsribe")
		return nil, err
	}

	return msgs, err
}
func (c *Client) Close() error {
	return c.ch.Close()
}

func (c *Client) DeclareQueue(name string) (amqp.Queue, error) {
	q, err := c.ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		c.logger.WithErr(err).Error("failed to declare exchange")
		return q, err
	}

	return q, nil
}

func (c *Client) DeclareExchange(name string, xType string) error {
	if err := c.ch.ExchangeDeclare(
		name,  // name
		xType, // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	); err != nil {
		c.logger.WithErr(err).Error("failed to declare exchange")
		return err
	}

	return nil
}
