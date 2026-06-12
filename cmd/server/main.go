// @title           Bank Server API
// @version         1.0
// @description     Production-grade banking backend service API
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/banking/bank-server/docs"
	"github.com/banking/bank-server/internal/auth"
	"github.com/banking/bank-server/internal/cache"
	"github.com/banking/bank-server/internal/config"
	"github.com/banking/bank-server/internal/handler"
	"github.com/banking/bank-server/internal/observability"
	"github.com/banking/bank-server/internal/repository"
	"github.com/banking/bank-server/internal/router"
	"github.com/banking/bank-server/internal/service"
	"github.com/banking/bank-server/internal/validator"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	log, err := observability.NewLogger(
		"bank-server",
		cfg.PubSub.ProjectID,
		cfg.PubSub.TopicID,
		cfg.Log.Level,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error creating logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("starting bank-server", zap.String("env", cfg.Server.Environment), zap.String("config", cfg.String()))

	ctx := context.Background()

	db, err := repository.NewDatabase(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	redisClient, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		return fmt.Errorf("connect redis: %w", err)
	}
	defer redisClient.Close()

	cacheStore := cache.New(redisClient)

	tokenService := auth.NewTokenService(cfg.JWT)

	userRepo := repository.NewUserRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	balanceRepo := repository.NewBalanceRepository(db)
	beneficiaryRepo := repository.NewBeneficiaryRepository(db)
	transferModeRepo := repository.NewTransferModeRepository(db)

	authService := service.NewAuthService(userRepo, tokenService)
	accountService := service.NewAccountService(accountRepo, balanceRepo, cacheStore, cfg.Cache, log)
	beneficiaryService := service.NewBeneficiaryService(beneficiaryRepo, cacheStore, cfg.Cache, log)
	transferModeService := service.NewTransferModeService(transferModeRepo, cacheStore, cfg.Cache, log)
	healthService := service.NewHealthService(db, redisClient)

	v := validator.New()
	handlers := router.Handlers{
		Auth:         handler.NewAuthHandler(authService, v),
		Account:      handler.NewAccountHandler(accountService),
		Beneficiary:  handler.NewBeneficiaryHandler(beneficiaryService),
		TransferMode: handler.NewTransferModeHandler(transferModeService),
		Health:       handler.NewHealthHandler(healthService),
	}

	_ = docs.SwaggerInfo

	r := router.New(cfg, handlers, tokenService, log)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("server listening", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		log.Info("shutdown signal received", zap.String("signal", sig.String()))
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	log.Info("server stopped gracefully")
	return nil
}
