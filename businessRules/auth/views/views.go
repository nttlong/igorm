package views

import (
	"context"

	"dbx"
	"dynacall"
	_ "dynacall"

	"unvs.br.auth/services"
)

type ViewService struct {
	dynacall.Caller
	Tenant      string
	TenantDb    *dbx.DBXTenant
	Context     context.Context
	AccessToken string

	services.TokenService
}

func init() {
	dynacall.RegisterCaller(&ViewService{
		Caller: dynacall.Caller{
			Path: "auth",
		},
	})
}
