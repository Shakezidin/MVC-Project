package middleware

import (
	"net/http"

	"github.com/banking/bank-server/internal/utils"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

// RequestID injects a unique request ID into context and response headers.
// Request IDs enable distributed tracing correlation across logs and services.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = "req-" + uuid.New().String()
		}

		ctx := utils.WithRequestID(r.Context(), requestID)
		w.Header().Set(RequestIDHeader, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
