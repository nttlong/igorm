package main

import (
	"dbx"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	handler "unvs/internal/app/handler"
	caller "unvs/internal/app/handler/callers"
	"unvs/internal/app/handler/inspector"
	oauthHandler "unvs/internal/app/handler/oauth"
	"unvs/internal/config"

	_ "unvs.common"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	_ "unvs/docs" // Import thư mục chứa docs đã tạo bởi swag

	"net/http"
	// _ "net/http/pprof"

	echoSwagger "github.com/swaggo/echo-swagger" // Thư viện tích hợp Swagger cho Echo
)

var appLogger *logrus.Logger

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
	appLogger = logrus.New()
	go func() {
		s := &http.Server{
			Addr:         "localhost:6060",
			Handler:      nil,               // Your main Echo router should be here, not nil
			ReadTimeout:  10 * time.Second,  // Ví dụ: timeout đọc 10 giây
			WriteTimeout: 10 * time.Second,  // Ví dụ: timeout ghi 10 giây
			IdleTimeout:  120 * time.Second, // Ví dụ: timeout kết nối rảnh 120 giây
		}
		log.Println(s.ListenAndServe())
	}()
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	logPath := config.AppConfigInstance.Logs
	logDir := filepath.Dir(logPath)
	err = os.MkdirAll(logDir, 0750)
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

	appLogger.SetOutput(mw)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// --------------------------------------------------------

	// 3. Cấu hình Formatter cho Logrus
	// Sử dụng JSONFormatter để có log cấu trúc, dễ dàng thêm các trường như "package"
	// appLogger.SetFormatter(&logrus.JSONFormatter{
	// 	TimestampFormat: time.RFC3339Nano, // Định dạng thời gian chi tiết hơn
	// })
	// appLogger.SetFormatter(&logrus.TextFormatter{
	// 	DisableColors:   false,        // Đặt true nếu bạn chạy trên môi trường không hỗ trợ màu (ví dụ: production server không có terminal)
	// 	FullTimestamp:   true,         // Hiển thị đầy đủ thông tin thời gian (ngày, giờ, phút, giây, nano giây)
	// 	TimestampFormat: time.RFC3339, // Định dạng thời gian (ví dụ: "2006-01-02T15:04:05Z07:00")
	// 	// Có thể thêm thêm các tùy chọn khác như:
	// 	// DisableSorting:  true, // Không sắp xếp các trường theo thứ tự abc
	// 	// QuoteEmptyFields: true, // Thêm dấu ngoặc kép cho các trường rỗng
	// 	// FieldMap:        logrus.FieldMap{ // Map lại tên các trường mặc định nếu muốn
	// 	// 	logrus.FieldKeyTime:  "@timestamp",
	// 	// 	logrus.FieldKeyLevel: "level",
	// 	// 	logrus.FieldKeyMsg:   "message",
	// 	// },
	// })
	// 4. Cấu hình Level cho Logrus
	// appLogger.SetLevel(logrus.TraceLevel | logrus.InfoLevel | logrus.ErrorLevel | logrus.DebugLevel | logrus.PanicLevel | logrus.FatalLevel | logrus.WarnLevel) // Chỉ hiển thị Info, Warn, Error, Fatal, Panic

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	angularDistPath := `E:\Docker\go\unvs\igorm\unvs-ng-ui\dist\unvs-ng-ui\browser`
	e.Static("/", angularDistPath)
	e.GET("/*", func(c echo.Context) error {
		//check if url has a file extension
		if filepath.Ext(c.Request().URL.Path) != "" {
			return c.File(filepath.Join(angularDistPath, c.Request().URL.Path))
		} else {
			return c.File(filepath.Join(angularDistPath, "index.html"))
		}
	})
	// Fallback cho tất cả các route không khớp với file tĩnh về index.html
	// Điều này là cần thiết cho Angular Router (Single Page Application)
	// CATCH-ALL ROUTE CHO ANGULAR ROUTER (SPA Fallback).
	// Đây phải là route CUỐI CÙNG được định nghĩa.
	// Nếu không có API, Swagger, hoặc file tĩnh nào khớp, nó sẽ trả về index.html
	// để Angular Router xử lý client-side routing.
	// e.GET("/*", func(c echo.Context) error {
	// 	return c.File(filepath.Join(angularDistPath, "index.html"))
	// })

	// Cấu hình CORS middleware
	// Đây là nơi bạn kiểm soát các chính sách cross-domain

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// Đảm bảo tên trường là chính xác: AllowedOrigins
		// Đây phải là một slice of strings
		AllowOrigins: []string{
			"*",
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		},
		MaxAge: 86400,
	}))

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
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5, // Cấp độ nén từ 1 (nhanh nhất) đến 9 (tốt nhất). Mặc định là -1 (LevelDefault).
	}))
	// Route để phục vụ Swagger UI
	// Sau khi chạy 'swag init', các file docs sẽ được tạo và route này sẽ hiển thị UI.

	oathHandler := &oauthHandler.OAuthHandler{}
	callHandler := &caller.CallerHandler{
		AppLogger: appLogger,
	}
	e.POST("/oauth/token", oathHandler.Token)
	apiV1 := e.Group("/api/v1")

	apiV1.POST("/invoke", callHandler.Call)
	apiV1.POST("/invoke-form", callHandler.FormSubmit)
	apiV1.POST("/oauth/token", oathHandler.Token)
	///get/{tenant}/{module}/{action}/{optionalPath}
	apiV1.GET("/get/:tenant/:module/:action/*optionalPath", callHandler.CallGet)
	// e.GET("/get/:tenant/:module/:action/*paramoptionalPaths", handler.CallGet)
	// apiV1.GET("/get/:tenant/:module/:action/:feature/:lan", callHandler.CallGet)

	handler.RegisterRoutes(e,

		&inspector.InspectorHandler{},
	)

	// Khởi chạy server
	log.Println("Server đang lắng nghe tại cổng :8080")
	log.Println("Truy cập Swagger UI tại: http://localhost:8080/swagger/index.html")
	log.Fatal(e.Start(config.AppConfigInstance.Server.Bind + ":" + config.AppConfigInstance.Server.Port))
}
