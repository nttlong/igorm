package example

import (
	"mime/multipart"
	"vapi"
	vapiFileUtils "vapi/fileUtils"
)

type FileUtilsInjector struct {
	vapi.Inject[vapiFileUtils.FileUtils]
}

func (m *Media) Upload2(ctx *struct {
	vapi.Handler `route:"method:post;uri:@/{Tenant}"`
	Tenant       string
}, data struct {
	File multipart.FileHeader
}, fileUtils vapi.Inject[vapiFileUtils.FileUtils]) (UploadResult, error) {
	return UploadResult{}, nil
}
