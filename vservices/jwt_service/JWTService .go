package jwt_service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserId   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userId string, tenantID string, secret string, expiry time.Duration) (string, error)
	VerifyToken(tokenString string, secret string) (*TokenClaims, error)
}
