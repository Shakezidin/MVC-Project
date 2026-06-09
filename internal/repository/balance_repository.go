package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/banking/bank-server/internal/model"
	"github.com/banking/bank-server/internal/response"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type balanceRepository struct {
	db *Database
}

// NewBalanceRepository creates a PostgreSQL-backed balance repository.
func NewBalanceRepository(db *Database) BalanceRepository {
	return &balanceRepository{db: db}
}

func (r *balanceRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Balance, error) {
	query := `
		SELECT b.id, b.account_id, b.available_balance, b.current_balance, b.currency, b.created_at, b.updated_at
		FROM balances b
		INNER JOIN bank_accounts a ON a.id = b.account_id
		WHERE a.user_id = $1
		ORDER BY b.created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get balances by user: %w", err)
	}
	defer rows.Close()

	var balances []model.Balance
	for rows.Next() {
		var bal model.Balance
		if err := rows.Scan(
			&bal.ID,
			&bal.AccountID,
			&bal.AvailableBalance,
			&bal.CurrentBalance,
			&bal.Currency,
			&bal.CreatedAt,
			&bal.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan balance: %w", err)
		}
		balances = append(balances, bal)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate balances: %w", err)
	}

	return balances, nil
}

func (r *balanceRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) (*model.Balance, error) {
	query := `
		SELECT id, account_id, available_balance, current_balance, currency, created_at, updated_at
		FROM balances
		WHERE account_id = $1
	`

	var bal model.Balance
	err := r.db.Pool.QueryRow(ctx, query, accountID).Scan(
		&bal.ID,
		&bal.AccountID,
		&bal.AvailableBalance,
		&bal.CurrentBalance,
		&bal.Currency,
		&bal.CreatedAt,
		&bal.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response.NewNotFoundError("balance not found")
		}
		return nil, fmt.Errorf("get balance by account: %w", err)
	}

	return &bal, nil
}
