package cache

import (
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
	"os"
)

func getRedisAddress() string {
	host := os.Getenv("redis_host")
	port := os.Getenv("redis_port")
	return host + ":" + port
}

func InitRedis() *redis_store.RedisStore {
	return redis_store.NewRedis(redis.NewClient(&redis.Options{
		Addr: getRedisAddress(),
	}))
}
