package auth

import (
	"dbx"
	"dynacall"
	_ "dynacall"

	"unvs.br.auth/services"
	service "unvs.br.auth/services"
)

type User struct {
	dynacall.Caller
	Tenant   string
	TenantDb *dbx.DBXTenant
	service.FeatureService
	services.TokenService
	service.PasswordService
	service.CacheService

	AccessToken string
}

func init() {
	dynacall.RegisterCaller(&User{
		Caller: dynacall.Caller{
			Path: "auth",
		},
	})
}
