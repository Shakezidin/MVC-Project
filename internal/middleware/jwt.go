package middleware

import (
	"net/http"
	"strings"

	"github.com/banking/bank-server/internal/auth"
	"github.com/banking/bank-server/internal/response"
	"github.com/banking/bank-server/internal/utils"
	"go.uber.org/zap"
)

// JWTAuth validates Bearer tokens and injects user context.
func JWTAuth(tokenService *auth.TokenService, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			logger.Info("Authorization is " + authHeader)
			if authHeader == "" {
				response.Error(w, r, http.StatusUnauthorized,
					"authorization header required", response.ErrCodeUnauthorized, nil)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Error(w, r, http.StatusUnauthorized,
					"invalid authorization header format", response.ErrCodeUnauthorized, nil)
				return
			}

			claims, err := tokenService.ValidateToken(parts[1])
			if err != nil {
				response.Error(w, r, http.StatusUnauthorized,
					"invalid or expired token", response.ErrCodeUnauthorized, nil)
				return
			}

			ctx := utils.WithUserID(r.Context(), claims.UserID.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
