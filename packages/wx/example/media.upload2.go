package example

import (
	"fmt"
	"mime/multipart"
	"wx"
)

type FileUtils struct {
}

func (fs *FileUtils) SaveFile() {

}

type FileUtilsService struct {
	FileUtil *wx.Scoped[FileUtils]
}

func (f *FileUtilsService) New() error {
	f.FileUtil.Init(func(ctx *wx.ServiceContext) (*FileUtils, error) {
		return &FileUtils{}, nil
	})

	return nil
}

func (m *Media) Upload2(ctx *struct {
	wx.Handler `route:"method:post;uri:@/{Tenant}"`
	Tenant     string
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
