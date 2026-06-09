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

// AccountService handles account and balance business logic.
type AccountService struct {
	accountRepo repository.AccountRepository
	balanceRepo repository.BalanceRepository
	cache       *cache.Cache
	cacheTTL    config.CacheConfig
	log         *zap.Logger
}

// NewAccountService creates an AccountService with injected dependencies.
func NewAccountService(
	accountRepo repository.AccountRepository,
	balanceRepo repository.BalanceRepository,
	cacheClient *cache.Cache,
	cacheTTL config.CacheConfig,
	log *zap.Logger,
) *AccountService {
	return &AccountService{
		accountRepo: accountRepo,
		balanceRepo: balanceRepo,
		cache:       cacheClient,
		cacheTTL:    cacheTTL,
		log:         log,
	}
}

type cachedAccountList struct {
	Accounts   []model.AccountResponse `json:"accounts"`
	TotalCount int                     `json:"total_count"`
}

// GetAccounts returns all accounts for a user with Redis caching.
func (s *AccountService) GetAccounts(ctx context.Context, userID uuid.UUID, pagination utils.Pagination) (*response.PaginatedData, error) {
	cacheKey := cache.AccountListKey(userID.String())

	var cached cachedAccountList
	found, err := s.cache.Get(ctx, cacheKey, &cached)
	if err != nil {
		s.log.Warn("cache read failed, falling back to DB", zap.Error(err))
	}

	if found {
		return &response.PaginatedData{
			Items:      cached.Accounts,
			Page:       pagination.Page,
			Limit:      pagination.Limit,
			TotalCount: cached.TotalCount,
			TotalPages: int(math.Ceil(float64(cached.TotalCount) / float64(pagination.Limit))),
		}, nil
	}

	accounts, total, err := s.accountRepo.GetByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, err
	}

	responses := make([]model.AccountResponse, 0, len(accounts))
	for _, acc := range accounts {
		responses = append(responses, toAccountResponse(acc))
	}

	// Cache the full list for the user (pagination applied at read time from DB;
	// for simplicity we cache the current page result).
	if err := s.cache.Set(ctx, cacheKey, cachedAccountList{
		Accounts:   responses,
		TotalCount: total,
	}, s.cacheTTL.AccountListTTL); err != nil {
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

// GetAllBalances returns balances for all accounts of a user.
func (s *AccountService) GetAllBalances(ctx context.Context, userID uuid.UUID) ([]model.BalanceResponse, error) {
	balances, err := s.balanceRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.BalanceResponse, 0, len(balances))
	for _, bal := range balances {
		responses = append(responses, toBalanceResponse(bal))
	}
	return responses, nil
}

// GetAccountBalance returns the balance for a specific account, validating ownership.
func (s *AccountService) GetAccountBalance(ctx context.Context, userID, accountID uuid.UUID) (*model.BalanceResponse, error) {
	// Ownership check: account must belong to the authenticated user.
	_, err := s.accountRepo.GetByIDAndUserID(ctx, accountID, userID)
	if err != nil {
		return nil, err
	}

	balance, err := s.balanceRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	resp := toBalanceResponse(*balance)
	return &resp, nil
}

func toAccountResponse(acc model.BankAccount) model.AccountResponse {
	return model.AccountResponse{
		AccountID:           acc.ID,
		AccountType:         acc.AccountType,
		BranchName:          acc.BranchName,
		MaskedAccountNumber: utils.MaskAccountNumber(acc.AccountNumber),
		Status:              acc.Status,
	}
}

func toBalanceResponse(bal model.Balance) model.BalanceResponse {
	return model.BalanceResponse{
		AccountID:        bal.AccountID,
		AvailableBalance: bal.AvailableBalance,
		CurrentBalance:   bal.CurrentBalance,
		Currency:         bal.Currency,
	}
}
