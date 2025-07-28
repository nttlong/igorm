package app

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

type AppService[T any] struct {
	Host      string
	Port      int
	App       *echo.Echo
	Container *T
}

const containerKey = "app_container"

func (appSvr *AppService[T]) Setup(onRequestFn func(owner *AppService[T], c echo.Context) error) *AppService[T] {
	e := echo.New()
	appSvr.App = e
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(containerKey, appSvr.Container)
			onRequestFn(appSvr, c)
			return next(c)
		}
	})
	return appSvr

}

func (appSvr *AppService[T]) Run() error {

	return appSvr.App.Start(appSvr.Host + ":" + strconv.Itoa(appSvr.Port))

}
