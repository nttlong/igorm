package services

import (
	"caching"
	"context"
	"crypto/rand"
	"dbx"
	"encoding/base64"
	"time"
)

type RefreshTokenService struct {
	Size          int
	Cache         caching.Cache
	TenantDb      *dbx.DBXTenant
	EncryptionKey string
	Context       context.Context
}

func (rf *RefreshTokenService) GenerateRefreshToken() (string, error) {
	if rf.Size == 0 {
		rf.Size = 32
	}
	b := make([]byte, rf.Size) // 256-bit
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	ret := base64.URLEncoding.EncodeToString(b)

	key := "refresh_token_" + rf.TenantDb.TenantDbName + ":" + ret
	rf.Cache.Set(rf.Context, key, ret, time.Minute*5)

	return base64.URLEncoding.EncodeToString(b), nil
}
