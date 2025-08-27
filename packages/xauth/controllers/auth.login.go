package controllers

import (
	"wx"
	"xauth/config"
	"xauth/repo"
	"xauth/services"
)

type AuthService struct {
	wx.Service
	services.AuthService
	repo.UserRepo
	DbContext *repo.DbContext
	Config    config.ConfigService
}

func (Auth *AuthService) New() error {
	var err error
	Auth.Config, err = config.NewYamlConfigService()
	if err != nil {
		return err
	}

	Auth.Config, err = config.NewYamlConfigService()
	if err != nil {
		return err
	}

	Auth.AuthService = services.NewAuthServiceArgon()
	Auth.UserRepo = repo.NewUserRepoSQL(Auth.DbContext)
	return nil
}

//git filter-repo --path packages/wx --path-rename packages/wx/:
//git remote remove origin   # xóa remote cũ (monorepo)
//git remote add origin https://github.com/nttlong/wx.git

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (Auth *Auth) Login(ctx *wx.Handler, data LoginData, authSvc *AuthService) (any, error) {
	ret, err := authSvc.HashPassword(data.Password + "@" + data.Username)
	if err != nil {
		return nil, err
	}
	data.Password = ret
	return data, nil
}
func (Auth *Auth) Auth(ctx *wx.Handler, formData wx.Form[LoginData], authSvc *AuthService) (any, error) {
	data := formData.Data
	ret, err := authSvc.HashPassword(data.Password + "@" + data.Username)
	if err != nil {
		return nil, err
	}
	data.Password = ret
	return data, nil
}
