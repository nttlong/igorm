package controllers

import (
	"wx"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (Auth *Auth) Login(ctx *wx.Handler, data LoginData) (any, error) {
	return data, nil
}
