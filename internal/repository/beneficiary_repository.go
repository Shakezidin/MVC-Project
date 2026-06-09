package repository

import (
	"context"
	"fmt"

	"github.com/banking/bank-server/internal/model"
	"github.com/banking/bank-server/internal/utils"
	"github.com/google/uuid"
)

type beneficiaryRepository struct {
	db *Database
}

// NewBeneficiaryRepository creates a PostgreSQL-backed beneficiary repository.
func NewBeneficiaryRepository(db *Database) BeneficiaryRepository {
	return &beneficiaryRepository{db: db}
}

func (r *beneficiaryRepository) GetByUserID(ctx context.Context, userID uuid.UUID, pagination utils.Pagination) ([]model.Beneficiary, int, error) {
	countQuery := `SELECT COUNT(*) FROM beneficiaries WHERE user_id = $1`
	var total int
	if err := r.db.Pool.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count beneficiaries: %w", err)
	}

	query := `
		SELECT id, user_id, beneficiary_name, bank_name, account_number, ifsc, nickname, created_at, updated_at
		FROM beneficiaries
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get beneficiaries by user: %w", err)
	}
	defer rows.Close()

	var beneficiaries []model.Beneficiary
	for rows.Next() {
		var b model.Beneficiary
		if err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.BeneficiaryName,
			&b.BankName,
			&b.AccountNumber,
			&b.IFSC,
			&b.Nickname,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan beneficiary: %w", err)
		}
		beneficiaries = append(beneficiaries, b)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate beneficiaries: %w", err)
	}

	return beneficiaries, total, nil
}
