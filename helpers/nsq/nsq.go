package nsqq

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/bitly/go-nsq"
)

var (
	nullLogger = log.New(ioutil.Discard, "", log.LstdFlags)
)

type Queue struct {
	Topic           string
	ServerAddresses []string
	Config          *nsq.Config
	Producers       map[string]*nsq.Producer
}

func NewQueue(topicname string, addresses ...[]string) *Queue {
	q := Queue{
		Topic: topicname,
	}

	if len(addresses) > 0 {
		q.ServerAddresses = addresses[0]
	}

	q.Init()
	return &q
}

func (mq *Queue) Init() error {
	if mq.Config == nil {
		mq.Config = nsq.NewConfig()
	}

	if len(mq.ServerAddresses) == 0 {
		mq.ServerAddresses = getDaemonHosts()
	}

	mq.Producers = make(map[string]*nsq.Producer)
	for _, addr := range mq.ServerAddresses {
		producer, err := nsq.NewProducer(addr, mq.Config)
		if err != nil {
			return err
		}
		producer.SetLogger(nullLogger, nsq.LogLevelInfo)
		mq.Producers[addr] = producer
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

func getDaemonHosts() []string {
	hostString := os.Getenv("NSQ_DAEMON_HOSTS")
	if hostString == "" {
		return []string{"127.0.0.1:4160"}
	}
	return strings.Split(hostString, ",")
}
