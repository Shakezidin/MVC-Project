package router

import (
	"net/http"

	"github.com/banking/bank-server/internal/auth"
	"github.com/banking/bank-server/internal/config"
	"github.com/banking/bank-server/internal/handler"
	"github.com/banking/bank-server/internal/middleware"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// Handlers groups all HTTP handlers for dependency injection.
type Handlers struct {
	Auth         *handler.AuthHandler
	Account      *handler.AccountHandler
	Beneficiary  *handler.BeneficiaryHandler
	TransferMode *handler.TransferModeHandler
	Health       *handler.HealthHandler
}

// New creates the application router with middleware chain and route registration.
func New(cfg *config.Config, handlers Handlers, tokenService *auth.TokenService, log *zap.Logger) http.Handler {
	r := mux.NewRouter()

	// Global middleware chain — order matters:
	// RequestID → Recovery → Logging → Security → CORS → Timeout → RateLimit
	rateLimiter := middleware.NewRateLimiter(cfg.Server.RateLimitRPS, cfg.Server.RateLimitBurst)

	r.Use(
		middleware.RequestID,
		middleware.Recovery(log),
		middleware.Logging(log),
		middleware.SecureHeaders,
		middleware.CORS,
		middleware.Timeout(cfg.Server.RequestTimeout),
		rateLimiter.Limit,
	)

	// Health endpoints (no auth required)
	r.HandleFunc("/health", handlers.Health.Health).Methods(http.MethodGet)
	r.HandleFunc("/ready", handlers.Health.Ready).Methods(http.MethodGet)
	r.HandleFunc("/live", handlers.Health.Live).Methods(http.MethodGet)

	// Swagger UI — serves OpenAPI docs for API exploration
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	// API v1 routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Public routes
	api.HandleFunc("/auth/login", handlers.Auth.Login).Methods(http.MethodPost)

	// Protected routes
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.JWTAuth(tokenService, log))

	protected.HandleFunc("/accounts", handlers.Account.GetAccounts).Methods(http.MethodGet)
	protected.HandleFunc("/accounts/balances", handlers.Account.GetAllBalances).Methods(http.MethodGet)
	protected.HandleFunc("/accounts/{accountId}/balance", handlers.Account.GetAccountBalance).Methods(http.MethodGet)
	protected.HandleFunc("/beneficiaries", handlers.Beneficiary.GetBeneficiaries).Methods(http.MethodGet)
	protected.HandleFunc("/transfer-modes", handlers.TransferMode.GetTransferModes).Methods(http.MethodGet)

	return r
}
