package services

import (
	"xauth/repo"
)

type AuthService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(encodedHash, password string) (bool, error)
}
type AuthServiceArgon struct {
	argonTime    uint32 // số vòng (t)
	argonMemory  uint32 // KiB (m): 64MB
	argonThreads uint8  // số luồng (p)
	saltLen      int    // bytes
	keyLen       uint32 // bytes (256-bit)
	userRepo     repo.UserRepo
}

func NewAuthServiceArgon() *AuthServiceArgon {
	return &AuthServiceArgon{
		argonTime:    3,
		argonMemory:  64 * 1024,
		argonThreads: 2,
		saltLen:      16,
		keyLen:       32,
	}
}
