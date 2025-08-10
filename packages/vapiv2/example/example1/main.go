package main

import (
	"vapi"
	"vapi/example"
	_ "vapi/example/example1/controllers"
	"vapi/mw"
)

type TestController struct {
}

func main() {
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
	server.Middleware(mw.Zip)
	err := server.Start()
	if err != nil {
		panic(err)
	}

}
