package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	redix "github.com/garyburd/redigo/redis"
)

const (
	PoolAllocationErr = "failed to allocate pool"
	Prefix            = "API"
	CacheTimeout      = 86400
)

func RedisPool(master bool) *redix.Pool {
	addr := "127.0.0.1:6379"
	password := os.Getenv("REDIS_PASSWORD")

	if master && os.Getenv("REDIS_MASTER_ADDRESS") != "" {
		addr = os.Getenv("REDIS_MASTER_ADDRESS")
	} else if os.Getenv("REDIS_MASTER_ADDRESS") != "" {
		addr = os.Getenv("REDIS_SLAVE_ADDRESS")
	}

	if master && os.Getenv("REDIS_MASTER_SERVICE_HOST") != "" {
		addr = fmt.Sprintf("%s", os.Getenv("REDIS_MASTER_SERVICE_HOST"))
	} else if os.Getenv("REDIS_SLAVE_SERVICE_HOST") != "" {
		addr = fmt.Sprintf("%s", os.Getenv("REDIS_SLAVE_SERVICE_HOST"))
	}
	return &redix.Pool{
		MaxIdle:     2,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redix.Conn, error) {
			c, err := redix.Dial("tcp", fmt.Sprintf("%s", addr))
			if err != nil {
				return nil, err
			}
			if password != "" && master {
				if _, err = c.Do("AUTH", password); err != nil {
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

	conn, err := pool.Dial()
	if err != nil {
		return data, err
	} else if conn.Err() != nil {
		return data, err
	}

	reply, err := conn.Do("GET", fmt.Sprintf("%s:%s", Prefix, key))
	if err != nil || reply == nil {
		return data, err
	}

	return redix.Bytes(reply, err)
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

	if pool == nil {
		return errors.New(PoolAllocationErr)
	}
	conn := pool.Get()
	if conn.Err() != nil {
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
	if conn.Err() != nil {
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
	if conn.Err() != nil {
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
	if conn.Err() != nil {
		return err
	}

	_, err = conn.Do("DEL", fmt.Sprintf("%s:%s", Prefix, key))
	return err
}

//Goadmin calls
func GetNamespaces() (namespaces map[string][]string, err error) {
	pool := RedisPool(true)
	if pool == nil {
		return namespaces, errors.New(PoolAllocationErr)
	}

	conn := pool.Get()
	if conn.Err() != nil {
		return namespaces, err
	}
	reply, err := redix.Strings(conn.Do("KEYS", "*"))

	namespaces = make(map[string][]string, 0)

	for _, key := range reply {
		keyArr := strings.Split(key, ":")
		if len(keyArr) > 0 {
			idx := keyArr[0]
			if _, ok := namespaces[idx]; !ok {
				namespaces[idx] = make([]string, 0)
			}
			namespaces[idx] = append(namespaces[idx], key)
		}
	}
	return
}

func DeleteFullPath(key string) error {
	var err error
	pool := RedisPool(true)
	if pool == nil {
		return errors.New(PoolAllocationErr)
	}

	conn := pool.Get()
	if conn.Err() != nil {
		return err
	}

	_, err = conn.Do("DEL", fmt.Sprintf("%s", key))
	return err
}

func GetFullPath(key string) ([]string, error) {
	data := make([]string, 0)
	pool := RedisPool(false)

	if pool == nil {
		return data, errors.New(PoolAllocationErr)
	}

	conn, err := pool.Dial()
	if err != nil {
		return data, err
	} else if conn.Err() != nil {
		return data, err
	}

	reply, err := redix.Strings(conn.Do("KEYS", "*"))
	if err != nil || reply == nil {
		return data, err
	}

	for _, r := range reply {
		if strings.Contains(r, key) {
			data = append(data, r)
		}
	}
	return data, err
}
