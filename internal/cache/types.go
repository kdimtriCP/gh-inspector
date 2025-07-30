package cache

import (
	"time"
)

type Cache interface {
	Get(key string) ([]byte, bool, error)
	Set(key string, value []byte, ttl time.Duration) error
	Delete(key string) error
	Clear() error
	Close() error
}
