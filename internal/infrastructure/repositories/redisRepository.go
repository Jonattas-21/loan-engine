package repositories

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisRepository struct {
	redis *redis.Client
}

func (r *RedisRepository) Set(key string, value interface{}, ttl time.Duration) error {
	err := r.redis.Set(key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) Get(key string) (interface{}, error) {
	val, err := r.redis.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
