package auth

import (
	"context"

	"dbx"
	"dynacall"
	_ "dynacall"

	"unvs.br.auth/services"
)

type RoleService struct {
	dynacall.Caller
	Tenant      string
	TenantDb    *dbx.DBXTenant
	Context     context.Context
	AccessToken string

	services.TokenService
}

func init() {
	dynacall.RegisterCaller(&RoleService{
		Caller: dynacall.Caller{
			Path: "auth",
		},
	})
}
