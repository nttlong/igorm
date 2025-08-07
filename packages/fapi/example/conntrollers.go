package example

import (
	"fapi"
	"mime/multipart"
)

type Media struct {
}

// func (m *Media) Upload(ctx *fapi.Context, file multipart.FileHeader) (*UploadResult, error) {
// 	panic("implement me")
// }

type MediaContext struct {
	fapi.Context `route:"uri:{FileID}.mp4;method:GET"`
	FileID       string `param:"FileID"`
}

func (m *Media) Content(ctx *MediaContext) ([]byte, error) {
	panic("implement me")
}

type UploadResult struct {
	FileID string `json:"file_id"`
	Name   string `json:"name"`
}

func (m *Media) Upload(ctx *fapi.Context, file multipart.FileHeader) (*UploadResult, error) {
	panic("implement me")
}

type DeptInfo struct {
	Code string //<-- 1 text box
}
type EmpInfo struct { // swager se render cau truc nay
	Name string //<-- 1 text bo
	Depr DeptInfo
}

func (m *Media) Register(
	ctx *struct {
		fapi.Context `route:"uri:/emp-ergister/{Name}" `
		Name         string                `json:"name"`
		Age          int                   `json:"age"`
		File         *multipart.FileHeader `json:"file"`
	},
) (*UploadResult, error) {
	panic("implement me")
}
