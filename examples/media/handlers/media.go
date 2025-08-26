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
