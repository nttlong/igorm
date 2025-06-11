package auth

import (
	"context"

	"dbx"
	"dynacall"
	_ "dynacall"
)

type User struct {
	dynacall.Caller
	Tenant    string
	TenantDb  *dbx.DBXTenant
	Context   context.Context
	JwtSecret []byte
}

func init() {
	dynacall.RegisterCaller(&User{
		Caller: dynacall.Caller{
			Path: "auth",
		},
	})
}
