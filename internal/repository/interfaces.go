package repository

import (
	"context"

	"github.com/banking/bank-server/internal/model"
	"github.com/banking/bank-server/internal/utils"
	"github.com/google/uuid"
)

// UserRepository defines data access operations for users.
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

// AccountRepository defines data access operations for bank accounts.
type AccountRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID, pagination utils.Pagination) ([]model.BankAccount, int, error)
	GetByIDAndUserID(ctx context.Context, accountID, userID uuid.UUID) (*model.BankAccount, error)
}

// BalanceRepository defines data access operations for balances.
type BalanceRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Balance, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) (*model.Balance, error)
}

// BeneficiaryRepository defines data access operations for beneficiaries.
type BeneficiaryRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID, pagination utils.Pagination) ([]model.Beneficiary, int, error)
}

// TransferModeRepository defines data access operations for transfer modes.
type TransferModeRepository interface {
	GetAllActive(ctx context.Context) ([]model.TransferMode, error)
}
