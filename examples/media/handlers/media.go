package handlers

import (
	"fmt"
	"media/services"
	"wx"
)

type Media struct {
	wx.ControllerContext
	File           services.FileService
	Directories    services.DirectoryService
	FileSvc        *wx.Depend[services.FileService, Media]
	DirectoriesSvc *wx.Depend[services.DirectoryService, Media]
}

func (m *Media) New(fs *wx.Depend[services.FileService, Media], ds *wx.Global[services.DirectoryService]) error {
	fmt.Println("Media.New")
	ds.Init(func() (*services.DirectoryService, error) {
		return &services.DirectoryService{}, nil
	})
	test, err := ds.Ins()
	if err != nil {
		return err
	}
	fmt.Println(test)
	m.DirectoriesSvc.Init(func(app *Media) (*services.DirectoryService, error) {
		ret := &services.DirectoryService{}
		ret.New()

		return ret, nil

	})
	m.FileSvc.Init(func(app *Media) (*services.FileService, error) {
		dirSvc, err := app.DirectoriesSvc.Ins()
		if err != nil {
			return nil, err
		}
		return &services.FileService{
			DirectorySvc: dirSvc,
		}, nil

	})

	return nil
}
