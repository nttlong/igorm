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
type Helper struct {
}

func (media *Media) Upload(ctx *wx.Handler, data struct {
	File *multipart.FileHeader
}, helper *wx.Depend[Helper]) (*UploadResult, error) {
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

	fileService.SaveFile(data.File)

	return &UploadResult{}, nil
}
