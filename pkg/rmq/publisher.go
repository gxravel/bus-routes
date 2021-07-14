package rmq

import (
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// Publisher implements the publisher methods.
type Publisher struct {
	*client
}

// NewPublisher creates new instance of publisher with its own connection and channels.
func NewPublisher(url string, logger Logger, channelsMaxNumber int) (*Publisher, error) {
	client, err := newClient(url, logger, channelsMaxNumber)
	if err != nil {
		return nil, err
	}

	return &Publisher{client: client}, nil
}

// NewTask produces message to a named queue.
func (p *Publisher) NewTask(qname string, body []byte) error {
	q, err := p.declareQueue(qname, false, false)
	if err != nil {
		return err
	}

	return p.produce(&Meta{QName: q.Name, Mode: amqp.Persistent}, body)
}

// Publish publishes a message to the exchange.
func (p *Publisher) Publish(meta *Meta, body []byte) error {
	if err := p.declareExchange(meta.XName, meta.XType, true); err != nil {
		return err
	}

	return p.produce(meta, body)
}

// CallRPC produces a message to a nammed queue, pre-installing ReplyTo property,
// which it expects to consume from.
func (p *Publisher) CallRPC(meta *Meta, body []byte) error {
	q, err := p.declareQueue("", false, false)
	if err != nil {
		return err
	}

	meta.QName = q.Name
	meta.CorrID = uuid.New().String()

	return p.produce(meta, body)
}

func (c *Publisher) Close() error {
	return c.conn.Close()
}

func (c *client) produce(meta *Meta, body []byte) error {
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
