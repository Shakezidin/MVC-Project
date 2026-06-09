package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/banking/bank-server/internal/config"
	"github.com/redis/go-redis/v9"
)

// RedisClient wraps go-redis with health check support.
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient creates a Redis connection from configuration.
func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	opts, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	if cfg.Password != "" {
		opts.Password = cfg.Password
	}
	opts.DB = cfg.DB

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &RedisClient{Client: client}, nil
}

// Ping checks Redis connectivity for health probes.
func (r *RedisClient) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

// Close gracefully shuts down the Redis client.
func (r *RedisClient) Close() error {
	return r.Client.Close()
}
