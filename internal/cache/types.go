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

type CacheEntry struct {
	Key       string
	Value     []byte
	ExpiresAt time.Time
	CreatedAt time.Time
}
