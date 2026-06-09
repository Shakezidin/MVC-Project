package service

import (
	"context"

	"github.com/banking/bank-server/internal/cache"
	"github.com/banking/bank-server/internal/repository"
)

// HealthService checks dependency health for observability endpoints.
type HealthService struct {
	db    *repository.Database
	redis *cache.RedisClient
}

// NewHealthService creates a HealthService.
func NewHealthService(db *repository.Database, redis *cache.RedisClient) *HealthService {
	return &HealthService{db: db, redis: redis}
}

// HealthStatus represents the health check result.
type HealthStatus struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

// CheckHealth verifies all dependencies are reachable.
func (s *HealthService) CheckHealth(ctx context.Context) HealthStatus {
	services := make(map[string]string)
	overall := "healthy"

	if err := s.db.Ping(ctx); err != nil {
		services["database"] = "unhealthy"
		overall = "unhealthy"
	} else {
		services["database"] = "healthy"
	}

	if err := s.redis.Ping(ctx); err != nil {
		services["redis"] = "unhealthy"
		overall = "unhealthy"
	} else {
		services["redis"] = "healthy"
	}

	return HealthStatus{
		Status:   overall,
		Services: services,
	}
}

// IsReady returns true when all dependencies are healthy.
func (s *HealthService) IsReady(ctx context.Context) bool {
	status := s.CheckHealth(ctx)
	return status.Status == "healthy"
}
