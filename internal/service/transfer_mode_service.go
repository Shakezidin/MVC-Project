package service

import (
	"context"

	"github.com/banking/bank-server/internal/cache"
	"github.com/banking/bank-server/internal/config"
	"github.com/banking/bank-server/internal/model"
	"github.com/banking/bank-server/internal/repository"
	"go.uber.org/zap"
)

// TransferModeService handles transfer mode business logic.
type TransferModeService struct {
	transferModeRepo repository.TransferModeRepository
	cache            *cache.Cache
	cacheTTL         config.CacheConfig
	log              *zap.Logger
}

// NewTransferModeService creates a TransferModeService with injected dependencies.
func NewTransferModeService(
	transferModeRepo repository.TransferModeRepository,
	cacheClient *cache.Cache,
	cacheTTL config.CacheConfig,
	log *zap.Logger,
) *TransferModeService {
	return &TransferModeService{
		transferModeRepo: transferModeRepo,
		cache:            cacheClient,
		cacheTTL:         cacheTTL,
		log:              log,
	}
}

// GetTransferModes returns all active transfer modes with Redis caching.
func (s *TransferModeService) GetTransferModes(ctx context.Context) ([]model.TransferModeResponse, error) {
	cacheKey := cache.TransferModesKey()

	var cached []model.TransferModeResponse
	found, err := s.cache.Get(ctx, cacheKey, &cached)
	if err != nil {
		s.log.Warn("cache read failed, falling back to DB", zap.Error(err))
	}

	if found {
		return cached, nil
	}

	modes, err := s.transferModeRepo.GetAllActive(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]model.TransferModeResponse, 0, len(modes))
	for _, m := range modes {
		responses = append(responses, model.TransferModeResponse{
			Code: m.Code,
			Name: m.Name,
		})
	}

	if err := s.cache.Set(ctx, cacheKey, responses, s.cacheTTL.TransferModesTTL); err != nil {
		s.log.Warn("cache write failed", zap.Error(err))
	}

	return responses, nil
}
