package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" // import để tự đăng ký pprof handlers
	"os"
	"reflect"
	"runtime/pprof"
	"wx"
	"wx/example"
	_ "wx/example/example1/controllers"
	"wx/mw"
)

type TestController struct {
}

func main() {
	go func() {
		f, _ := os.Create("mem.pprof")
		pprof.WriteHeapProfile(f)
		f.Close()
		log.Println("pprof listening on :6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	wx.Routes("/api/v1", reflect.TypeFor[example.Media]())

	server := wx.NewHtttpServer("/api/v1", 8080, "localhost")
	uri, err := wx.GetUriOfHandler[example.Auth](server, "Oauth")
	if err != nil {
		panic(err)
	}
	log.Println(uri)

	swagger := wx.CreateSwagger(server, "/swagger")
	swagger.Info(wx.SwaggerInfo{
		Title:       "Swagger Example API",
		Description: "This is a sample server Petstore server.",
		Version:     "1.0.0",
	})

	swagger.OAuth2Password(uri)
	swagger.Build()
	server.Middleware(mw.LogAccessTokenClaims)
	server.Middleware(mw.Cors)
	//server.Middleware(mw.Zip)
	err = server.Start()
	if err != nil {
		panic(err)
	}

}
