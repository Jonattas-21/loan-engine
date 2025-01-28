package repositories

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisRepository struct {
	Redis *redis.Client
}

func (r *RedisRepository) Set(key string, value interface{}, ttl time.Duration) error {
	err := r.Redis.Set(key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) Get(key string) (interface{}, error) {
	val, err := r.Redis.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
