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
	if ex.Type == "" {
		return errors.New("Must specify exchange type")
	}
	if ex.RoutingKey == "" {
		return errors.New("Must specify routing key")
	}
	return nil
}
