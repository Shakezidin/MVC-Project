package observability

import (
	"context"

	"github.com/banking/bank-server/internal/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(
	serviceName string,
	projectID string,
	topicID string,
	logLevel string,
) (*zap.Logger, error) {
	config := zap.NewProductionConfig()

	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		level = zapcore.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(level)

	baseLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	publisher, err := NewPubSubPublisher(
		projectID,
		topicID,
	)

	if err != nil {
		return baseLogger, nil
	}

	pubsubCore := NewPubSubCore(
		zapcore.NewNopCore(),
		publisher,
		serviceName,
	)

	core := zapcore.NewTee(
		baseLogger.Core(),
		pubsubCore,
	)

	logger := zap.New(core)

	return logger, nil
}

// FromContext returns a logger with fields from the context (request ID, user ID)
func FromContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	fields := []zap.Field{}

	if requestID := utils.RequestIDFromContext(ctx); requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}

	if userID := utils.UserIDFromContext(ctx); userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}

	return logger.With(fields...)
}

