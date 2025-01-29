package repositories

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type RedisRepository struct {
	Redis  *redis.Client
	Logger *logrus.Logger
}

func (r *RedisRepository) Set(key string, value []byte, ttl time.Duration) error {
	err := r.Redis.Set(key, value, ttl).Err()
	if err != nil {
		return err
	}
	r.Logger.Infoln(fmt.Printf("Value set in cache key %v", key))
	return nil
}

func (r *RedisRepository) Get(key string) (string, error) {
	val, err := r.Redis.Get(key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisRepository) Ping() error {
	_, err := r.Redis.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
