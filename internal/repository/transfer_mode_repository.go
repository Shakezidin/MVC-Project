package repository

import (
	"context"
	"fmt"

	"github.com/banking/bank-server/internal/model"
)

type transferModeRepository struct {
	db *Database
}

// NewTransferModeRepository creates a PostgreSQL-backed transfer mode repository.
func NewTransferModeRepository(db *Database) TransferModeRepository {
	return &transferModeRepository{db: db}
}

func (r *transferModeRepository) GetAllActive(ctx context.Context) ([]model.TransferMode, error) {
	query := `
		SELECT id, code, name, description, is_active, created_at, updated_at
		FROM transfer_modes
		WHERE is_active = true
		ORDER BY name ASC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get transfer modes: %w", err)
	}
	defer rows.Close()

	var modes []model.TransferMode
	for rows.Next() {
		var m model.TransferMode
		if err := rows.Scan(
			&m.ID,
			&m.Code,
			&m.Name,
			&m.Description,
			&m.IsActive,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan transfer mode: %w", err)
		}
		modes = append(modes, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate transfer modes: %w", err)
	}

	return modes, nil
}
