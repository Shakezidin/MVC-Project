package auth

import (
	"fmt"
	"time"

	"github.com/banking/bank-server/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims holds JWT payload for authenticated users.
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// TokenService handles JWT generation and validation.
type TokenService struct {
	secret []byte
	expiry time.Duration
	issuer string
}

// NewTokenService creates a JWT token service from config.
func NewTokenService(cfg config.JWTConfig) *TokenService {
	return &TokenService{
		secret: []byte(cfg.Secret),
		expiry: cfg.Expiry,
		issuer: cfg.Issuer,
	}
}

// GenerateToken creates a signed JWT for the given user.
func (s *TokenService) GenerateToken(userID uuid.UUID, email string) (string, int64, error) {
	now := time.Now()
	expiresAt := now.Add(s.expiry)

	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", 0, fmt.Errorf("sign token: %w", err)
	}

	return signed, int64(s.expiry.Seconds()), nil
}

// ValidateToken parses and validates a JWT, returning claims on success.
func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
