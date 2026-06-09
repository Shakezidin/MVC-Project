package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/banking/bank-server/internal/model"
	"github.com/banking/bank-server/internal/response"
	"github.com/banking/bank-server/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type accountRepository struct {
	db *Database
}

// NewAccountRepository creates a PostgreSQL-backed account repository.
func NewAccountRepository(db *Database) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) GetByUserID(ctx context.Context, userID uuid.UUID, pagination utils.Pagination) ([]model.BankAccount, int, error) {
	countQuery := `SELECT COUNT(*) FROM bank_accounts WHERE user_id = $1`
	var total int
	if err := r.db.Pool.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count accounts: %w", err)
	}

	query := `
		SELECT id, user_id, account_type, branch_name, account_number, status, created_at, updated_at
		FROM bank_accounts
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get accounts by user: %w", err)
	}
	defer rows.Close()

	var accounts []model.BankAccount
	for rows.Next() {
		var acc model.BankAccount
		if err := rows.Scan(
			&acc.ID,
			&acc.UserID,
			&acc.AccountType,
			&acc.BranchName,
			&acc.AccountNumber,
			&acc.Status,
			&acc.CreatedAt,
			&acc.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan account: %w", err)
		}
		accounts = append(accounts, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate accounts: %w", err)
	}

	return accounts, total, nil
}

func (r *accountRepository) GetByIDAndUserID(ctx context.Context, accountID, userID uuid.UUID) (*model.BankAccount, error) {
	query := `
		SELECT id, user_id, account_type, branch_name, account_number, status, created_at, updated_at
		FROM bank_accounts
		WHERE id = $1 AND user_id = $2
	`

	var acc model.BankAccount
	err := r.db.Pool.QueryRow(ctx, query, accountID, userID).Scan(
		&acc.ID,
		&acc.UserID,
		&acc.AccountType,
		&acc.BranchName,
		&acc.AccountNumber,
		&acc.Status,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, response.NewNotFoundError("account not found")
		}
		return nil, fmt.Errorf("get account by id and user: %w", err)
	}

	return &acc, nil
}
