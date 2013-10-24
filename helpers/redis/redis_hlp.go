package redis

var (
	RedisClient *Client
)

func GetClient() *Client {
	if RedisClient == nil {
		RedisClient = NewRedisClient()
	}
	return RedisClient
}

func NewRedisClient() *Client {
	c := NewClient(50)

	// c.Addr = "168.61.40.178:6379"
	c.Addr = "127.0.0.1:6379"
	c.Db = 13

	return c
}
