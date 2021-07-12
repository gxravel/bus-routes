package rmq

import (
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

const (
	defaultContentType = "application/json"
)

// Client wraps connection and channel of amqp.
type Client struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger Logger
}

// Client creates new instance of client with its own connection and channel.
func NewClient(cfg Config, logger Logger) (*Client, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		logger.Fatal("failed to connect to RMQ")
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal("failed to open a channel")
		return nil, err
	}

	return &Client{
		conn:   conn,
		ch:     ch,
		logger: logger,
	}, nil
}

// NewTask produces message to a named queue.
func (c *Client) NewTask(qname string, body []byte) error {
	q, err := c.declareQueue(qname, false, false)
	if err != nil {
		return err
	}

	return c.produce(&Meta{QName: q.Name, Mode: amqp.Persistent}, body)
}

// WorkOnTask consumes message from a named queue and defined QoS.
func (c *Client) WorkOnTask(qname string, prefetchCount int) (<-chan amqp.Delivery, error) {
	q, err := c.declareQueue(qname, false, false)
	if err != nil {
		return nil, err
	}

	if err := c.setQos(prefetchCount, 0); err != nil {
		return nil, err
	}

	return c.consume(q.Name, true, false)
}

// Publish publishes a message to the exchange.
func (c *Client) Publish(meta *Meta, body []byte) error {
	if err := c.declareExchange(meta.XName, meta.XType, true); err != nil {
		return err
	}

	return c.produce(meta, body)
}

// Subscribe subscribes for a messages from the exchange.
func (c *Client) Subscribe(meta *Meta) (<-chan amqp.Delivery, error) {
	if err := c.declareExchange(meta.XName, meta.XType, true); err != nil {
		return nil, err
	}

	q, err := c.declareQueue(meta.QName, false, true)
	if err != nil {
		return nil, err
	}

	if err := c.bindQueue(meta); err != nil {
		return nil, err
	}

	return c.consume(q.Name, true, false)
}

// CallRPC produces a message to a nammed queue, pre-installing ReplyTo property, which it consumes from.
func (c *Client) CallRPC(qname string, body []byte) (<-chan amqp.Delivery, string, error) {
	q, err := c.declareQueue("", false, true)
	if err != nil {
		return nil, "", err
	}

	delivery, err := c.consume(q.Name, true, false)
	if err != nil {
		return nil, "", err
	}

	meta := &Meta{
		Key:    qname,
		QName:  q.Name,
		CorrID: uuid.New().String(),
	}

	return delivery, meta.CorrID, c.produce(meta, body)
}

// Close closes connection of amqp client.
func (c *Client) Close() error {
	return c.ch.Close()
}

func (c *Client) produce(meta *Meta, body []byte) error {
	pub := amqp.Publishing{
		DeliveryMode:  meta.Mode,
		ContentType:   defaultContentType,
		CorrelationId: meta.CorrID,
		ReplyTo:       meta.QName,
		Body:          body,
	}

	if err := c.ch.Publish(
		meta.XName, // exchange
		meta.Key,   // routing key
		false,      // mandatory
		false,      // immediate
		pub,
	); err != nil {
		c.logger.Error("failed to produce a message")
		return err
	}

	return nil
}

func (c *Client) consume(qname string, autoAck bool, exclusive bool) (<-chan amqp.Delivery, error) {
	delivery, err := c.ch.Consume(
		qname,     // queue
		"",        // consumer
		autoAck,   // auto ack
		exclusive, // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		c.logger.Error("failed to consume")
		return nil, err
	}

	return delivery, err
}

func (c *Client) declareQueue(name string, durable bool, exclusive bool) (amqp.Queue, error) {
	q, err := c.ch.QueueDeclare(
		name,
		durable,
		false, // delete when unused
		exclusive,
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		c.logger.Error("failed to declare exchange")
		return q, err
	}

	return q, nil
}

func (c *Client) bindQueue(meta *Meta) error {
	if err := c.ch.QueueBind(
		meta.QName,
		meta.Key, // routing key
		meta.XName,
		false,
		nil,
	); err != nil {
		c.logger.Error("failed to declare exchange")
		return err
	}

	return nil
}

func (c *Client) declareExchange(name string, xtype string, durable bool) error {
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

func (c *Client) setQos(prefetchCount int, prefetchSize int) error {
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
