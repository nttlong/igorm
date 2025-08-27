package services

import (
	"errors"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type PasswordBcrypt struct {
	cost int // bcrypt cost
}

var newAuthServiceBcryptOnce sync.Once
var authServiceBcrypt *PasswordBcrypt

// Constructor
func NewAuthServiceBcrypt(cost int) *PasswordBcrypt {
	newAuthServiceBcryptOnce.Do(func() {
		authServiceBcrypt = &PasswordBcrypt{cost: cost}
	})
	return authServiceBcrypt
}

// HashPassword tạo hash bcrypt
func (s *PasswordBcrypt) HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("empty password")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyPassword kiểm tra password với hash
func (s *PasswordBcrypt) VerifyPassword(encodedHash, password string) (bool, error) {
	if encodedHash == "" || password == "" {
		return false, errors.New("empty hash or password")
	}
	err := bcrypt.CompareHashAndPassword([]byte(encodedHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
