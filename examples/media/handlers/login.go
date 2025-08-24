package handlers

import (
	"wx"
)

type Logins struct {
}
type UserLoginInfo struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
}

func (c *Logins) Login(ctx *wx.Handler, data *struct {
	Username string `json:"username"`
	Password string `json:"password"`
}, loginSvc *wx.Depend[LoginService]) (string, error) {
	//string {
	// Simulate a login process
	if data.Username == "admin" && data.Password == "password" {
		loginSvcIns, err := loginSvc.Ins()
		if err != nil {
			return "", err
		}
		ret, err := loginSvcIns.GenerateJWT(data, nil)
		if err != nil {
			return "", err
		}

		return ret, nil
		//return "Login successful"
	}
	return "Login failed", nil
}
