package unvsauth

import (
	"errors"
	models "unvs-auth/models"
	"vdb"
	di "vdi"

	"github.com/google/uuid"
)

type AuthService struct {
	TenantDb di.Singleton[*AuthService, *vdb.TenantDB]
	UserRepo di.Singleton[*AuthService, UserRepository]

	JwtProvider di.Singleton[*AuthService, *JwtProvider]
	Hasher      di.Singleton[*AuthService, Hasher] // ← dùng interface
}

func (svc *AuthService) Register(email, username, password string) (*models.User, error) {

	hash, err := svc.Hasher.Get().Hash(password)
	if err != nil {
		return nil, err
	}
	user := models.User{
		UserId:       uuid.NewString(),
		Email:        email,
		Username:     username,
		HashPassword: hash,
		IsActive:     true,
	}
	err = svc.UserRepo.Get().Create(&user) //<-- neu muon su dung cache phai sua cho nay uh?
	if err != nil {
		return nil, err
	}
	return &user, nil
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
