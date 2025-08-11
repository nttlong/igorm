package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" // import để tự đăng ký pprof handlers
	"os"
	"runtime/pprof"
	"vapi"
	"vapi/example"
	_ "vapi/example/example1/controllers"
	"vapi/mw"
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
	vapi.Controller(func() (*example.Media, error) {
		return &example.Media{}, nil
	})
	server := vapi.NewHtttpServer("/api/v1", 8080, "localhost")
	vapi.SwaggerUtils.OAuth2Password(
		"/api/oauth/token",
		"",
	)
	server.Swagger()
	server.Middleware(mw.LogAccessTokenClaims)
	server.Middleware(mw.Cors)
	//server.Middleware(mw.Zip)
	err := server.Start()
	if err != nil {
		panic(err)
	}

}
