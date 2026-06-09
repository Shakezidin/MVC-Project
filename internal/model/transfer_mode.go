package model

import (
	"time"

	"github.com/google/uuid"
)

type TransferMode struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TransferModeResponse is the API representation of a transfer mode.
type TransferModeResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
