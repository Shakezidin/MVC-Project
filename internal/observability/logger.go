package observability

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// func NewLogger(
// 	serviceName string,
// 	projectID string,
// 	topicID string,
// ) (*zap.Logger, error) {

// 	config := zap.NewProductionConfig()

// 	baseLogger, err := config.Build()
// 	if err != nil {
// 		return nil, err
// 	}

// 	publisher, err := NewPubSubPublisher(
// 		projectID,
// 		topicID,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	pubsubCore := NewPubSubCore(
// 		baseLogger.Core(),
// 		publisher,
// 		serviceName,
// 	)

// 	logger := zap.New(pubsubCore)

// 	return logger, nil
// }

func NewLogger(
	serviceName string,
	projectID string,
	topicID string,
) (*zap.Logger, error) {

	/*
		========================================
		BASE LOGGER
		========================================
	*/
	config := zap.NewProductionConfig()

	baseLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	/*
		========================================
		PUBSUB PUBLISHER
		========================================
	*/
	publisher, err := NewPubSubPublisher(
		projectID,
		topicID,
	)

	if err != nil {
		return nil, err
	}

	/*
		========================================
		PUBSUB CORE
		========================================
	*/
	pubsubCore := NewPubSubCore(
		zapcore.NewNopCore(),
		publisher,
		serviceName,
	)

	/*
		========================================
		TEE BOTH CORES
		========================================
	*/
	core := zapcore.NewTee(
		baseLogger.Core(),
		pubsubCore,
	)

	logger := zap.New(core)

	return logger, nil
}
