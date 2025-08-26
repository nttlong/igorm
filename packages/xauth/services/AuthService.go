package services

import (
	"wx"
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

func (authService *AuthServiceArgon) New(dbService *wx.Depend[repo.UserRepoSQL]) (AuthService, error) {
	var err error
	authService.userRepo, err = dbService.Ins()
	if err != nil {
		return nil, err
	}
	authService.argonTime = 3           // số vòng (t)
	authService.argonTime = 3           // số vòng (t)
	authService.argonMemory = 64 * 1024 // KiB (m): 64MB
	authService.argonThreads = 2        // số luồng (p)
	authService.saltLen = 16            // bytes
	authService.keyLen = 32
	return authService, nil
}
