package handlers

import (
	"wx"
)

func (media *Media) ListFiles(ctx *wx.Handler) (*[]string, error) {

	ret, err := media.File.ListAllFiles()
	if err != nil {
		return nil, err
	}

	return &ret, nil

}
