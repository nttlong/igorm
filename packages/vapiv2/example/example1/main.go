package main

import (
	"vapi"
	"vapi/example"
)

type TestController struct {
}

func main() {
	vapi.Controller("", "", func() (*example.Media, error) {
		return &example.Media{}, nil
	})
	server := vapi.NewHtttpServer("/api", 8080, "localhost")
	vapi.SwaggerUtils.OAuth2Password(
		"api/oauth/token",
		"",
	)
	server.Swagger()
	server.Middleware(vapi.Cors)
	server.Middleware(vapi.Zip)
	err := server.Start()
	if err != nil {
		panic(err)
	}

}
