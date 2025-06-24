package unvscore

import (
	"context"
	"dbx"
	"sync"

	config "unvs.core/config"

	cacher "unvs.core/cacher"
	services "unvs.core/services"
)

type factory struct {
}

// type TokenServiceParams struct {
// 	Cache         cacher.Cache
// 	Context       context.Context
// 	TenantDb      *dbx.DBXTenant
// 	EncryptionKey string
// 	Language      string
// }

//	type PasswordServiceParams struct {
//		// Cache   cacher.Cache
//		Context context.Context
//	}
// type AuthServiceParams struct {
// 	Cache         cacher.Cache
// 	Context       context.Context
// 	EncryptionKey string
// 	Language      string

// 	TenantDb *dbx.DBXTenant
// }

func (f *factory) GetTokenService(ctx context.Context, Language string, tenantName string) (*services.TokenService, error) {
	db := Config.GetDb()
	db.Open()
	defer db.Close()
	tenantDb, err := db.GetTenant(tenantName)
	if err != nil {
		return nil, err
	}
	tenantDb.Open()
	return &services.TokenService{
		Cache:         Config.GetCache(),
		Context:       ctx,
		TenantDb:      tenantDb,
		EncryptionKey: Config.AppConfig.EncryptionKey,
		Language:      Language,
	}, nil
}

func (f *factory) GetPasswordService(ctx context.Context) *services.PasswordService {

	return &services.PasswordService{
		Cache:   Config.GetCache(),
		Context: ctx,
	}
}
func (f *factory) GetAuthService(ctx context.Context, Language string, tenantName string) (*services.AuthenticateService, error) {
	db := Config.GetDb()
	db.Open()
	defer db.Close()
	tenantDb, err := db.GetTenant(tenantName)
	if err != nil {
		return nil, err
	}
	retSvc := &services.AuthenticateService{
		TokenService: services.TokenService{
			Cache:         Config.GetCache(),
			Context:       ctx,
			TenantDb:      tenantDb,
			EncryptionKey: Config.AppConfig.EncryptionKey,
			Language:      Language,
		},
		PasswordService: services.PasswordService{
			Cache:   Config.GetCache(),
			Context: ctx,
		},
		Context: ctx,
		Cache:   Config.GetCache(),
	}
	return retSvc, nil
}

var Factory = &factory{}

type _config struct {
	AppConfig *config.Config
}

func (c *_config) GetCache() cacher.Cache {
	if c.AppConfig == nil {
		panic("config not loaded yet, call Config.LoadConfig first")
	}
	return config.GetCache()
}
func (c *_config) GetDb() *dbx.DBX {
	if c.AppConfig == nil {
		panic("config not loaded yet, call Config.LoadConfig first")
	}
	return config.CreateDbx()
}
func (c *_config) LoadConfig(path string) {
	config.LoadConfig(path)
	c.AppConfig = config.AppConfigInstance
}

var Config = &_config{}

type serviceBuilderInfo[T any] struct {
	tenantDb      *dbx.DBXTenant
	cache         cacher.Cache
	Svc           *T
	err           error
	EncryptionKey string
}
type serviceInfoWithLanguage[T any] struct {
	owner    *serviceBuilderInfo[T]
	Language string
}

var serviceBuilderCache sync.Map

func ServiceBuilder[T any](tenantName string) *serviceBuilderInfo[T] {
	if v, ok := serviceBuilderCache.Load(tenantName); ok {
		return v.(*serviceBuilderInfo[T])
	}
	ret := &serviceBuilderInfo[T]{}
	db := Config.GetDb()
	db.Open()
	defer db.Close()
	tenantDb, err := db.GetTenant(tenantName)
	if err != nil {
		ret.err = err
		return ret
	}
	ret.tenantDb = tenantDb
	ret.cache = Config.GetCache()
	ret.EncryptionKey = Config.AppConfig.EncryptionKey
	return ret

}
func (s *serviceBuilderInfo[T]) WithLan(language string) *serviceInfoWithLanguage[T] {

	return &serviceInfoWithLanguage[T]{
		owner:    s,
		Language: language,
	}

}
