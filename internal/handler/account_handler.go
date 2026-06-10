package handler

import (
	"fmt"
	"net/http"

	"github.com/banking/bank-server/internal/response"
	"github.com/banking/bank-server/internal/service"
	"github.com/banking/bank-server/internal/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// AccountHandler handles account-related HTTP endpoints.
type AccountHandler struct {
	accountService *service.AccountService
}

// NewAccountHandler creates an AccountHandler.
func NewAccountHandler(accountService *service.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// GetAccounts godoc
// @Summary      List user accounts
// @Description  Get all bank accounts of the logged-in user
// @Tags         accounts
// @Produce      json
// @Security     BearerAuth
// @Param        page  query int false "Page number"
// @Param        limit query int false "Items per page"
// @Success      200 {object} response.APIResponse
// @Failure      401 {object} response.APIResponse
// @Router       /api/v1/accounts [get]
func (h *AccountHandler) GetAccounts(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromContext(r)
	if err != nil {
		handleError(w, r, err)
		return
	}

	pagination := utils.ParsePagination(r)
	result, err := h.accountService.GetAccounts(r.Context(), userID, pagination)
	if err != nil {
		handleError(w, r, err)
		return
	}

	response.Success(w, r, http.StatusOK, "accounts fetched successfully", result)
}

// GetAllBalances godoc
// @Summary      List all account balances
// @Description  Get balances of all accounts of the logged-in user
// @Tags         accounts
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.APIResponse
// @Failure      401 {object} response.APIResponse
// @Router       /api/v1/accounts/balances [get]
func (h *AccountHandler) GetAllBalances(w http.ResponseWriter, r *http.Request) {

	userID, err := userIDFromContext(r)
	if err != nil {
		handleError(w, r, err)
		return
	}

	balances, err := h.accountService.GetAllBalances(r.Context(), userID)
	if err != nil {
		fmt.Println("error ", err)
		handleError(w, r, err)
		return
	}

	response.Success(w, r, http.StatusOK, "balances fetched successfully", balances)
}

// GetAccountBalance godoc
// @Summary      Get account balance
// @Description  Get balance of a specific account (must belong to logged-in user)
// @Tags         accounts
// @Produce      json
// @Security     BearerAuth
// @Param        accountId path string true "Account ID"
// @Success      200 {object} response.APIResponse
// @Failure      401 {object} response.APIResponse
// @Failure      404 {object} response.APIResponse
// @Router       /api/v1/accounts/{accountId}/balance [get]
func (h *AccountHandler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromContext(r)
	if err != nil {
		handleError(w, r, err)
		return
	}

	accountID, err := parseUUIDParam(r, "accountId")
	if err != nil {
		response.Error(w, r, http.StatusBadRequest,
			"invalid account ID", response.ErrCodeBadRequest, nil)
		return
	}

	balance, err := h.accountService.GetAccountBalance(r.Context(), userID, accountID)
	if err != nil {
		handleError(w, r, err)
		return
	}

	response.Success(w, r, http.StatusOK, "balance fetched successfully", balance)
}

func parseUUIDParam(r *http.Request, name string) (uuid.UUID, error) {
	vars := mux.Vars(r)
	return uuid.Parse(vars[name])
}
