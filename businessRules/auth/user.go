package auth

import (
	"dbx"
	"dynacall"
	_ "dynacall"
	"fmt"
	"time"
)

type User struct {
	dynacall.Caller
	Tenant   string
	TenantDb *dbx.DBXTenant
}

func (u *User) Login(username string, password string, loginOn time.Time) bool {
	fmt.Println(loginOn.Year())
	return username == "admin" && password == "password"
}

func init() {
	dynacall.RegisterCaller(&User{
		Caller: dynacall.Caller{
			Path: "auth",
		},
	})
}
