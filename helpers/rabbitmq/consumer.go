package rabbitmq

import (
	"errors"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	config   *Config
	Name     string
	Exchange Exchange

	incomingMessages <-chan amqp.Delivery
	DoneChan         chan error
}

type Handler interface {
	HandleMessage(message *amqp.Delivery) error
}

type HandlerFunc func(message *amqp.Delivery) error

func (h HandlerFunc) HandleMessage(m *amqp.Delivery) error {
	return h(m)
}

func (c *Consumer) AddHandler(handler Handler) {
	go c.handlerLoop(handler)
}

func (c *Consumer) handlerLoop(handler Handler) {
	for msg := range c.incomingMessages {
		handler.HandleMessage(&msg)
		msg.Ack(false)
	}
	c.DoneChan <- nil
}

func (c *Consumer) Close() error {
	var err error
	if c.channel != nil {
		if err = c.channel.Cancel(c.Name, true); err != nil {
			return err
		}
	}
	if c.conn != nil {
		if err = c.conn.Close(); err != nil {
			return err
		}
	}

	return nil
}

func NewConsumer(consumerName string, queueName string, exchange Exchange, config *Config) (consumer *Consumer, err error) {
	if consumerName == "" {
		err = errors.New("Must give the consumer a name")
		return
	}
	if queueName == "" {
		err = errors.New("Must give the queue name.")
		return
	}
	if len(queueName) > 255 {
		err = errors.New("Must give a queue name that contains 1-255 characters.")
		return
	}

	if err = exchange.Validate(); err != nil {
		return
	}

	if config == nil {
		config = NewConfig()
	}

	conn, err := amqp.Dial(config.GetConnectionString())
	if err != nil {
		return
	}

	//setup the channel
	var channel *amqp.Channel
	if channel, err = conn.Channel(); err != nil {
		return
	}

	//setup the exchange
	if err = channel.ExchangeDeclare(
		exchange.Name, //exchange name
		exchange.Type, //exchange type
		true,          //durable
		false,         //remove when complete
		false,         //internal
		false,         //noWait
		nil,           //arguments
	); err != nil {
		return
	}

	//setup the queue
	var queue amqp.Queue
	if queue, err = channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	); err != nil {
		return
	}
	if err = channel.QueueBind(
		queue.Name,          //queue name
		exchange.RoutingKey, //routing key ("binding key")
		exchange.Name,       //exchange (source)
		false,               //noWait
		nil,                 //arguments
	); err != nil {
		return
	}

	//setup the deliverables
	var messages <-chan amqp.Delivery

	if messages, err = channel.Consume(
		queue.Name,   //queue name
		consumerName, //consumer name
		false,        //auto acknowledge
		false,        //exclusive
		false,        //not local
		false,        //no wait
		nil,          //arguments
	); err != nil {
		return
	}

	consumer = new(Consumer)
	consumer.Name = consumerName
	consumer.config = config
	consumer.conn = conn
	consumer.channel = channel
	consumer.Exchange = exchange
	consumer.incomingMessages = messages

	return
}
