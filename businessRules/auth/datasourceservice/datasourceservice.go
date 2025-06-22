package datasourceservice

import (
	"dynacall"

	"dbx"

	services "unvs.br.auth/services"
)

type DataSource struct {
	dynacall.Caller
	Tenant   string
	TenantDb *dbx.DBXTenant
	services.FeatureService
	services.TokenService

	services.CacheService

	AccessToken string
}
