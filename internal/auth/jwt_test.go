package auth

import (
	"testing"
	"time"

	"github.com/banking/bank-server/internal/config"
	"github.com/google/uuid"
)

func TestGenerateAndValidateToken(t *testing.T) {
	cfg := config.JWTConfig{
		Secret: "test-secret-key-with-minimum-32-characters",
		Expiry: 1 * time.Hour,
		Issuer: "test-issuer",
	}

	svc := NewTokenService(cfg)
	userID := uuid.New()
	email := "test@example.com"

	token, expiresIn, err := svc.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("token should not be empty")
	}

	if expiresIn != 3600 {
		t.Errorf("expected expiresIn 3600, got %d", expiresIn)
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	cfg := config.JWTConfig{
		Secret: "test-secret-key-with-minimum-32-characters",
		Expiry: 1 * time.Hour,
		Issuer: "test-issuer",
	}

	svc := NewTokenService(cfg)

	_, err := svc.ValidateToken("invalid.token.here")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}
