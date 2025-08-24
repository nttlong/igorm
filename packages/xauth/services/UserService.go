package services

import (
	"vdb"
	"wx"
	dbModels "xauth/dbModels"
)

type UserService struct {
	authService *AuthService
}

func (userService *UserService) New(
	authService *wx.Global[AuthService],
) error {
	var err error
	userService.authService, err = authService.Ins()
	if err != nil {
		return err
	}

	return nil
}
func (userService *UserService) CreateUser(db *vdb.TenantDB, user *dbModels.Users, password string) error {
	var err error
	if password != "" {
		user.HashedPassword, err = userService.authService.HashPassword(user.Username + "@" + password)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	err = db.Create(user)
	return err

}
