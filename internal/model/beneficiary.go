package model

import (
	"time"

	"github.com/google/uuid"
)

type Beneficiary struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"-"`
	BeneficiaryName string  `json:"beneficiary_name"`
	BankName      string    `json:"bank_name"`
	AccountNumber string    `json:"-"`
	IFSC          string    `json:"ifsc"`
	Nickname      string    `json:"nickname"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BeneficiaryResponse is the API representation of a beneficiary.
type BeneficiaryResponse struct {
	BeneficiaryName     string `json:"beneficiary_name"`
	BankName            string `json:"bank_name"`
	AccountNumberMasked string `json:"account_number_masked"`
	IFSC                string `json:"ifsc"`
	Nickname            string `json:"nickname"`
}
