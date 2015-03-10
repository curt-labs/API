package nsq

import (
	"encoding/json"
	"fmt"
	nsqq "github.com/bitly/go-nsq"
	"os"
)

func Push(topic string, data interface{}) error {
	config := nsqq.NewConfig()
	w, err := nsqq.NewProducer(getDaemonHosts(), config)
	if err != nil {
		return err
	}
	if w == nil {
		return fmt.Errorf("%s", "failed to connect to producer")
	}

	js, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	err = w.Publish(topic, js)
	if err != nil {
		return err
	}

	return nil
}

func getDaemonHosts() string {
	hostString := os.Getenv("NSQ_HOST")
	if hostString == "" {
		return "127.0.0.1:4150"
	}
	return hostString
}
