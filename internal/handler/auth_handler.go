package handler

import (
	"encoding/json"
	"net/http"

	"github.com/banking/bank-server/internal/model"
	"github.com/banking/bank-server/internal/response"
	"github.com/banking/bank-server/internal/service"
	"github.com/banking/bank-server/internal/validator"
)

// AuthHandler handles authentication HTTP endpoints.
type AuthHandler struct {
	authService *service.AuthService
	validator   *validator.Validator
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(authService *service.AuthService, v *validator.Validator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   v,
	}
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.LoginRequest true "Login credentials"
// @Success      200 {object} response.APIResponse
// @Failure      401 {object} response.APIResponse
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, r, http.StatusBadRequest,
			"invalid request body", response.ErrCodeBadRequest, nil)
		return
	}

	if details := h.validator.ValidateStruct(req); details != nil {
		response.Error(w, r, http.StatusUnprocessableEntity,
			"validation failed", response.ErrCodeValidation, details)
		return
	}

	result, err := h.authService.Login(r.Context(), req)
	if err != nil {
		handleError(w, r, err)
		return
	}

	response.Success(w, r, http.StatusOK, "login successful", result)
}
