package services

import (
	"fmt"
	"strings"
	"sync"
	"xauth/models"
	"xauth/repo"
)

type Login interface {
	DoLogin(username string, password string) (*models.User, error)
	HashUserPassword(username string, password string) (string, error)
	VerifyUserPassword(username string, password string, hasPass string) (bool, error)
}
type LoginService struct {
	userRepo    repo.UserRepo
	passwordSvr PasswordService
}

func (login *LoginService) HashUserPassword(username string, password string) (string, error) {
	return login.passwordSvr.HashPassword(fmt.Sprintf("%s@%s", strings.ToLower(username), password))
}
func (login *LoginService) VerifyUserPassword(username string, password string, hasPass string) (bool, error) {
	return login.passwordSvr.VerifyPassword(hasPass, fmt.Sprintf("%s@%s", strings.ToLower(username), password))
}

var hashDefaultPassOfAdminOnce sync.Once

func (login *LoginService) DoLogin(username string, password string) (*models.User, error) {
	var err error

	hashDefaultPassOfAdminOnce.Do(func() {
		var hashDefaultPassOfAdmin string
		hashDefaultPassOfAdmin, err = login.HashUserPassword("admin", "123456")
		if err != nil {
			return
		}
		login.userRepo.CreateDefaultlUser(hashDefaultPassOfAdmin)
	})
	if err != nil {
		return nil, err
	}

	user, err := login.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	ok, err := login.VerifyUserPassword(username, password, user.Password)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}
	return user, nil

}
func NewLonginService(userRepo repo.UserRepo, passwordSvr PasswordService) Login {
	return &LoginService{
		userRepo:    userRepo,
		passwordSvr: passwordSvr,
	}
}
