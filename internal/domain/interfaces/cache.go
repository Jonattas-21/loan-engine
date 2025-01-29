package interfaces

import "time"

type CacheRepository interface {
	Get(key string) (interface{}, error)
	Set(key string, item interface{}, ttl time.Duration) error
	Ping() error
}