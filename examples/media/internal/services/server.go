package services

import (
	"media/handlers"
	"wx"
)

type Server struct {
}

func (s *Server) Start() error {
	wx.Controller(func() (*handlers.Media, error) {
		return &handlers.Media{}, nil
	})

	server := wx.NewHtttpServer("/api", 8081, "0.0.0.0")
	swagger := wx.CreateSwagger(server, "swagger")
	swagger.OAuth2Password(server.BaseUrl + "oauth/token")
	swagger.Info(wx.SwaggerInfo{
		Title:       "Exmaple Media API",
		Description: "This is a sample server Petstore server.",
		Version:     "1.0.0",
	})
	swagger.Build()

	err := server.Start()
	if err != nil {
		return err
	}
	return nil

}
