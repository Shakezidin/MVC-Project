package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/banking/bank-server/internal/response"
	"github.com/banking/bank-server/internal/utils"
	"go.uber.org/zap"
)

// Recovery catches panics and returns a 500 without crashing the server.
// This is essential in production to maintain availability during unexpected failures.
func Recovery(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					requestID := utils.RequestIDFromContext(r.Context())
					log.Error("panic recovered",
						zap.Any("panic", rec),
						zap.String("request_id", requestID),
						zap.String("stack", string(debug.Stack())),
					)
					response.Error(w, r, http.StatusInternalServerError,
						"something went wrong", response.ErrCodeInternal, nil)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
