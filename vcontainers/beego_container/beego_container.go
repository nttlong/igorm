package beego_container

import (
	"core_container"
	"vdi"

	"github.com/beego/beego/v2/server/web"
)

type BeegoContainer struct {
	Core  vdi.Singleton[BeegoContainer, *core_container.CoreContainer]
	Beego vdi.Singleton[BeegoContainer, *web.HttpServer]
}

func NewBeegoContainer(configFilePath string) *BeegoContainer {
	ret := vdi.NewContainer(func(owner *BeegoContainer) error {

		owner.Core.Init = func(owner *BeegoContainer) *core_container.CoreContainer {
			return core_container.NewCoreContainer(configFilePath)
		}

		owner.Beego.Init = func(owner *BeegoContainer) *web.HttpServer {
			return nil
		}

		return nil
	})
	return ret
}

func (c *BeegoContainer) StartServer() error {
	panic("implement me")
}
