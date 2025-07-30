package cache

import (
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=../mock/mock_cache/mock_$GOFILE -package=mock_cache

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
