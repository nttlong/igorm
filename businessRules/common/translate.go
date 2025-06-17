package common

import (
	"context"

	"dbx"
	"dynacall"
	_ "dynacall"
)

type TranslateService struct {
	dynacall.Caller
	Tenant      string
	TenantDb    *dbx.DBXTenant
	Context     context.Context
	AccessToken string
}

func init() {
	dynacall.RegisterCaller(&TranslateService{})
}
