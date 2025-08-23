package handlers

import "wx"

type LoginController struct {
}

func (c *LoginController) Login(ctx *wx.Handler, data *struct {
	Username string `json:"username"`
	Password string `json:"password"`
}) string {
	// Simulate a login process
	if data.Username == "admin" && data.Username == "password" {
		return "Login successful"
	}
	return "Login failed"
}
