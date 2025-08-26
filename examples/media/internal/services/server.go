package services

import (
	"media/handlers"
	"os"
	"wx"
	"wx/mw"
)

type Server struct {
}

func (s *Server) Start() error {
	err := wx.Routes("/api", handlers.Media{}, handlers.Users{}, handlers.Logins{})
	if err != nil {
		return err
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := wx.NewHtttpServer("/api", port, "127.0.0.1")
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
		return err
	}
	return nil

}
