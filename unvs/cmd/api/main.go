package main

import (
	"context"
	"dbx"
	"log"
	"reflect"
	"time"
	handler "unvs/internal/app/handler"
	caller "unvs/internal/app/handler/callers"
	"unvs/internal/app/handler/inspector"
	_ "unvs/internal/model/base"

	_ "unvs/internal/app/middleware/auth"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "unvs/docs" // Import thư mục chứa docs đã tạo bởi swag

	cache "caching"

	_ "unvs/internal/app/handler/inspector"
	_ "unvs/views_business/auth"
	_ "unvs/views_business/users"

	"net/http"
	_ "net/http/pprof"
	oauthHandler "unvs/internal/app/handler/oauth"

	echoSwagger "github.com/swaggo/echo-swagger" // Thư viện tích hợp Swagger cho Echo
)

func getMssqlConfig() dbx.Cfg {
	return dbx.Cfg{
		Driver:   "mssql",
		Host:     "localhost",
		Port:     0,
		User:     "sa",
		Password: "123456",
	}

}
func getPgConfig() dbx.Cfg {
	return dbx.Cfg{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "123456",
	}
}
func getMysqlConfig() dbx.Cfg {
	return dbx.Cfg{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "123456",
	}
}
func getMemoryCache(ownerType reflect.Type) cache.Cache {
	return cache.NewInMemoryCache(
		ownerType,
		time.Minute*5, time.Minute*10,
	)

}
func getBadgerCache(ownerType reflect.Type) cache.Cache {
	ret, err := cache.NewBadgerCache(
		ownerType,
		"unvs",
	)
	if err != nil {
		panic(err)
	}
	return ret
}
func getMemcachedServer(ownerType reflect.Type) cache.Cache {
	return cache.NewMemcachedCache(
		ownerType,
		[]string{"127.0.0.1:11211"},
	)
}
func getRedisCached(ownerType reflect.Type) cache.Cache {
	return cache.NewRedisCache(
		context.Background(),
		ownerType,
		"127.0.0.1:6379",
		"",
		0,
	)
}

// func createTenantDb(tenant string) (*dbx.DBXTenant, error) {
// 	cfg := getMssqlConfig()
// 	db := dbx.NewDBX(cfg)
// 	db.Open()
// 	defer db.Close()
// 	if err := db.Ping(); err != nil {
// 		return nil, err
// 	}
// 	tenantDB, err := db.GetTenant(tenant)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return tenantDB, nil

// }
// func createUserRepo(tenantDB *dbx.DBXTenant) user_repo.UserRepository {
// 	return user_repo.NewUserRepo(tenantDB)
// }
// func createAccService(tenantDB *dbx.DBXTenant) *account.AccountService {
// 	return account.NewAccountService(
// 		createUserRepo(tenantDB),
// 		//getBadgerCache(reflect.TypeOf(account.AccountService{})),
// 		getMemoryCache(reflect.TypeOf(account.AccountService{})),
// 		//getMemcachedServer(reflect.TypeOf(account.AccountService{})),
// 		//getRedisCached(reflect.TypeOf(account.AccountService{})),
// 	)
// }
// func createAccHandler(tenantDB *dbx.DBXTenant) accHandler.AccountHandler {
// 	return *accHandler.NewAccountHandler(createAccService(tenantDB))
// }
// func createOAuthHandler(tenantDB *dbx.DBXTenant) oauthHandler.OAuthHandler {

// 	return *oauthHandler.NewOAuthHandler(createAccService(tenantDB))
// }

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
	defer dbx.CloseAll()
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	// var tenantDB *dbx.DBXTenant
	// tenantDB, err := createTenantDb("tenant1")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// tenantDB.Open()
	// defer tenantDB.Close()
	// //accHandlers := createAccHandler(tenantDB)
	// oauthHandler := createOAuthHandler(tenantDB)
	// Khởi tạo Echo
	// if err != nil {
	// 	log.Fatal(err)
	// }

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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
	log.Fatal(e.Start(":8080"))
}
