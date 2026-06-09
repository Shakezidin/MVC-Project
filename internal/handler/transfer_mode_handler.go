package handler

import (
	"net/http"

	"github.com/banking/bank-server/internal/response"
	"github.com/banking/bank-server/internal/service"
)

// TransferModeHandler handles transfer mode HTTP endpoints.
type TransferModeHandler struct {
	transferModeService *service.TransferModeService
}

// NewTransferModeHandler creates a TransferModeHandler.
func NewTransferModeHandler(transferModeService *service.TransferModeService) *TransferModeHandler {
	return &TransferModeHandler{transferModeService: transferModeService}
}

// GetTransferModes godoc
// @Summary      List transfer modes
// @Description  Get available transfer modes (UPI, NEFT, RTGS, IMPS)
// @Tags         transfer-modes
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.APIResponse
// @Failure      401 {object} response.APIResponse
// @Router       /api/v1/transfer-modes [get]
func (h *TransferModeHandler) GetTransferModes(w http.ResponseWriter, r *http.Request) {
	modes, err := h.transferModeService.GetTransferModes(r.Context())
	if err != nil {
		handleError(w, r, err)
		return
	}

	response.Success(w, r, http.StatusOK, "transfer modes fetched successfully", modes)
}
