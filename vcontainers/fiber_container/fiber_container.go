package fiber_container

import (
	"core_container"
	"strconv"
	"vdi"

	"github.com/gofiber/fiber/v2"
)

type FiberContainer struct {
	Core  vdi.Singleton[FiberContainer, *core_container.CoreContainer]
	Fiber vdi.Singleton[FiberContainer, *fiber.App]
}

func NewFiberContainer(configFilePath string) *FiberContainer {
	ret := vdi.NewContainer(func(owner *FiberContainer) error {

		owner.Core.Init = func(owner *FiberContainer) *core_container.CoreContainer {
			return core_container.NewCoreContainer(configFilePath)
		}

		owner.Fiber.Init = func(owner *FiberContainer) *fiber.App {
			app := fiber.New()
			// Có thể gắn middleware, logger tại đây nếu cần
			return app
		}

		return nil
	})
	return ret
}

func (c *FiberContainer) StartServer() error {
	core := c.Core.Get()
	cfg := core.Config.Get()
	appConfig := cfg.GetAppConfig()

	app := c.Fiber.Get()
	err := app.Listen(appConfig.Host + ":" + strconv.Itoa(appConfig.Port))
	return err
}
