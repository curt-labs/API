package redis

var (
	RedisClient = NewRedisClient()
)

func NewRedisClient() (c Client) {

	// client.Addr = "137.117.72.189:6379"
	c.Addr = "127.0.0.1:6379"
	c.Db = 13

	return
}
