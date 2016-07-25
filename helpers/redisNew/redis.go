package redis

import (
	"os"

	redisio "github.com/monnand/goredis"
)

var (
	RedisClient = NewRedisClient()
	RedisMaster = NewRedisMaster()
)

func NewClient(poolsize int) (c *redisio.Client) {
	redisio.MaxPoolSize = poolsize
	_c := new(redisio.Client)

	return _c
}

func NewRedisClient() *redisio.Client {
	c := NewClient(50)

	c.Addr = "127.0.0.1:6379"
	if addr := os.Getenv("REDIS_CLIENT_ADDRESS"); addr != "" {
		c.Addr = addr
	}

	c.Db = 13
	return c
}

func NewRedisMaster() *redisio.Client {
	c := NewClient(50)

	c.Addr = "127.0.0.1:6379"
	if addr := os.Getenv("REDIS_MASTER_ADDRESS"); addr != "" {
		c.Addr = addr
	}

	c.Db = 13
	return c
}
