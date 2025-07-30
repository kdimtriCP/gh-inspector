package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

func New(cacheDir string) (Cache, error) {
	if cacheDir == "" {
		cacheDir = ".gh-inspector-cache"
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	dbPath := filepath.Join(cacheDir, "cache.db")
	return NewSQLiteCache(dbPath)
}

func GenerateKey(parts ...string) string {
	h := sha256.New()
	for _, part := range parts {
		h.Write([]byte(part))
	}
	return hex.EncodeToString(h.Sum(nil))
}
