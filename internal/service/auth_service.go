package service

import (
	"context"

	"github.com/banking/bank-server/internal/auth"
	"github.com/banking/bank-server/internal/model"
	"github.com/banking/bank-server/internal/repository"
	"github.com/banking/bank-server/internal/response"
)

// AuthService handles authentication business logic.
type AuthService struct {
	userRepo     repository.UserRepository
	tokenService *auth.TokenService
}

// NewAuthService creates an AuthService with injected dependencies.
func NewAuthService(userRepo repository.UserRepository, tokenService *auth.TokenService) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Login authenticates a user and returns a JWT token.
func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		return nil, response.NewUnauthorizedError("invalid email or password")
	}

	token, expiresIn, err := s.tokenService.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, response.NewInternalError(err)
	}

	return &model.LoginResponse{
		Token:     token,
		ExpiresIn: expiresIn,
		TokenType: "Bearer",
	}, nil
}
