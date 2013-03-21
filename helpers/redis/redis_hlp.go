package redis

var (
	client Client
)

func NewRedisClient() Client {
	// client.Addr = "137.117.72.189:6379"
	client.Addr = "127.0.0.1:6379"
	client.Db = 13

	return client
}