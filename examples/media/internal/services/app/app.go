package app

import (
	"media/internal/services"
	"wx"
)

type App struct {
	Files  *wx.Depend[services.FileService, App]
	Server *wx.Depend[services.Server, App]
}

func (app *App) New() error {
	app.Files.Init(func(app *App) services.FileService {
		return services.FileService{}
	})
	app.Server.Init(func(app *App) services.Server {
		return services.Server{}
	})
	return nil
}
