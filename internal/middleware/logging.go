package middleware

import (
	"net/http"
	"time"

	"github.com/banking/bank-server/internal/utils"
	"go.uber.org/zap"
)

// responseWriter wraps http.ResponseWriter to capture status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging logs each request with structured fields including latency and status.
func Logging(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			latency := time.Since(start)
			requestID := utils.RequestIDFromContext(r.Context())
			userID := utils.UserIDFromContext(r.Context())

			reqLog := log.With(zap.String("request_id", requestID), zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.String("user_id", userID))
			reqLog.Info("request completed",
				zap.Int("status_code", wrapped.statusCode),
				zap.Duration("latency", latency),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			)

		})
	}
}
