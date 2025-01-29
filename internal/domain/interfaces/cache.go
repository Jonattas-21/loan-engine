package interfaces

import "time"

type CacheRepository interface {
	Get(key string) (string, error)
	Set(key string, item []byte, ttl time.Duration) error
	Ping() error
}