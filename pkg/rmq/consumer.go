package rmq

import "github.com/streadway/amqp"

// Consumer implements the consumer methods.
type Consumer struct {
	*client
}

// NewConsumer creates new instance of consumer with its own connection and channels.
func NewConsumer(url string, logger Logger, channelsMaxNumber int) (*Consumer, error) {
	client, err := newClient(url, logger, channelsMaxNumber)
	if err != nil {
		return nil, err
	}

	return &Consumer{client: client}, nil
}

// Consume consumes message from a named queue.
func (c *Consumer) Consume(qname string) (<-chan amqp.Delivery, error) {
	return c.consume(qname, true, false)
}

// WorkOnTask consumes message from a named queue and defined QoS.
func (c *Consumer) WorkOnTask(qname string, prefetchCount int) (<-chan amqp.Delivery, error) {
	q, err := c.declareQueue(qname, false, false)
	if err != nil {
		return nil, err
	}

	if err := c.setQos(prefetchCount, 0); err != nil {
		return nil, err
	}

	return c.consume(q.Name, false, false)
}

// Subscribe subscribes for a messages from the exchange.
func (c *Consumer) Subscribe(meta *Meta) (<-chan amqp.Delivery, error) {
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

func (c *Consumer) consume(qname string, autoAck bool, exclusive bool) (<-chan amqp.Delivery, error) {
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
