package main

import (
	//_ "vapi/internal/bootstrap"
	_ "vapi/docs"
	bs "vapi/internal/bootstrap"

	//<-- co import o day
	"vapi/controllers"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Go API Example
// @version 1.0
// @description This is a sample API for demonstration purposes.
// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.bearer BearerAuth
// @description "JWT Authorization header using the Bearer scheme. Enter your token in the format 'Bearer <token>'"
// @name Authorization

// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl /oauth/token
// @in header
// @name Authorization
// @description "OAuth2 Password Flow - Enter email/username and password in the popup to get token."

// @in header
// @name Authorization
// @description "OAuth2 Password Flow (Form Submit) - Use for explicit form data submission."

func main() {
	container := bs.GetAppContainer("./config.yaml")
	appService := container.App.Get()

	/*
		App Service trong container la 1 dich vu de start echo, theo cau truc sau:
		type AppService struct {
			Host string
			Port int
			App  *echo.Echo
		}

	*/
	_echo := appService.App //<-- lay echo app
	/*
	 _echo luc nahy chinh la *echo.Echo
	*/
	_echo.POST("/download/code/:code", controllers.Test)
	_echo.POST("/oauth/token", controllers.OAuthToken)
	_echo.GET("/swagger/*", echoSwagger.WrapHandler)
	err := appService.Run() // chay web server
	if err != nil {
		panic(err)
	}
}

// const containerKey = "app_container"

// func InjectContainer(container *bootstrap.AppContainer) echo.MiddlewareFunc {
// 	return func(next echo.HandlerFunc) echo.HandlerFunc {
// 		return func(c echo.Context) error {
// 			c.Set(containerKey, container)
// 			return next(c)
// 		}
// 	}
// }

// func GetContainer(c echo.Context) *bootstrap.AppContainer {
// 	container, ok := c.Get(containerKey).(*bootstrap.AppContainer)
// 	if !ok {
// 		return nil
// 	}
// 	return container
// }
