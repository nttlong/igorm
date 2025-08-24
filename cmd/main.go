package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"wx"
	_ "wx"
	"wx/mw"
)

func main() {
	log.Print("starting server...")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	portNumber, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	server := wx.NewHtttpServer("/api", portNumber, "127.0.0.1")
	swagger := wx.CreateSwagger(server, "swagger")
	swagger.OAuth2Password(server.BaseUrl + "oauth/token")
	swagger.Info(wx.SwaggerInfo{
		Title:       "Exmaple Media API",
		Description: "This is a sample server Petstore server.",
		Version:     "1.0.0",
	})
	swagger.Build()
	//server.Middleware(mw.Zip)
	server.Middleware(mw.Cors)
	err = server.Start()
	if err != nil {
		panic(err)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
