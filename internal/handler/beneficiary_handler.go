package handler

import (
	"net/http"

	"github.com/banking/bank-server/internal/response"
	"github.com/banking/bank-server/internal/service"
	"github.com/banking/bank-server/internal/utils"
)

// BeneficiaryHandler handles beneficiary HTTP endpoints.
type BeneficiaryHandler struct {
	beneficiaryService *service.BeneficiaryService
}

// NewBeneficiaryHandler creates a BeneficiaryHandler.
func NewBeneficiaryHandler(beneficiaryService *service.BeneficiaryService) *BeneficiaryHandler {
	return &BeneficiaryHandler{beneficiaryService: beneficiaryService}
}

// GetBeneficiaries godoc
// @Summary      List beneficiaries
// @Description  Get beneficiary list of the logged-in user
// @Tags         beneficiaries
// @Produce      json
// @Security     BearerAuth
// @Param        page  query int false "Page number"
// @Param        limit query int false "Items per page"
// @Success      200 {object} response.APIResponse
// @Failure      401 {object} response.APIResponse
// @Router       /api/v1/beneficiaries [get]
func (h *BeneficiaryHandler) GetBeneficiaries(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromContext(r)
	if err != nil {
		handleError(w, r, err)
		return
	}

	pagination := utils.ParsePagination(r)
	result, err := h.beneficiaryService.GetBeneficiaries(r.Context(), userID, pagination)
	if err != nil {
		handleError(w, r, err)
		return
	}

	response.Success(w, r, http.StatusOK, "beneficiaries fetched successfully", result)
}
