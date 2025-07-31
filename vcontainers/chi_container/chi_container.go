package chi_container

import (
	"core_container"
	"net/http"
	"strconv"
	"vdi"

	"github.com/go-chi/chi/v5"
)

type ChiContainer struct {
	Core vdi.Singleton[ChiContainer, *core_container.CoreContainer]
	Chi  vdi.Singleton[ChiContainer, chi.Router]
}

func NewChiContainer(configFilePath string) *ChiContainer {
	ret := vdi.NewContainer(func(owner *ChiContainer) error {

		owner.Core.Init = func(owner *ChiContainer) *core_container.CoreContainer {
			return core_container.NewCoreContainer(configFilePath)
		}

		owner.Chi.Init = func(owner *ChiContainer) chi.Router {
			r := chi.NewRouter()

			// Middleware ví dụ (nếu cần):
			// r.Use(middleware.Logger)
			// r.Use(middleware.Recoverer)

			return r
		}

		return nil
	})
	return ret
}

func (c *ChiContainer) StartServer() error {
	core := c.Core.Get()
	cfg := core.Config.Get()
	appConfig := cfg.GetAppConfig()

	r := c.Chi.Get()
	addr := appConfig.Host + ":" + strconv.Itoa(appConfig.Port)

	err := http.ListenAndServe(addr, r)
	return err
}
