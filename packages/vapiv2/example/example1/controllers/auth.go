package controllers

import (
	"fmt"
	"vapi"
)

type Auth struct {
}

func (a *Auth) Oauth(ctx *struct {
	vapi.Handler `route:"uri:/api/@/token"`
}, data struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}) (*struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}, error) {
	fmt.Println(data.Password)
	fmt.Println(data.UserName)
	ret := struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}{}
	ret.AccessToken = "12345556"
	ret.TokenType = "Bearer"
	ret.ExpiresIn = 3600
	return &ret, nil
}
func init() {
	vapi.Controller(func() (*Auth, error) {
		return &Auth{}, nil
	})

}
