package rabbitmq

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Username string
	Password string
	Hostname string
	Port     int
}

func NewConfig() *Config {
	c := new(Config)
	c.Hostname = "localhost"
	c.Port = 5672
	if os.Getenv("AMQP_HOST") != "" {
		c.Hostname = os.Getenv("AMQP_HOST")
	}
	if os.Getenv("AMQP_PORT") != "" {
		if port, err := strconv.Atoi(os.Getenv("AMQP_PORT")); err == nil {
			c.Port = port
		}
	}
	if os.Getenv("AMQP_USER") != "" {
		c.Username = os.Getenv("AMQP_USER")
	}
	if os.Getenv("AMQP_PASSWORD") != "" {
		c.Password = os.Getenv("AMQP_PASSWORD")
	}
	return c
}

func (c *Config) GetConnectionString() string {
	if c.Username != "" && c.Password != "" {
		return fmt.Sprintf("amqp://%s:%s@%s:%d", c.Username, c.Password, c.Hostname, c.Port)
	}
	return fmt.Sprintf("amqp://%s:%d", c.Hostname, c.Port)
}
