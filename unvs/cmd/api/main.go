package main

import (
	"log"
	handler "unvs/internal/app/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "unvs/docs" // Import thư mục chứa docs đã tạo bởi swag

	echoSwagger "github.com/swaggo/echo-swagger" // Thư viện tích hợp Swagger cho Echo
)

// @title Echo Hello World API
// @version 1.0
// @description Đây là một API ví dụ đơn giản với Echo và Swagger.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
//

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route để phục vụ Swagger UI
	// Sau khi chạy 'swag init', các file docs sẽ được tạo và route này sẽ hiển thị UI.
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Định nghĩa route API "Hello World"
	e.GET("/hz", handler.HzHandler) // Gọi handler tại đây

	// Khởi chạy server
	log.Println("Server đang lắng nghe tại cổng :8080")
	log.Println("Truy cập Swagger UI tại: http://localhost:8080/swagger/index.html")
	log.Fatal(e.Start(":8080"))
}
