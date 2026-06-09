package handler

import (
	"net/http"

	"github.com/banking/bank-server/internal/response"
	"github.com/banking/bank-server/internal/utils"
	"github.com/google/uuid"
)

// handleError maps domain errors to HTTP responses.
func handleError(w http.ResponseWriter, r *http.Request, err error) {
	appErr := response.AsAppError(err)
	response.Error(w, r, appErr.HTTPStatus, appErr.Message, appErr.Code, appErr.Details)
}

// userIDFromContext extracts and parses the authenticated user ID.
func userIDFromContext(r *http.Request) (uuid.UUID, error) {
	userIDStr := utils.UserIDFromContext(r.Context())
	if userIDStr == "" {
		return uuid.Nil, response.NewUnauthorizedError("user not authenticated")
	}
	return uuid.Parse(userIDStr)
}
