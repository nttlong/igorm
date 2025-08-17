package handlers

import (
	"mime/multipart"
	"wx"
)

type UploadResult struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	FileType string `json:"file_type"`
}

func (media *Media) Upload(ctx *wx.Handler, data struct {
	File *multipart.FileHeader
}) (*UploadResult, error) {
	// if data.File == nil {
	// 	return nil, wx.Errors.RequireErr("file")
	// }
	// file, err := data.File.Open()
	// if err != nil {
	// 	return nil, err
	// }
	fileService := media.File
	// if err != nil {
	// 	return nil, err
	// }

	directoryService := media.Directories

	fileService.SaveFile(data.File, &directoryService)

	return &UploadResult{}, nil
}
