package service

import (
	"context"
	"math"

	"github.com/banking/bank-server/internal/cache"
	"github.com/banking/bank-server/internal/config"
	"github.com/banking/bank-server/internal/model"
	"github.com/banking/bank-server/internal/repository"
	"github.com/banking/bank-server/internal/response"
	"github.com/banking/bank-server/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// BeneficiaryService handles beneficiary business logic.
type BeneficiaryService struct {
	beneficiaryRepo repository.BeneficiaryRepository
	cache           *cache.Cache
	cacheTTL        config.CacheConfig
	log             *zap.Logger
}

// NewBeneficiaryService creates a BeneficiaryService with injected dependencies.
func NewBeneficiaryService(
	beneficiaryRepo repository.BeneficiaryRepository,
	cacheClient *cache.Cache,
	cacheTTL config.CacheConfig,
	log *zap.Logger,
) *BeneficiaryService {
	return &BeneficiaryService{
		beneficiaryRepo: beneficiaryRepo,
		cache:           cacheClient,
		cacheTTL:        cacheTTL,
		log:             log,
	}
}

type cachedBeneficiaryList struct {
	Beneficiaries []model.BeneficiaryResponse `json:"beneficiaries"`
	TotalCount    int                         `json:"total_count"`
}

// GetBeneficiaries returns beneficiaries for a user with Redis caching.
func (s *BeneficiaryService) GetBeneficiaries(ctx context.Context, userID uuid.UUID, pagination utils.Pagination) (*response.PaginatedData, error) {
	cacheKey := cache.BeneficiaryListKey(userID.String())

	var cached cachedBeneficiaryList
	found, err := s.cache.Get(ctx, cacheKey, &cached)
	if err != nil {
		s.log.Warn("cache read failed, falling back to DB", zap.Error(err))
	}

	if found {
		return &response.PaginatedData{
			Items:      cached.Beneficiaries,
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			TotalCount: cached.TotalCount,
			TotalPages: int(math.Ceil(float64(cached.TotalCount) / float64(pagination.Limit))),
		}, nil
	}

	beneficiaries, total, err := s.beneficiaryRepo.GetByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, err
	}

	responses := make([]model.BeneficiaryResponse, 0, len(beneficiaries))
	for _, b := range beneficiaries {
		responses = append(responses, toBeneficiaryResponse(b))
	}

	if err := s.cache.Set(ctx, cacheKey, cachedBeneficiaryList{
		Beneficiaries: responses,
		TotalCount:    total,
	}, s.cacheTTL.BeneficiaryListTTL); err != nil {
		s.log.Warn("cache write failed", zap.Error(err))
	}

	return &response.PaginatedData{
		Items:      responses,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalCount: total,
		TotalPages: int(math.Ceil(float64(total) / float64(pagination.Limit))),
	}, nil
}

func toBeneficiaryResponse(b model.Beneficiary) model.BeneficiaryResponse {
	return model.BeneficiaryResponse{
		BeneficiaryName:     b.BeneficiaryName,
		BankName:            b.BankName,
		AccountNumberMasked: utils.MaskAccountNumber(b.AccountNumber),
		IFSC:                b.IFSC,
		Nickname:            b.Nickname,
	}
}
