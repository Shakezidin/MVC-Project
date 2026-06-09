package model

import (
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeSavings  AccountType = "SAVINGS"
	AccountTypeCurrent  AccountType = "CURRENT"
	AccountTypeSalary   AccountType = "SALARY"
	AccountTypeFixed    AccountType = "FIXED_DEPOSIT"
)

type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "ACTIVE"
	AccountStatusInactive AccountStatus = "INACTIVE"
	AccountStatusFrozen   AccountStatus = "FROZEN"
	AccountStatusClosed   AccountStatus = "CLOSED"
)

type BankAccount struct {
	ID            uuid.UUID     `json:"account_id"`
	UserID        uuid.UUID     `json:"-"`
	AccountType   AccountType   `json:"account_type"`
	BranchName    string        `json:"branch_name"`
	AccountNumber string        `json:"-"`
	Status        AccountStatus `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type Balance struct {
	ID               uuid.UUID `json:"-"`
	AccountID        uuid.UUID `json:"account_id"`
	AvailableBalance float64   `json:"available_balance"`
	CurrentBalance   float64   `json:"current_balance"`
	Currency         string    `json:"currency"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// AccountResponse is the API representation of a bank account.
type AccountResponse struct {
	AccountID            uuid.UUID     `json:"account_id"`
	AccountType          AccountType   `json:"account_type"`
	BranchName           string        `json:"branch_name"`
	MaskedAccountNumber  string        `json:"masked_account_number"`
	Status               AccountStatus `json:"status"`
}

// BalanceResponse is the API representation of an account balance.
type BalanceResponse struct {
	AccountID        uuid.UUID `json:"account_id"`
	AvailableBalance float64   `json:"available_balance"`
	CurrentBalance   float64   `json:"current_balance"`
	Currency         string    `json:"currency"`
}
