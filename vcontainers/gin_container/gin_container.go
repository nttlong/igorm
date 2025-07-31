package gin_container

import (
	"core_container"
	"strconv"
	"vdi"

	"github.com/gin-gonic/gin"
)

type GinContainer struct {
	Core vdi.Singleton[GinContainer, *core_container.CoreContainer]
	Gin  vdi.Singleton[GinContainer, *gin.Engine]
}

func NewGinContainer(configFilePath string) *GinContainer {
	ret := vdi.NewContainer(func(owner *GinContainer) error {

		owner.Core.Init = func(owner *GinContainer) *core_container.CoreContainer {
			return core_container.NewCoreContainer(configFilePath)
		}

		owner.Gin.Init = func(owner *GinContainer) *gin.Engine {
			r := gin.Default()
			// Có thể gắn middleware, logger, recovery ở đây nếu cần
			return r
		}

		return nil
	})
	return ret
}

func (c *GinContainer) StartServer() error {
	core := c.Core.Get()
	cfg := core.Config.Get()
	appConfig := cfg.GetAppConfig()

	r := c.Gin.Get()
	err := r.Run(appConfig.Host + ":" + strconv.Itoa(appConfig.Port))
	return err
}
