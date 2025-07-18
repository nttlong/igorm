package unvsauth

import (
	"errors"
	di "vdi"

	"github.com/google/uuid"
)

type AuthService struct {
	UserRepo    di.Singleton[*AuthService, *UserRepository]
	JwtProvider di.Singleton[*AuthService, *JwtProvider]
	Hasher      di.Singleton[*AuthService, Hasher] // ← dùng interface
}

func (svc *AuthService) Register(email, username, password string) (*User, error) {

	hash, err := svc.Hasher.Get().Hash(password)
	if err != nil {
		return nil, err
	}
	user := &User{
		UserId:       uuid.NewString(),
		Email:        email,
		Username:     username,
		HashPassword: hash,
		IsActive:     true,
	}
	err = svc.UserRepo.Get().Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *AuthService) Login(identifier, password string) (string, error) {
	user, err := svc.UserRepo.Get().FindByEmailOrUsername(identifier)
	if err != nil || !user.IsActive {
		return "", errors.New("invalid credentials")
	}
	if !svc.Hasher.Get().Verify(user.HashPassword, password) {
		return "", errors.New("invalid credentials")
	}
	return svc.JwtProvider.Get().GenerateToken(user.UserId)

}
