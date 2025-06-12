package main

import (
	"dbx"
	"io"
	"log"
	"os"
	"path/filepath"
	handler "unvs/internal/app/handler"
	caller "unvs/internal/app/handler/callers"
	"unvs/internal/app/handler/inspector"
	oauthHandler "unvs/internal/app/handler/oauth"
	"unvs/internal/config"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "unvs/docs" // Import thư mục chứa docs đã tạo bởi swag

	"net/http"
	_ "net/http/pprof"

	echoSwagger "github.com/swaggo/echo-swagger" // Thư viện tích hợp Swagger cho Echo
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
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	logPath := config.AppConfigInstance.Logs
	logDir := filepath.Dir(logPath)
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		panic(err)
	}
	// logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("Không thể mở file log '%s': %v", logPath, err)
	}
	// defer logFile.Close()
	defer dbx.CloseAll()
	lumberjackLogger := &lumberjack.Logger{
		Filename: logPath, // Tên file log chính

		MaxSize:    2,    // log xoay vòng khi đạt 10 MB
		MaxBackups: 10,   // Giữ lại tối đa 3 file log cũ đã xoay vòng
		MaxAge:     28,   // Xóa các file log cũ hơn 28 ngày
		Compress:   true, // Nén các file log cũ (.gz)
	}
	mw := io.MultiWriter(os.Stdout, lumberjackLogger)
	log.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	e := echo.New()

	// Middleware

	e.Use(middleware.Recover())
	// Cấu hình logger của Echo
	e.Logger.SetOutput(mw)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}",` +
			`"method":"${method}","uri":"${uri}","status":${status},` +
			`"latency_human":"${latency_human}"}` + "\n",
		Output: mw,
	}))
	// Route để phục vụ Swagger UI
	// Sau khi chạy 'swag init', các file docs sẽ được tạo và route này sẽ hiển thị UI.
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	oathHandler := &oauthHandler.OAuthHandler{}
	callHandler := &caller.CallerHandler{}
	e.POST("/oauth/token", oathHandler.Token)
	apiV1 := e.Group("/api/v1")
	apiV1.POST("/invoke/:action", callHandler.Call)
	handler.RegisterRoutes(e,

		&inspector.InspectorHandler{},
	)

	// Khởi chạy server
	log.Println("Server đang lắng nghe tại cổng :8080")
	log.Println("Truy cập Swagger UI tại: http://localhost:8080/swagger/index.html")
	log.Fatal(e.Start(config.AppConfigInstance.Server.Bind + ":" + config.AppConfigInstance.Server.Port))
}
