package handlers

import (
	"media/services"
	"wx"
)

type ListFile3 struct {
	wx.Handler
	FileManager services.FileManager
}

func (lst *ListFile3) New() error {
	lst.FileManager = services.NewFileManagerLocal()

	return nil
}
func (media *Media) ListFiles3(ctx *struct {
	ListFile3 `route:"method:get"`
}) ([]string, error) {
	uriOfFile, err := wx.GetUriOfHandler[Media]("Files")
	if err != nil {
		return nil, err
	}
	// 2. Base URL
	baseUrl := ctx.GetAbsRootUri() + "/api" + uriOfFile + "/"
	lst, err := ctx.FileManager.GetFileList()
	for i := 0; i < len(lst); i++ {
		lst[i] = baseUrl + lst[i]
	}
	return lst, err
}

type ListFile4 struct {
	wx.HttpContext
}
