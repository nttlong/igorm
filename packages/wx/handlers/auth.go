package handlers

import (
	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	AccessToken string `json:"access_token"`
}

type UserClaims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}
