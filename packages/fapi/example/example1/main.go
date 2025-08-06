package main

import (
	"fapi"
	"fapi/example"
)

type TestController struct {
}

func main() {
	fapi.Controller("", "", func() (*example.Media, error) {
		return &example.Media{}, nil
	})
	server := fapi.NewHtttpServer("/api", 8080, "localhost")
	server.Swagger()
	server.Middleware(fapi.Cors)
	server.Middleware(fapi.Zip)
	err := server.Start()
	if err != nil {
		panic(err)
	}

}
