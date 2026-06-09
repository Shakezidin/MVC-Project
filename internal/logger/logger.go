package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a production-ready structured Zap logger.
// Zap is chosen for its zero-allocation JSON encoding in hot paths.
func New(level string) (*zap.Logger, error) {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if os.Getenv("APP_ENV") == "development" {
		cfg.Development = true
		cfg.Encoding = "console"
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return cfg.Build()
}

// WithRequestContext enriches a logger with request-scoped fields.
func WithRequestContext(log *zap.Logger, requestID, method, path string, userID string) *zap.Logger {
	fields := []zap.Field{
		zap.String("request_id", requestID),
		zap.String("method", method),
		zap.String("path", path),
	}
	if userID != "" {
		fields = append(fields, zap.String("user_id", userID))
	}
	return log.With(fields...)
}
