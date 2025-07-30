package cache

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteCache struct {
	db *sql.DB
}

func NewSQLiteCache(dbPath string) (*SQLiteCache, error) {
	if dbPath == "" {
		dbPath = filepath.Join(".gh-inspector", "cache.db")
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	cache := &SQLiteCache{db: db}
	if err := cache.createTable(); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			return nil, fmt.Errorf("failed to create table: %w (also failed to close db: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	go cache.cleanupExpired()

	return cache, nil
}

func (c *SQLiteCache) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS cache (
		key TEXT PRIMARY KEY,
		value BLOB NOT NULL,
		expires_at INTEGER NOT NULL,
		created_at INTEGER NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_expires_at ON cache(expires_at);
	`
	_, err := c.db.Exec(query)
	return err
}

func (c *SQLiteCache) Get(key string) ([]byte, bool, error) {
	var value []byte
	var expiresAt int64

	query := "SELECT value, expires_at FROM cache WHERE key = ?"
	err := c.db.QueryRow(query, key).Scan(&value, &expiresAt)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("failed to get cache entry: %w", err)
	}

	if time.Now().Unix() > expiresAt {
		_ = c.Delete(key)
		return nil, false, nil
	}

	return value, true, nil
}

func (c *SQLiteCache) Set(key string, value []byte, ttl time.Duration) error {
	expiresAt := time.Now().Add(ttl).Unix()
	createdAt := time.Now().Unix()

	query := `
	INSERT OR REPLACE INTO cache (key, value, expires_at, created_at)
	VALUES (?, ?, ?, ?)
	`
	_, err := c.db.Exec(query, key, value, expiresAt, createdAt)
	if err != nil {
		return fmt.Errorf("failed to set cache entry: %w", err)
	}

	return nil
}

func (c *SQLiteCache) Delete(key string) error {
	query := "DELETE FROM cache WHERE key = ?"
	_, err := c.db.Exec(query, key)
	if err != nil {
		return fmt.Errorf("failed to delete cache entry: %w", err)
	}
	return nil
}

func (c *SQLiteCache) Clear() error {
	query := "DELETE FROM cache"
	_, err := c.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}
	return nil
}

func (c *SQLiteCache) Close() error {
	return c.db.Close()
}

func (c *SQLiteCache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		query := "DELETE FROM cache WHERE expires_at < ?"
		_, _ = c.db.Exec(query, time.Now().Unix())
	}
}
