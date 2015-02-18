package rabbitmq

import (
	"errors"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
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
	producer.conn = conn
	producer.Exchange = exchange
	producer.channel, err = conn.Channel()

	return
}

func (p *Producer) SendMessage(mess []byte) error {
	if p.channel == nil {
		return errors.New("Invalid channel")
	}
	if err := p.Exchange.Validate(); err != nil {
		return err
	}
	if len(mess) < 1 {
		return errors.New("Message cannot be empty")
	}
	return p.channel.Publish(
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
}
