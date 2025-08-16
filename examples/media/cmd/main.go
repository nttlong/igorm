package main

import (
	apps "media/internal/services/app"
	"wx"
)

func main() {
	err := wx.Start(func(app *apps.App) error {
		server, err := app.Server.Ins()
		if err != nil {
			return err
		}
		server.Start()
		return nil
	})
	if err != nil {
		panic(err)
	}
}
