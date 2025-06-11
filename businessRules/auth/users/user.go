package auth

import (
	"caching"
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
	Cache     caching.Cache
}

func init() {
	dynacall.RegisterCaller(&User{
		Caller: dynacall.Caller{
			Path: "auth",
		},
	})
}
