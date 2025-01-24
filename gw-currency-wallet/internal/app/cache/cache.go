package cache

import (
	redis_store "github.com/eko/gocache/store/redis/v4"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

func getRedisAddress() string {
	e := godotenv.Load(".env")
	if e != nil {
		log.Println("No env")
	}
	host := os.Getenv("redis_host")
	port := os.Getenv("redis_port")
	return host + ":" + port
}

func InitRedis() *redis_store.RedisStore {
	return redis_store.NewRedis(redis.NewClient(&redis.Options{
		Addr: getRedisAddress(),
	}))
}
