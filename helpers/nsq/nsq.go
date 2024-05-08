package nsq

import (
	"encoding/json"
	"fmt"
	nsqq "github.com/nsqio/go-nsq"
	"os"
)

var (
	NsqHost = os.Getenv("NSQ_HOST")
)

type nopLogger struct{}

func (*nopLogger) Output(int, string) error {
	return nil
}

func Push(topic string, data interface{}) error {
	config := nsqq.NewConfig()
	w, err := nsqq.NewProducer(getDaemonHosts(), config)
	if w == nil && err == nil {
		return fmt.Errorf("%s", "failed to connect to producer")
	}
	w.SetLogger(&nopLogger{}, nsqq.LogLevelError)
	defer w.Stop()

	if err != nil {
		return err
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
	if NsqHost == "" {
		return "127.0.0.1:4150"
	}
	return NsqHost
}
