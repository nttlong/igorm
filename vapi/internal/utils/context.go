package utils

import (
	"vapi/internal/bootstrap"

	"github.com/labstack/echo/v4"
)

func GetContainer(c echo.Context) *bootstrap.AppContainer {
	val := c.Get("app_container")
	if val == nil {
		panic("AppContainer not found in context")
	}
	return val.(*bootstrap.AppContainer)
}
