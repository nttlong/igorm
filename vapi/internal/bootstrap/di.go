package bootstrap

import (
	"vapi/internal/cache"
	"vapi/internal/config"
	ctxSvc "vapi/internal/context_service"
	"vapi/internal/dbcontext"
	"vapi/internal/security"
	"vapi/internal/tenant"
	"vcache"
	"vdi"
)

type AppContainer struct {
	Config    vdi.Singleton[AppContainer, *config.ConfigService]
	Cache     vdi.Singleton[AppContainer, vcache.Cache]
	DbContext vdi.Singleton[AppContainer, *dbcontext.DbContext]
	Tenant    vdi.Singleton[AppContainer, *tenant.TenantService]
	Context   vdi.Transient[AppContainer, *ctxSvc.ContextService]
	Security  vdi.Transient[AppContainer, *security.SecurityPolicyService]
}

// Hàm khởi tạo và đăng ký tất cả dịch vụ
func NewAppContainer() (*AppContainer, error) {
	app, err := vdi.RegisterContainer(func(c *AppContainer) error {
		c.Config.Init = func(owner *AppContainer) *config.ConfigService {
			cfg, err := config.NewConfigService("./../config/config.yaml")
			if err != nil {
				panic(err)
			}
			return cfg
		}
		c.Cache.Init = func(owner *AppContainer) vcache.Cache {
			cache, err := cache.NewCacheService(c.Config.Get())
			if err != nil {
				return nil
			}
			return cache
		}
		c.DbContext.Init = func(owner *AppContainer) *dbcontext.DbContext {
			db, err := dbcontext.NewDbContext(c.Config.Get())
			if err != nil {
				panic(err)
			}
			return db
		}
		c.Tenant.Init = func(owner *AppContainer) *tenant.TenantService {
			config := c.Config.Get()

			return tenant.NewTenantService(c.DbContext.Get().DB, config.Get().Database.Manager)
		}
		c.Context.Init = func(owner *AppContainer) *ctxSvc.ContextService {
			return ctxSvc.NewContextService()
		}
		c.Security.Init = func(owner *AppContainer) *security.SecurityPolicyService {
			ret := security.NewSecurityPolicyService(c.Cache.Get(), c.DbContext.Get().DB)
			ret.CtxSvc = c.Context.Get()
			return ret
		}
		return nil
	})
	return app.Get(), err
}
