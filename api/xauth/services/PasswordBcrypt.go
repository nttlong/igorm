package services

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type PasswordBcrypt struct {
	cost int // bcrypt cost
}

// Constructor
func NewAuthServiceBcrypt(cost int) *PasswordBcrypt {
	if cost == 0 {
		cost = bcrypt.DefaultCost // default = 10
	}
	return &PasswordBcrypt{cost: cost}
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
