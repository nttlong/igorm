package jwtservice

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	AccountID int    `json:"account_id"`
	TenantID  string `json:"tenant_id"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(accountID int, tenantID string, secret string, expiry time.Duration) (string, error)
	VerifyToken(tokenString string, secret string) (*TokenClaims, error)
}
