package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/banking/bank-server/internal/response"
)

// Timeout enforces a maximum request duration to prevent resource exhaustion.
func Timeout(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()

			done := make(chan struct{})
			panicChan := make(chan interface{}, 1)

			go func() {
				defer func() {
					if rec := recover(); rec != nil {
						panicChan <- rec
					}
				}()
				next.ServeHTTP(w, r.WithContext(ctx))
				close(done)
			}()

			select {
			case p := <-panicChan:
				panic(p)
			case <-done:
				return
			case <-ctx.Done():
				if ctx.Err() == context.DeadlineExceeded {
					response.Error(w, r, http.StatusGatewayTimeout,
						"request timed out", response.ErrCodeRequestTimeout, nil)
				}
			}
		})
	}
}
