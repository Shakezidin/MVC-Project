package service

import (
	"context"
	"testing"

	"github.com/banking/bank-server/internal/auth"
	"github.com/banking/bank-server/internal/config"
	"github.com/banking/bank-server/internal/model"
	"github.com/banking/bank-server/internal/response"
	"github.com/google/uuid"
)

type mockUserRepo struct {
	user *model.User
	err  error
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.user, m.err
}

func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return m.user, m.err
}

func TestAuthService_Login_Success(t *testing.T) {
	hash, _ := auth.HashPassword("password123")
	user := &model.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: hash,
	}

	repo := &mockUserRepo{user: user}
	tokenSvc := auth.NewTokenService(config.JWTConfig{
		Secret: "test-secret-key-with-minimum-32-characters",
		Issuer: "test",
	})

	svc := NewAuthService(repo, tokenSvc)

	result, err := svc.Login(context.Background(), model.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if result.Token == "" {
		t.Error("expected non-empty token")
	}

	if result.TokenType != "Bearer" {
		t.Errorf("expected Bearer token type, got %s", result.TokenType)
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	hash, _ := auth.HashPassword("password123")
	user := &model.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: hash,
	}

	repo := &mockUserRepo{user: user}
	tokenSvc := auth.NewTokenService(config.JWTConfig{
		Secret: "test-secret-key-with-minimum-32-characters",
		Issuer: "test",
	})

	svc := NewAuthService(repo, tokenSvc)

	_, err := svc.Login(context.Background(), model.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Fatal("expected error for invalid password")
	}

	appErr := response.AsAppError(err)
	if appErr.Code != response.ErrCodeUnauthorized {
		t.Errorf("expected UNAUTHORIZED, got %s", appErr.Code)
	}
}
