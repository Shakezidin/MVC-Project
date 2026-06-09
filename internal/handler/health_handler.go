package handler

import (
	"net/http"

	"github.com/banking/bank-server/internal/response"
	"github.com/banking/bank-server/internal/service"
)

// HealthHandler handles health check endpoints.
type HealthHandler struct {
	healthService *service.HealthService
}

// NewHealthHandler creates a HealthHandler.
func NewHealthHandler(healthService *service.HealthService) *HealthHandler {
	return &HealthHandler{healthService: healthService}
}

// Health returns overall service health including dependencies.
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	status := h.healthService.CheckHealth(r.Context())
	httpStatus := http.StatusOK
	if status.Status != "healthy" {
		httpStatus = http.StatusServiceUnavailable
	}
	response.Success(w, r, httpStatus, "health check completed", status)
}

// Ready returns readiness probe — used by orchestrators before routing traffic.
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	if h.healthService.IsReady(r.Context()) {
		response.Success(w, r, http.StatusOK, "service is ready", map[string]string{"status": "ready"})
		return
	}
	response.Error(w, r, http.StatusServiceUnavailable,
		"service is not ready", response.ErrCodeInternal, nil)
}

// Live returns liveness probe — confirms the process is running.
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	response.Success(w, r, http.StatusOK, "service is alive", map[string]string{"status": "alive"})
}
