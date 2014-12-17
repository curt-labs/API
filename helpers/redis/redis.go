package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	redix "github.com/garyburd/redigo/redis"
	"os"
	"time"
)

const (
	Db                = 13
	PoolAllocationErr = "failed to allocate pool"
	Prefix            = "goapi"
	CacheTimeout      = 86400
)

func RedisPool(master bool) *redix.Pool {
	addr := "127.0.0.1:6379"
	password := os.Getenv("REDIS_PASSWORD")
	if master {
		if ad := os.Getenv("REDIS_MASTER_ADDRESS"); ad != "" {
			addr = ad
		}
	} else {
		if ad := os.Getenv("REDIS_CLIENT_ADDRESS"); ad != "" {
			addr = ad
		}
		if ad := os.Getenv("REDIS_SLAVE_SERVICE_HOST"); ad != "" {
			addr = fmt.Sprintf("%s:%s", ad, os.Getenv("REDIS_SLAVE_SERVICE_PORT"))
		}
	}

	return &redix.Pool{
		MaxIdle:     2,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redix.Conn, error) {
			c, err := redix.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redix.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func Get(key string) ([]byte, error) {
	data := make([]byte, 0)
	pool := RedisPool(false)
	if pool == nil {
		return data, errors.New(PoolAllocationErr)
	}

	conn := pool.Get()
	conn.Send("select", Db)

	return redix.Bytes(conn.Do("GET", fmt.Sprintf("%s:%s", Prefix, key)))
}

func Setex(key string, obj interface{}, exp int) error {
	pool := RedisPool(true)
	if pool == nil {
		return errors.New(PoolAllocationErr)
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	conn := pool.Get()
	if err := conn.Send("select", Db); err != nil {
		return err
	}

	_, err = conn.Do("SETEX", fmt.Sprintf("%s:%s", Prefix, key), exp, data)
	return err
}

func Set(key string, obj interface{}) error {
	pool := RedisPool(true)
	if pool == nil {
		return errors.New(PoolAllocationErr)
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	conn := pool.Get()
	if err := conn.Send("select", Db); err != nil {
		return err
	}

	_, err = conn.Do("SET", fmt.Sprintf("%s:%s", Prefix, key), data)
	return err
}

func Lpush(key string, obj interface{}) error {
	pool := RedisPool(true)
	if pool == nil {
		return errors.New(PoolAllocationErr)
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	conn := pool.Get()
	if err := conn.Send("select", Db); err != nil {
		return err
	}

	_, err = conn.Do("LPUSH", fmt.Sprintf("%s:%s", Prefix, key), data)
	return err
}

func Delete(key string) error {
	var err error
	pool := RedisPool(true)
	if pool == nil {
		return errors.New(PoolAllocationErr)
	}

	conn := pool.Get()
	if err := conn.Send("select", Db); err != nil {
		return err
	}

	_, err = conn.Do("DEL", fmt.Sprintf("%s:%s", Prefix, key))
	return err
}
