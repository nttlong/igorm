package controllers

import (
	"wx"

	"xauth/services"
)

type Auth struct {
	AuthService *services.AuthService
}

func (auth *Auth) New(authSvc *wx.Global[services.AuthService]) error {
	var err error
	auth.AuthService, err = authSvc.Ins()
	if err != nil {
		return err
	}
	return nil
}
