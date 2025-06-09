package handlers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	// ... import các handler và middleware khác
	// Middleware
)

// RegisterRoutes đăng ký tất cả các route API vào Echo engine
func RegisterRoutes(e *echo.Echo, handlers ...interface{}) {
	// Swagger UI Route
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	apiV1 := e.Group("/api/v1")
	for _, handler := range handlers {
		// ... register handle into echo engine
		// detect handler is echo.HandlerFunc
		fmt.Println("Registering handler: ", handler)
		handlerType := reflect.TypeOf(handler)
		handlerStaticType := handlerType
		if handlerType.Kind() == reflect.Ptr {
			handlerStaticType = handlerType.Elem()
		}
		for i := 0; i < handlerType.NumMethod(); i++ {
			method := handlerType.Method(i)
			fmt.Println(handlerStaticType.PkgPath(), handlerStaticType.Name())
			fmt.Printf("Method Name: %s, Type: %s\n", method.Name, method.Type)
			//fx := method.Func.Addr().Interface().(func(echo.Context) error)
			//fmt.Println(fx)
			// method.Func.Call()
			apiPath := strings.Split(strings.Split(handlerStaticType.PkgPath(), "/handler/")[1], "/")[0]
			apiPath += "/" + method.Name
			apiPath = strings.ToLower(apiPath)
			fmt.Println("API Path: ", apiPath)

			apiV1.POST(fmt.Sprintf("/%s", apiPath), func(c echo.Context) error {
				args := []reflect.Value{
					reflect.ValueOf(handler),
					reflect.ValueOf(c)}
				ret := method.Func.Call(args)
				if len(ret) == 0 {
					return nil
				}
				if len(ret) == 1 {
					if ret[0].Interface() == nil {
						return nil
					}
					return ret[0].Interface().(error)
				}
				return nil

			})

		}
	}

	// apiV1 := e.Group("/api/v1")
	// {
	// 	apiV1.GET("/hz", healthHandler.HzHandler)
	// 	apiV1.POST("/accounts/create", accHandlers.CreateAccount, auth.JWTAuthMiddleware) // Sử dụng middleware từ internal/auth
	// 	apiV1.POST("/accounts/login", accHandlers.Login)
	// 	apiV1.POST("/oauth/token", accHandlers.LoginByFormSubmit) // Đây là endpoint trong api/v1 group
	// }

	// // Nếu bạn vẫn muốn có /oauth/token ở gốc (không khuyến khích nếu tất cả là api/v1)
	// e.POST("/oauth/token", accHandlers.LoginByFormSubmit)
}
