//go:build debug
// +build debug

package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" // import để tự đăng ký pprof handlers
	"os"
	"reflect"
	"runtime/pprof"
	"time"
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
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		if err := f.Close(); err != nil {
			log.Fatal("could not close memory profile file: ", err)
		}
		log.Println("pprof listening on :6060")
		server := &http.Server{
			Addr: "localhost:6060",
			// Handler:      s.handler,
			ReadTimeout:  10 * time.Second, // Giới hạn đọc request
			WriteTimeout: 10 * time.Second, // Giới hạn ghi response
			IdleTimeout:  60 * time.Second, // Cho keep-alive
		}

		if err := server.ListenAndServe(); err != nil {
			log.Fatal("could not start pprof server: ", err)
		}

	}()
	if err := wx.Routes("/api/v1", reflect.TypeFor[example.Media]()); err != nil {
		panic(err)
	}

	server := wx.NewHtttpServer("/api/v1", "8080", "localhost")
	uri, err := wx.GetUriOfHandler[example.Auth]("Oauth")
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
	if err := swagger.Build(); err != nil {
		panic(err)
	}
	server.Middleware(mw.LogAccessTokenClaims)
	server.Middleware(mw.Cors)
	//server.Middleware(mw.Zip)
	err = server.Start()
	if err != nil {
		panic(err)
	}

}
