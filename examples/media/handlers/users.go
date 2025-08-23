package handlers

import "wx"

type Users struct {
}
type UserData struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func (c *Users) CreatetUser(ctx *wx.Handler, data UserData) (string, error) {
	// Simulate fetching user data
	r, err := HashPassword(data.Password)
	if err != nil {
		return "", err
	}
	return r, nil
}
