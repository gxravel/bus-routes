package rmq

import (
	"github.com/streadway/amqp"
)

// Publisher implements the publisher methods.
type Publisher struct {
	*client
	ch *channel
}

// NewPublisher creates new instance of publisher with its own connection and channels.
func NewPublisher(url string, logger Logger, channelsMaxNumber int) (*Publisher, error) {
	client, err := newClient(url, logger, channelsMaxNumber)
	if err != nil {
		return nil, err
	}

	return &Publisher{client: client}, nil
}

// withChannel creates Publisher with channel.
func (p *Publisher) withChannel(ch *channel) *Publisher {
	return &Publisher{
		client: p.client,
		ch:     ch,
	}
}

// WithFreeChannel returns Publisher with channel and function to free the channel.
func (p *Publisher) WithFreeChannel() (*Publisher, func()) {
	ch := p.useFreeChannel()

	return p.withChannel(ch),
		func() { p.freeChannel(ch) }
}

// NewTask produces a message to a named queue.
func (p *Publisher) NewTask(qname string, body []byte) error {
	q, err := p.ch.declareQueue(qname, true)
	if err != nil {
		return err
	}

	return p.Produce(&Meta{QName: q.Name, Mode: amqp.Persistent}, body)
}

// Publish publishes a message to the exchange.
func (p *Publisher) Publish(meta *Meta, body []byte) error {
	if err := p.ch.declareExchange(meta.XName, meta.XType, true); err != nil {
		return err
	}

	return p.Produce(meta, body)
}

// Produce produces a message.
func (p *Publisher) Produce(meta *Meta, body []byte) error {
	pub := amqp.Publishing{
		DeliveryMode:  meta.Mode,
		ContentType:   defaultContentType,
		CorrelationId: meta.CorrID,
		ReplyTo:       meta.QName,
		Body:          body,
	}

	if err := p.ch.Publish(
		meta.XName, // exchange
		meta.Key,   // routing key
		false,      // mandatory
		false,      // immediate
		pub,
	); err != nil {
		p.logger.Error("failed to produce a message")
		return err
	}

	return nil
}
