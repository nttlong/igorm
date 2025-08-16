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
	File multipart.FileHeader
}) (*UploadResult, error) {
	// if data.File == nil {
	// 	return nil, wx.Errors.RequireErr("file")
	// }
	return &UploadResult{}, nil
}
