package controllers

import (
	"wx"

	"xauth/services"
)

type Auth struct {
	AuthService services.AuthService
}

func (auth *Auth) New(authSvc *wx.Global[services.AuthServiceArgon]) error {
	var err error
	AuthService, err := authSvc.Ins()
	if err != nil {
		return err
	}
	auth.AuthService = &AuthService
	if err != nil {
		return err
	}
	return nil
}
