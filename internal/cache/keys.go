package cache

import (
	"context"
	"fmt"
)

// Cache key prefixes — centralized to avoid key collision and enable bulk invalidation.
const (
	prefixAccountList    = "accounts:user:"
	prefixBeneficiaryList = "beneficiaries:user:"
	keyTransferModes     = "transfer_modes:all"
)

// AccountListKey returns the cache key for a user's account list.
func AccountListKey(userID string) string {
	return prefixAccountList + userID
}

// BeneficiaryListKey returns the cache key for a user's beneficiary list.
func BeneficiaryListKey(userID string) string {
	return prefixBeneficiaryList + userID
}

// TransferModesKey returns the cache key for transfer modes.
func TransferModesKey() string {
	return keyTransferModes
}

// InvalidateUserCache removes all cached data for a user.
func InvalidateUserCache(ctx context.Context, client *RedisClient, userID string) error {
	keys := []string{
		AccountListKey(userID),
		BeneficiaryListKey(userID),
	}
	return client.Client.Del(ctx, keys...).Err()
}

// InvalidateTransferModes removes cached transfer modes.
func InvalidateTransferModes(ctx context.Context, client *RedisClient) error {
	return client.Client.Del(ctx, TransferModesKey()).Err()
}

// InvalidateAccountList removes cached account list for a user.
func InvalidateAccountList(ctx context.Context, client *RedisClient, userID string) error {
	return client.Client.Del(ctx, AccountListKey(userID)).Err()
}

// InvalidateBeneficiaryList removes cached beneficiary list for a user.
func InvalidateBeneficiaryList(ctx context.Context, client *RedisClient, userID string) error {
	return client.Client.Del(ctx, BeneficiaryListKey(userID)).Err()
}

// Pattern for debugging — not used in production paths.
func userCachePattern(userID string) string {
	return fmt.Sprintf("*:user:%s", userID)
}
