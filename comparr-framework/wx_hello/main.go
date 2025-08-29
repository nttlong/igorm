package main

import (
	"fmt"
	"mime/multipart"
	"wx"
)

type Address struct {
	City     string `json:"city"`
	District string `json:"district"`
	Street   string `json:"street"`
}

type UserInput struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Email   string   `json:"email"`
	Phones  []string `json:"phones"`
	Address Address  `json:"address"`
}

type TestApi struct {
}

func (t *TestApi) DoSayHelloInVn(name string) string {
	return "xin chao," + name
}
func (t *TestApi) DoSayHelloInEn(name string) string {
	return "hello" + name
}
func (t *TestApi) Hello(h *struct {
	wx.Handler `route:"hello/{name}/{LangCode};method:get"`
	Name       string
	LangCode   string
}) (interface{}, error) {
	var err error
	var msg string
	if h.Name == "" {
		err = wx.Errors.RequireErr("name")
	}
	if h.LangCode == "" {
		err = wx.Errors.RequireErr("langCode")
	}
	if h.LangCode == "vn" {
		msg = t.DoSayHelloInVn(h.Name)
	}
	if h.LangCode == "en" {
		msg = t.DoSayHelloInEn(h.Name)
	}
	if err != nil {
		return nil, err
	}
	return map[string]string{"msg": msg}, nil
}
func (t *TestApi) CreateUser(h *wx.Handler,
	Body UserInput, //<-- cach 1 dat body o day wx tu parse
) (string, error) {
	return fmt.Sprintf("User %s, %d tuổi, sống ở %s",
		Body.Name, Body.Age, Body.Address.City), nil
}
func (t *TestApi) Upload(
	h *wx.Handler, //<-- day la handler, mac dinh la POST, hanlder uri lay theo ten Method
	data *struct {
		File *multipart.FileHeader `form:"file"`
	},
) error {
	if data.File != nil {
		fs, err := data.File.Open()
		if err != nil {
			return err
		}
		defer fs.Close()

		return nil
	}
	return wx.Errors.RequireErr("file")

}
func main() {
	wx.Routes("/api", &TestApi{})
	server := wx.NewHtttpServer("/api", "5000", "0.0.0.0")

	swagger := wx.CreateSwagger(server, "docs")
	swagger.Info(wx.SwaggerInfo{
		Title:       "Swagger Example API",
		Description: "This is a sample server Petstore server.",
		Version:     "1.0.0",
	})
	err := swagger.Build()
	if err != nil {
		panic(err)
	}
	err = server.Start()
	if err != nil {
		panic(err)
	}

}
