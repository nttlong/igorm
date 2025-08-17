package handlers

import (
	"fmt"
	"media/services"
	"wx"
)

type Media struct {
	wx.ControllerContext
	File        *services.FileService
	Directories *services.DirectoryService
}

func (m *Media) New(
	fs *wx.Depend[services.FileService, Media],
	ds *wx.Global[services.DirectoryService],
	urlSvc *wx.Depend[services.UrlResolverService, Media],
) error {
	fmt.Println("Media.New")
	ds.Init(func() (*services.DirectoryService, error) {
		return &services.DirectoryService{}, nil
	})

	ds.Init(func() (*services.DirectoryService, error) {
		return &services.DirectoryService{
			DirUpload:     "./uploads",
			DirUploadName: "uploads",
		}, nil
	})
	urlSvc.Init(func(app *Media) (*services.UrlResolverService, error) {
		return &services.UrlResolverService{
			BaseUrl: m.BaseUrl,
		}, nil
	})
	fs.Init(func(app *Media) (*services.FileService, error) {
		dirSvc, err := ds.Ins()
		if err != nil {
			return nil, err
		}
		urlSvc, err := urlSvc.Ins()
		if err != nil {
			return nil, err
		}

		return &services.FileService{
			DirectorySvc: dirSvc,
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
