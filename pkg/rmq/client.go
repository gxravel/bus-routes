package rmq

import (
	"sync"

	"github.com/streadway/amqp"
)

const (
	defaultContentType       = "application/json"
	maxChannelsMaxNumber     = 30
	defaultChannelsMaxNumber = 4
)

// client wraps amqp connection and implements the methods to handle channels.
type client struct {
	conn   *amqp.Connection
	logger Logger

	channels          chan *channel
	channelsMaxNumber int

	chm                   sync.Mutex
	channelsCurrentNumber int
}

// newClient creates new instance of client and adds one channel.
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

	c.channels = make(chan *channel, channelsMaxNumber)

	if err := c.addChannel(); err != nil {
		return nil, err
	}

	return c, nil
}

// addChannel creates new channel for the connection
// and adds it to the channels.
func (c *client) addChannel() error {
	ch, err := c.conn.Channel()
	if err != nil {
		c.logger.Fatal("failed to open a channel")
		return err
	}

	c.channels <- &channel{ch}
	c.channelsCurrentNumber += 1

	return nil
}

// useFreeChannel returns currently available channel.
// If possible it creates new one.
func (c *client) useFreeChannel() *channel {
	for {
		select {
		case ch := <-c.channels:
			return ch

		default:
			c.chm.Lock()
			if c.channelsCurrentNumber < c.channelsMaxNumber {
				if err := c.addChannel(); err != nil {
					c.logger.Fatalf(err.Error(), "could not create new channel")
				}
			}
			c.chm.Unlock()
		}
	}
}

// freeChannel frees the current channel so it can be used by others.
func (c *client) freeChannel(ch *channel) {
	if ch.Channel != nil {
		c.channels <- ch
	} else {
		c.chm.Lock()
		c.channelsCurrentNumber -= 1
		c.chm.Unlock()
	}
}

func (c *client) Close() error {
	return c.conn.Close()
}

type channel struct {
	*amqp.Channel
}

func (ch *channel) declareQueue(name string, durable bool) (amqp.Queue, error) {
	var exclusive = name == ""

	q, err := ch.QueueDeclare(
		name,
		durable,
		false,     // delete when unused
		exclusive, // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return q, err
	}

	return q, nil
}

func (ch *channel) bindQueue(meta *Meta) error {
	if err := ch.QueueBind(
		meta.QName,
		meta.Key, // routing key
		meta.XName,
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (ch *channel) declareExchange(name string, xtype string, durable bool) error {
	if err := ch.ExchangeDeclare(
		name,
		xtype,
		durable,
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	); err != nil {
		return err
	}

	return nil
}

func (ch *channel) setQos(prefetchCount int, prefetchSize int) error {
	if err := ch.Qos(
		prefetchCount,
		prefetchSize,
		false, // global
	); err != nil {
		return err
	}

	return nil
}
