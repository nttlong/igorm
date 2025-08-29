package handlers

import (
	"fmt"
	"media/services"
	"wx"
)

type Media struct {
	FileDirectory string
	wx.ControllerContext
	File        *services.FileService
	Directories *services.DirectoryService
}

func (m *Media) New(
	fs *wx.Depend[services.FileService],
	ds *wx.Global[services.DirectoryService],
	urlSvc *wx.Depend[services.UrlResolverService],
) error {
	fmt.Println("Media.New")
	m.FileDirectory = "./uploads"
	ds.Init(func() (*services.DirectoryService, error) {
		return &services.DirectoryService{}, nil
	})

	ds.Init(func() (*services.DirectoryService, error) {
		return &services.DirectoryService{
			DirUpload:     "./uploads",
			DirUploadName: "uploads",
		}, nil
	})
	urlSvc.Init(func() (*services.UrlResolverService, error) {
		return &services.UrlResolverService{
			BaseUrl: m.BaseUrl,
		}, nil
	})
	fs.Init(func() (*services.FileService, error) {
		dirSvc, err := ds.Ins()
		if err != nil {
			return nil, err
		}
		urlSvc, err := urlSvc.Ins()
		if err != nil {
			return nil, err
		}

		return &services.FileService{
			DirectorySvc: &dirSvc,
			UrlSvc:       urlSvc,
		}, nil

	})
	fsSvc, err := fs.Ins()
	if err != nil {
		return err
	}
	m.File = fsSvc

	return nil
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
	wx.Handler `route:"/hello/{name}/{LangCode}"`
	Name       string
	LangCode   string
}) (string, error) {
	if h.Name == "" {
		return "", wx.Errors.RequireErr("name")
	}
	if h.LangCode == "" {
		return "", wx.Errors.RequireErr("langCode")
	}
	if h.LangCode == "vn" {
		return t.DoSayHelloInVn(h.Name), nil
	}
	if h.LangCode == "en" {
		return t.DoSayHelloInEn(h.Name), nil
	}
	return "", wx.Errors.UnSupportError(fmt.Sprintf("%s is not support", h.LangCode))
}
