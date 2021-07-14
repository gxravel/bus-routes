package rmq

import (
	"github.com/streadway/amqp"
)

const (
	defaultContentType       = "application/json"
	maxChannelsMaxNumber     = 30
	defaultChannelsMaxNumber = 4
)

// client wraps connection and channel of amqp.
type client struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger Logger

	channels          chan *amqp.Channel
	channelsMaxNumber int
}

// newClient creates new instance of client with its own connection.
func newClient(url string, logger Logger, channelsMaxNumber int) (*client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Fatal("failed to connect to RMQ")
		return nil, err
	}

	if channelsMaxNumber > maxChannelsMaxNumber {
		channelsMaxNumber = maxChannelsMaxNumber
	} else if channelsMaxNumber < 1 {
		channelsMaxNumber = defaultChannelsMaxNumber
	}

	c := &client{
		conn:              conn,
		logger:            logger,
		channelsMaxNumber: channelsMaxNumber,
	}

	c.channels = make(chan *amqp.Channel, channelsMaxNumber)

	if err := c.newChannels(1); err != nil {
		return nil, err
	}

	return c, nil
}

// newChannels creates new channels for the connection.
func (c *client) newChannels(count int) error {
	if len(c.channels)+count > c.channelsMaxNumber {
		count = c.channelsMaxNumber - len(c.channels)
	}

	for i := 0; i < count; i++ {
		ch, err := c.conn.Channel()
		if err != nil {
			c.logger.Fatal("failed to open a channel")
			return err
		}
		c.channels <- ch
	}

	return nil
}

// UseFreeChannel anchors the client with free channel.
// It creates new channel if possible.
func (c *client) UseFreeChannel() {
	select {
	case ch := <-c.channels:
		c.useChannel(ch)
	default:
		if cap(c.channels) < c.channelsMaxNumber {
			if err := c.newChannels(1); err != nil {
				c.logger.Fatalf(err.Error(), "could not create new channel")
			}

			c.useChannel(<-c.channels)
		}
	}
}

// FreeChannel frees the current channel so it can be used by others.
func (c *client) FreeChannel() {
	c.channels <- c.ch
}

// useChannel anchors the client with channel.
func (c *client) useChannel(ch *amqp.Channel) {
	c.ch = ch
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

func (c *client) Close() error {
	return c.conn.Close()
}
