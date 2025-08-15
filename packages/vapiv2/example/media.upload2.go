package example

import (
	"fmt"
	"mime/multipart"
	"vapi"
	vapiFileUtils "vapi/fileUtils"
)

type FileUtilsService struct {
	FileUtil *vapi.Scoped[vapiFileUtils.FileUtils]
}

func (f *FileUtilsService) New() error {
	f.FileUtil.Init(func(ctx *vapi.ServiceContext) (*vapiFileUtils.FileUtils, error) {
		fmt.Println("OK")
		return &vapiFileUtils.FileUtils{}, nil
	})

	return nil
}

func (m *Media) Upload2(ctx *struct {
	vapi.Handler `route:"method:post;uri:@/{Tenant}"`
	Tenant       string
}, data struct {
	File multipart.FileHeader
}, fileUtils *FileUtilsService) (UploadResult, error) {
	files, err := fileUtils.FileUtil.GetInstance()
	if err != nil {
		return UploadResult{}, err
	}
	files.SaveFile()

	fmt.Println(files)
	return UploadResult{}, nil
}
