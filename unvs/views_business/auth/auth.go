package auth

import (
	"unvs/views"
)

type Auth struct {
	views.BaseView
}
type AuthParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *Auth) Login(authLogin AuthParam) AuthParam {
	return authLogin
	// Code to handle login
}
func init() {
	views.AddView(&Auth{
		BaseView: views.BaseView{
			ViewPath: "auth",
			//IsAuth:   true,
		},
	})
}
