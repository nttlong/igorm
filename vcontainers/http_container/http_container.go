package http_container

import (
	"core_container"
	"net/http"
	"strconv"
	"vdi"
)

type HttpContainer struct {
	Core    vdi.Singleton[HttpContainer, *core_container.CoreContainer]
	Handler vdi.Singleton[HttpContainer, http.Handler] // Cho phép inject router
}

func NewHttpContainer(configFilePath string) *HttpContainer {
	ret := vdi.NewContainer(func(owner *HttpContainer) error {

		owner.Core.Init = func(owner *HttpContainer) *core_container.CoreContainer {
			return core_container.NewCoreContainer(configFilePath)
		}

		owner.Handler.Init = func(owner *HttpContainer) http.Handler {
			// Mặc định là mux, có thể override bởi user
			mux := http.NewServeMux()
			mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("pong"))
			})
			return mux
		}

		return nil
	})
	return ret
}

func (c *HttpContainer) StartServer() error {
	core := c.Core.Get()
	cfg := core.Config.Get()
	appCfg := cfg.GetAppConfig()

	addr := appCfg.Host + ":" + strconv.Itoa(appCfg.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: c.Handler.Get(),
	}

	return server.ListenAndServe()
}
