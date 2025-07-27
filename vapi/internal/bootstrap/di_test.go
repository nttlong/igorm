package bootstrap

import (
	"context"
	"testing"
	"time"
	"vapi/internal/cache"
	"vapi/internal/config"
	ctxSvc "vapi/internal/context_service"
	"vapi/internal/dbcontext"
	"vapi/internal/security"
	"vapi/internal/security/models"
	"vapi/internal/tenant"
	"vcache"
	"vdb"

	"github.com/stretchr/testify/assert"
)

func TestDI(t *testing.T) {
	// Create a new container
	c, err := NewAppContainer()
	assert.NoError(t, err)
	assert.IsType(t, &AppContainer{}, c)
	// Get the config service
}
func BenchmarkTestDI(b *testing.B) {
	c := (&AppContainer{}).New(func(svc *AppContainer) error {
		var err error
		svc.Config.Init = func(owner *AppContainer) *config.ConfigService {
			ret := &config.ConfigService{}
			ret.New("./../config/config.yaml")
			return ret

		}
		svc.DbContext.Init = func(owner *AppContainer) *dbcontext.DbContext {
			val := &dbcontext.DbContext{}
			err = val.New(svc.Config.Get().Get().Database.Driver, svc.Config.Get().Get().Database.Dsn)
			if err != nil {
				panic(err)
			}
			svc.DbContext.Set(val)

			return val
		}
		svc.Tenant.Init = func(owner *AppContainer) *tenant.TenantService {
			ret := &tenant.TenantService{}
			err := ret.New(svc.DbContext.Get().DB, svc.Config.Get().Get().Database.Manager)
			if err != nil {
				panic(err)
			}
			return ret
		}
		svc.Cache.Init = func(owner *AppContainer) vcache.Cache {

			ret, err := cache.NewCacheService(svc.Config.Get())
			if err != nil {

				panic(err)
			}
			return ret
		}
		svc.Security.Init = func(owner *AppContainer) *security.SecurityPolicyService {
			ret := security.NewSecurityPolicyService(
				func() context.Context {
					return context.Background()
				}(),
				svc.Cache.Get(),
				func() *vdb.TenantDB {
					ret, err := svc.DbContext.Get().DB.CreateDB("test0001")
					if err != nil {
						panic(err)
					}
					return ret
				}(),
			)
			return ret
		}
		svc.Context.Init = func(owner *AppContainer) *ctxSvc.ContextService {
			svc.Context.Value = &ctxSvc.ContextService{}
			svc.Context.Value.Ctx = b.Context()
			return svc.Context.Value
		}

		return err
	})

	assert.NoError(b, c.Error)
	assert.IsType(b, &AppContainer{}, c)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.ResovleContext(func() context.Context {
			return context.Background()
		})

		err := c.Security.Get().CreateOrUpdate(&models.SecurityPolicy{
			TenantID:         "tenant1",
			JwtSecret:        "abc123",
			MaxLoginFailures: 5,
			LockoutMinutes:   15,
			JwtExpireMinutes: 60,
			CreatedAt:        time.Now().UTC(),
		})
		assert.NoError(b, err)
	}

}
