package core_container

import (
	"config_service"
	_ "config_service"
	"vcache"
	"vdi"
)

type CoreContainer struct {
	*vdi.Container[CoreContainer]
	Config vdi.Singleton[CoreContainer, config_service.ConfigService]
	Cache  vdi.Singleton[CoreContainer, vcache.Cache]
}

/*
Setup the container with the necessary dependencies.
*/
func NewCoreContainer(configFilePath string) *CoreContainer {
	// c := &CoreContainer{}
	c := vdi.NewContainer[CoreContainer](func(owner *CoreContainer) error {
		owner.Config.Init = func(owner *CoreContainer) config_service.ConfigService {
			ret, err := config_service.NewConfigService(configFilePath)
			if err != nil {
				panic(err)
			}
			return ret

		}
		owner.Cache.Init = func(owner *CoreContainer) vcache.Cache {
			appConfig := (owner.Config.Get()).GetAppConfig()
			switch appConfig.CacheType {
			case "redis":
				return vcache.NewRedisCache(appConfig.Redis.Nodes, appConfig.Redis.Password, appConfig.Redis.PrefixKey, appConfig.Redis.DB, appConfig.Redis.Timeout)
			case "memcached":
				return vcache.NewMemcachedCache(appConfig.Memcached.Nodes, appConfig.Memcached.PrefixKey)
			case "badger":
				ret, err := vcache.NewBadgerCache(appConfig.Badger.Path, appConfig.Badger.PrefixKey)
				if err != nil {
					panic(err)
				}
				return ret
			case "inmemory":
				return vcache.NewInMemoryCache(0, 0)
			default:
				panic("Invalid cache type")
			}
		}
		return nil
	})

	return c
}
