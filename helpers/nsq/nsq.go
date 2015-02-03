package nsqq

import (
	"github.com/bitly/go-nsq"
)

type Queue struct {
	Topic           string
	ServerAddresses []string
	Config          *nsq.Config
	Producers       map[string]*nsq.Producer
}

func NewQueue(topicname string, addresses []string) *Queue {
	q := Queue{
		Topic:           topicname,
		ServerAddresses: addresses,
	}
	q.Init()
	return &q
}

func (mq *Queue) Init() error {
	if mq.Config == nil {
		mq.Config = nsq.NewConfig()
	}

	if len(mq.ServerAddresses) == 0 {
		mq.ServerAddresses = []string{"127.0.0.1:4160"}
	}

	mq.Producers = make(map[string]*nsq.Producer)
	for _, addr := range mq.ServerAddresses {
		producer, err := nsq.NewProducer(addr, mq.Config)
		if err == nil {
			mq.Producers[addr] = producer
		}
	}

	return nil
}

func (mq *Queue) Dispose() {
	for _, p := range mq.Producers {
		p.Stop()
	}
}

func (mq *Queue) Push(data []byte) error {
	for _, p := range mq.Producers {
		if err := p.Publish(mq.Topic, data); err != nil {
			return err
		}
	}
	return nil
}
