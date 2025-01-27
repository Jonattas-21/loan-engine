package cache

import (
	"github.com/go-redis/redis"
	"os"
	"fmt"
)

func NewCache() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return rdb
}