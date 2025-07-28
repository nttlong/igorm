package bootstrap

import (
	"context"
	ps "vapi/internal/PasswordService"
	ss "vapi/internal/SharedSecret"
	"vapi/internal/account"
	"vapi/internal/app"
	"vapi/internal/cache"
	"vapi/internal/config"
	"vapi/internal/dbcontext"
	jwtservice "vapi/internal/jwt_service"
	vLogger "vapi/internal/logger"
	"vapi/internal/security"
	"vapi/internal/tenant"
	"vcache"
	"vdb"
	"vdi"

	"github.com/labstack/echo/v4"
)

type AppContainer struct {
	*vdi.Container[AppContainer]
	Config       vdi.Singleton[AppContainer, *config.ConfigService]
	Cache        vdi.Singleton[AppContainer, vcache.Cache]
	Db           vdi.Singleton[AppContainer, *vdb.TenantDB]
	CtxSvc       vdi.Transient[AppContainer, context.Context]
	DbService    vdi.Singleton[AppContainer, *dbcontext.DbContext]
	TenantDb     vdi.Singleton[AppContainer, *tenant.TenantService]
	SharedSecret vdi.Singleton[AppContainer, ss.SharedSecretService]
	Jwtservices  vdi.Singleton[AppContainer, jwtservice.JWTService]
	PwdSvc       vdi.Singleton[AppContainer, *ps.BcryptPasswordService]

	AccountSvc vdi.Singleton[AppContainer, *account.AccountService]

	Security      vdi.Singleton[AppContainer, *security.SecurityPolicyService]
	GetContext    func() context.Context
	GetDb         func() *vdb.TenantDB
	GetTenantName func() string
	App           vdi.Singleton[AppContainer, *app.AppService[AppContainer]]
	Logger        vdi.Singleton[AppContainer, *vLogger.LoggerService]
	Error         error
}

// Hàm khởi tạo và đăng ký tất cả dịch vụ
func GetAppContainer(
	configPath string,

) *AppContainer {
	c := (&AppContainer{}).New(func(svc *AppContainer) error {
		var err error
		svc.DbService.Init = func(owner *AppContainer) *dbcontext.DbContext {
			ret := &dbcontext.DbContext{}
			driver := svc.Config.Get().Get().Database.Driver
			dsn := svc.Config.Get().Get().Database.Dsn
			manager := svc.Config.Get().Get().Database.Manager
			vdb.SetManagerDb(driver, manager)
			db, err := vdb.Open(driver, dsn)
			if err != nil {
				svc.Error = err
				return nil
			}
			ret.DB = db
			return ret
		}
		svc.TenantDb.Init = func(owner *AppContainer) *tenant.TenantService {
			ret := &tenant.TenantService{}
			ret.Db = svc.Db.Get()

			return ret
		}
		svc.Db.Init = func(owner *AppContainer) *vdb.TenantDB {
			driver := svc.Config.Get().Get().Database.Driver
			dsn := svc.Config.Get().Get().Database.Dsn
			manager := svc.Config.Get().Get().Database.Manager
			vdb.SetManagerDb(driver, manager)
			db, err := vdb.Open(driver, dsn)
			if err != nil {
				svc.Error = err
				return nil
			}

			return db
		}
		svc.Config.Init = func(owner *AppContainer) *config.ConfigService {
			ret := &config.ConfigService{}
			ret.New(configPath)
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
			if svc.GetContext == nil {
				panic("GetContext is not set")
			}

			ret := &security.SecurityPolicyService{}
			ret.Cache = svc.Cache.Get()
			ret.GetDb = svc.GetDb
			ret.GetCtx = svc.GetContext

			return ret
		}
		svc.GetDb = func() *vdb.TenantDB {
			if svc.GetTenantName == nil {
				panic("GetTenantName is not set")
			}

			tenantName := svc.GetTenantName()
			ret, err := svc.TenantDb.Get().Tenant(tenantName)
			if err != nil {
				panic(err)
			}
			return ret

		}
		svc.SharedSecret.Init = func(owner *AppContainer) ss.SharedSecretService {
			ret := ss.NewSharedSecretService()

			return ret
		}
		svc.PwdSvc.Init = func(owner *AppContainer) *ps.BcryptPasswordService {
			ret := ps.NewBcryptPasswordService(0)
			return ret
		}
		svc.Jwtservices.Init = func(owner *AppContainer) jwtservice.JWTService {
			ret := jwtservice.NewJWTService()
			return ret
		}

		svc.AccountSvc.Init = func(owner *AppContainer) *account.AccountService {
			ret := &account.AccountService{}
			ret.Cache = svc.Cache.Get()
			ret.GetDb = svc.GetDb
			ret.GetCtx = svc.GetContext
			ret.PasswordHasher = svc.PwdSvc.Get()
			ret.PolicySvc = svc.Security.Get()
			ret.JwtSvc = svc.Jwtservices.Get()

			return ret
		}
		svc.Logger.Init = func(owner *AppContainer) *vLogger.LoggerService {
			log := owner.Config.Get().Get().Log
			ret := vLogger.NewLoggerService(&vLogger.LoggerConfig{
				FileName:   log.FileName,
				MaxSize:    log.MaxSize,
				MaxAge:     log.MaxAge,
				MaxBackups: log.MaxBackups,
				Compress:   log.Compress,
				AlsoStdout: log.AlsoStdout,
			})

			return ret
		}

		svc.App.Init = func(owner *AppContainer) *app.AppService[AppContainer] {
			host := svc.Config.Get().Get().Host
			port := svc.Config.Get().Get().Port
			ret := &app.AppService[AppContainer]{
				Host:      host,
				Port:      port,
				Container: owner,
			}
			ret.Setup(func(owner *app.AppService[AppContainer], c echo.Context) error {
				tenantName := "tenant1"
				owner.Container.GetContext = func() context.Context {
					return c.Request().Context()
				}
				owner.Container.GetTenantName = func() string {

					return tenantName
				}
				return nil

			})
			owner.Logger.Get().Apply(ret.App)

			return ret
		}

		return err
	})

	return c
}
