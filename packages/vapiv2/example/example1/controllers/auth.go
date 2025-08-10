package controllers

import (
	"fmt"
	"vapi"
)

type Auth struct {
}

func (a *Auth) Oauth(ctx *struct {
	vapi.Handler `route:"uri:/api/@/token"`
}, data *struct {
	UserName string
	Password string
}) (interface{}, error) {
	fmt.Println(data.Password)
	fmt.Println(data.UserName)

	return data, nil
}
func init() {
	vapi.Controller("", "", func() (*Auth, error) {
		return &Auth{}, nil
	})

}
