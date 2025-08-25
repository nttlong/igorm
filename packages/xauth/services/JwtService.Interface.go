package services

import "time"

type JwtConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}
type IJWTService interface {
	Generate(cfg JwtConfig, data any) (*JWTPair, error)
	GetConfig() (*JwtConfig, error)
}
