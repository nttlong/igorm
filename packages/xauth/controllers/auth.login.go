package controllers

import (
	"wx"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (Auth *Auth) Login(ctx *wx.Handler, data LoginData) (any, error) {
	ret, err := Auth.AuthService.HashPassword(data.Password + "@" + data.Username)
	if err != nil {
		return nil, err
	}
	data.Password = ret
	return data, nil
}
func (Auth *Auth) Auth(ctx *wx.Handler, formData wx.Form[LoginData]) (any, error) {
	data := formData.Data
	ret, err := Auth.AuthService.HashPassword(data.Password + "@" + data.Username)
	if err != nil {
		return nil, err
	}
	data.Password = ret
	return data, nil
}
