package redis

import (
	"os"
)

var (
	RedisClient = NewRedisClient()
	RedisMaster = NewRedisMaster()
)

func NewRedisClient() *Client {
	c := NewClient(50)

	c.Addr = "127.0.0.1:6379"
	if addr := os.Getenv("REDIS_CLIENT_ADDRESS"); addr != "" {
		c.Addr = addr
	}

	c.Db = 13
	c.Password = os.Getenv("REDIS_PASSWORD")

	return c
}

func NewRedisMaster() *Client {
	c := NewClient(50)

	c.Addr = "127.0.0.1:6379"
	if addr := os.Getenv("REDIS_MASTER_ADDRESS"); addr != "" {
		c.Addr = addr
	}

	c.Db = 13
	c.Password = os.Getenv("REDIS_PASSWORD")

	return c
}
