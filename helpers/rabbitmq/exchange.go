package rabbitmq

import (
	"errors"
)

type Exchange struct {
	Name       string
	Type       string
	RoutingKey string
}

func (ex *Exchange) Validate() error {
	if ex.Name == "" {
		return errors.New("Must specify exchange name")
	}

	//exchange types can be: direct, fanout, topic, headers
	//we're only going to support direct for now
	if ex.Type != "direct" {
		ex.Type = "direct"
	}

	if ex.Type == "direct" && ex.RoutingKey == "" {
		return errors.New("Must specify routing key")
	}

	return nil
}
