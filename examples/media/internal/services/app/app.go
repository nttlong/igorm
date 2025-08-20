package app

import (
	"media/internal/services"
	"wx"
)

type App struct {
	Server wx.Depend[services.Server]
}

func (app *App) New() error {
	// func(s*services.Server) Start(){

	// }

	app.Server.Init(func() (*services.Server, error) {
		return &services.Server{}, nil
	})
	return nil
}
