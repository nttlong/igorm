package handlers

import (
	"fmt"
	"wx"
)

func (media *Media) ListFiles(ctx *wx.Handler) (*[]string, error) {
	fs, err := media.FileSvc.Ins()
	if err != nil {
		return nil, err
	}
	fmt.Println(fs)

	// fmt.Println(media.BaseUrl)
	// directoryService, err := media.Directories.Ins()
	// if err != nil {
	// 	return nil, err
	// }
	// ret, err := directoryService.ListAllDirectories()
	// if err != nil {
	// 	return nil, err
	// }

	dirs := media.Directories

	ret, err := media.File.ListAllFiles(&dirs)
	if err != nil {
		return nil, err
	}

	return &ret, nil

}
