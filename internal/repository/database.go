package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/banking/bank-server/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database wraps a pgx connection pool with lifecycle management.
type Database struct {
	Pool *pgxpool.Pool
}

// NewDatabase creates a context-aware PostgreSQL connection pool.
// Connection pooling is critical in production to avoid exhausting DB connections.
func NewDatabase(ctx context.Context, cfg config.DatabaseConfig) (*Database, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &Database{Pool: pool}, nil
}

// Close gracefully shuts down the connection pool.
func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

// Ping checks database connectivity for health probes.
func (db *Database) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return db.Pool.Ping(ctx)
}
