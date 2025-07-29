package echo_container

import (
	"core_container"
	_ "core_container"
	"strconv"
	"vdi"

	"github.com/labstack/echo"
)

type EchoContainer struct {
	Core vdi.Singleton[EchoContainer, *core_container.CoreContainer]
	Echo vdi.Singleton[EchoContainer, *echo.Echo]
}

func NewEchoContainer(configFilePath string) *EchoContainer {

	ret := vdi.NewContainer(func(owner *EchoContainer) error {

		owner.Core.Init = func(owner *EchoContainer) *core_container.CoreContainer {
			return core_container.NewCoreContainer(configFilePath)
		}
		owner.Echo.Init = func(owner *EchoContainer) *echo.Echo {
			ret := echo.New()
			return ret
		}
		return nil
	})
	// ret := EchoContainer{
	// 	Core: core,
	// }

	// ret.New(func(owner *EchoContainer) error {
	// 	//nothing for setting up echo container

	// 	owner.Echo.Init = func(owner *EchoContainer) *echo.Echo {

	// 		ret := echo.New()
	// 		return ret
	// 	}
	// 	return nil
	// })

	return ret
}
func (c *EchoContainer) StartServer() error {
	core := c.Core.Get()
	cfg := core.Config.Get()
	appConfig := (cfg).GetAppConfig()
	e := c.Echo.Get()
	err := e.Start(appConfig.Host + ":" + strconv.Itoa(appConfig.Port))
	return err

}
