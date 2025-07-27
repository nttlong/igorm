package jwtservice

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type DefaultJWTService struct{}

func NewJWTService() *DefaultJWTService {
	return &DefaultJWTService{}
}

func (s *DefaultJWTService) GenerateToken(accountID int, tenantID string, secret string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		AccountID: accountID,
		TenantID:  tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *DefaultJWTService) VerifyToken(tokenString string, secret string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok {
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}
