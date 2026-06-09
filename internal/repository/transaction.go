package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// TxFunc executes a function within a database transaction.
// Automatically rolls back on error and commits on success.
func (db *Database) TxFunc(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}
