package example

import (
	"mime/multipart"
	"vapi"
	vapiFileUtils "vapi/fileUtils"
)

type FileUtilsService struct {
	FileUtil vapi.Singleton[vapiFileUtils.FileUtils, FileUtilsService]
}

func (m *Media) Upload2(ctx *struct {
	vapi.Handler `route:"method:post;uri:@/{Tenant}"`
	Tenant       string
}, data struct {
	File multipart.FileHeader
}, fileUtils *FileUtilsService) (UploadResult, error) {
	return UploadResult{}, nil
}
