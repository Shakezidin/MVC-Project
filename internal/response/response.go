package response

import (
	"encoding/json"
	"net/http"

	"github.com/banking/bank-server/internal/utils"
)

// APIResponse is the standard envelope for all API responses.
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorBody  `json:"error,omitempty"`
}

// ErrorBody contains structured error details.
type ErrorBody struct {
	Code    string            `json:"code"`
	Details map[string]string `json:"details,omitempty"`
}

// PaginatedData wraps paginated list responses.
type PaginatedData struct {
	Items      interface{} `json:"items"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalCount int         `json:"total_count"`
	TotalPages int         `json:"total_pages"`
}

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, r *http.Request, status int, message string, data interface{}) {
	requestID := utils.RequestIDFromContext(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(status)

	resp := APIResponse{
		Success:   status >= 200 && status < 300,
		Message:   message,
		RequestID: requestID,
		Data:      data,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// Success writes a successful JSON response.
func Success(w http.ResponseWriter, r *http.Request, status int, message string, data interface{}) {
	JSON(w, r, status, message, data)
}

// Error writes an error JSON response.
func Error(w http.ResponseWriter, r *http.Request, status int, message, code string, details map[string]string) {
	requestID := utils.RequestIDFromContext(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Request-ID", requestID)
	w.WriteHeader(status)

	resp := APIResponse{
		Success:   false,
		Message:   message,
		RequestID: requestID,
		Error: &ErrorBody{
			Code:    code,
			Details: details,
		},
	}

	_ = json.NewEncoder(w).Encode(resp)
}
