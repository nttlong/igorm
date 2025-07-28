package passwordservice

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordService interface {
	Hash(password string) (string, error)
	Verify(password, hashed string) bool
}

type BcryptPasswordService struct {
	cost int // default 12
}

func NewBcryptPasswordService(cost int) *BcryptPasswordService {
	if cost <= 0 {
		cost = 12
	}
	return &BcryptPasswordService{cost: cost}
}

func (s *BcryptPasswordService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	return string(bytes), err
}

func (s *BcryptPasswordService) Verify(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
