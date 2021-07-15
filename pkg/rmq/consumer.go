package rmq

import (
	"github.com/streadway/amqp"
)

const (
	maxQueuesNumber     = 30
	defaultQueuesNumber = 4
)

// Consumer implements the consumer methods.
type Consumer struct {
	*client
	ch     *channel
	queues chan string
}

// NewConsumer creates new instance of consumer with its own connection and channels.
// It uses one channel and attaches <queuesNumber> exclusive queues to the Consumer.
func NewConsumer(url string, logger Logger, queuesNumber int) (*Consumer, error) {
	client, err := newClient(url, logger, 1)
	if err != nil {
		return nil, err
	}

	if queuesNumber > maxQueuesNumber {
		queuesNumber = maxQueuesNumber
	} else if queuesNumber < 1 {
		queuesNumber = defaultQueuesNumber
	}

	c := &Consumer{client: client}
	c, _ = c.WithFreeChannel()

	c.queues = make(chan string, queuesNumber)

	for i := 0; i < queuesNumber; i++ {
		if err := c.addExclusiveQueue(); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// newChannels creates new channels for the connection.
func (c *Consumer) addExclusiveQueue() error {
	q, err := c.ch.declareQueue("", false)
	if err != nil {
		return err
	}

	c.queues <- q.Name

	return nil
}

// WithFreeChannel returns Consumer with currently available channel
// and function to free the channel.
func (c *Consumer) WithFreeChannel() (*Consumer, func()) {
	ch := c.useFreeChannel()

	return c.withChannel(ch),
		func() { c.freeChannel(ch) }
}

// withChannel creates Consumer with channel.
func (c *Consumer) withChannel(ch *channel) *Consumer {
	return &Consumer{
		client: c.client,
		ch:     ch,
	}
}

// ListAllQueues returns the list of all queues attached to this Consumer.
func (c *Consumer) ListAllQueues() []string {
	queues := make([]string, 0, cap(c.queues))
	for i := 0; i < cap(c.queues); i++ {
		qname := <-c.queues
		queues = append(queues, qname)

		defer func(qname string) { c.queues <- qname }(qname)
	}

	return queues
}

// GetFreeQueue returns currently availabe queue name
// and function to free the queue.
func (c *Consumer) GetFreeQueue() (string, func()) {
	qname := <-c.queues
	return qname, func() { c.freeQueue(qname) }
}

// freeChannel frees the current channel so it can be used by others.
func (c *Consumer) freeQueue(qname string) {
	c.queues <- qname
}

// WorkOnTask consumes message from a named queue with defined QoS.
func (c *Consumer) WorkOnTask(qname string, prefetchCount int) (<-chan amqp.Delivery, error) {
	q, err := c.ch.declareQueue(qname, false)
	if err != nil {
		return nil, err
	}

	if err := c.ch.setQos(prefetchCount, 0); err != nil {
		return nil, err
	}

	return c.Consume(q.Name, false, false)
}

// Subscribe subscribes for a messages from the exchange.
func (c *Consumer) Subscribe(meta *Meta) (<-chan amqp.Delivery, error) {
	if err := c.ch.declareExchange(meta.XName, meta.XType, true); err != nil {
		return nil, err
	}

	q, err := c.ch.declareQueue(meta.QName, false)
	if err != nil {
		return nil, err
	}

	if err := c.ch.bindQueue(meta); err != nil {
		return nil, err
	}

	return c.Consume(q.Name, true, false)
}

func (c *Consumer) Consume(qname string, autoAck bool, exclusive bool) (<-chan amqp.Delivery, error) {
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
