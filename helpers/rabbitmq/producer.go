package rabbitmq

import (
	"errors"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	config   *Config
	Exchange Exchange
}

func NewProducer(exchange Exchange, config *Config) (producer *Producer, err error) {
	if exchange.Type == "" {
		exchange.Type = "direct"
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

	producer = new(Producer)
	producer.config = config
	producer.conn = conn
	producer.Exchange = exchange
	producer.channel, err = conn.Channel()

	return
}

func (p *Producer) reconnect() error {
	var err error
	if p.config != nil {
		if p.conn, err = amqp.Dial(p.config.GetConnectionString()); err != nil {
			return err
		}

		if p.channel, err = p.conn.Channel(); err != nil {
			return err
		}
	}
	return err
}

func (p *Producer) SendMessage(mess []byte) error {
	var err error
	if p.conn == nil || p.channel == nil {
		if err = p.reconnect(); err != nil {
			return err
		}
	}
	if err = p.Exchange.Validate(); err != nil {
		return err
	}
	if len(mess) < 1 {
		return errors.New("Message cannot be empty")
	}
	err = p.channel.Publish(
		p.Exchange.Name,
		p.Exchange.RoutingKey,
		false, //mandatory?
		false, //immediate?
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "UTF-8",
			Body:            mess,
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
	)
	if err != nil {
		//problem with the connection or channel
		//clear them out and try sending the message again
		if err == amqp.ErrClosed {
			p.conn = nil
			p.channel = nil
			return p.SendMessage(mess)
		}
	}
	return err
}
