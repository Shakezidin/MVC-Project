package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache provides generic JSON serialization helpers for Redis.
type Cache struct {
	client *RedisClient
}

// New creates a Cache wrapper around RedisClient.
func New(client *RedisClient) *Cache {
	return &Cache{client: client}
}

// Get retrieves and deserializes a cached value. Returns false on cache miss.
func (c *Cache) Get(ctx context.Context, key string, dest interface{}) (bool, error) {
	data, err := c.client.Client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, fmt.Errorf("cache get: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return false, fmt.Errorf("cache unmarshal: %w", err)
	}
	return true, nil
}

// Set serializes and stores a value with TTL.
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal: %w", err)
	}

	if err := c.client.Client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("cache set: %w", err)
	}
	return nil
}

// Delete removes a key from cache.
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Client.Del(ctx, key).Err()
}
