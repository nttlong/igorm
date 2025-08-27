package handlers

import (
	"media/services"
	"os"
	"path/filepath"
	"strings"
	"wx"
)

//const rootPath = `D:\code\go\news2\igorm\examples\media\cmd\uploads`

type DirectoryReader struct {
	Files []string
	Dirs  []string
}
type UserInfo struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
}

func (user *UserInfo) New(ctx *wx.AuthContext) error {
	// Initialize user info if needed
	return nil
}

type RoleChecker struct {
}

func (role *RoleChecker) New(ctx *wx.AuthContext) error {
	// Implement role checking logic here
	return nil
}
func (media *Media) ListFiles(ctx *struct {
	wx.Handler `route:"method:get"`
}, dr *wx.Depend[DirectoryReader], userSvc *wx.Auth[UserInfo]) (*[]string, error) {

	// 1. Kiểm tra thư mục gốc
	if _, err := os.Stat(media.FileDirectory); os.IsNotExist(err) {
		// http.Error(ctx.Res, "Directory not found.", http.StatusNotFound)
		return nil, err
	}
	uriOfFile, err := wx.GetUriOfHandler[Media]("Files")
	if err != nil {
		return nil, err
	}
	// 2. Base URL
	baseUrl := ctx.GetAbsRootUri() + "/api" + uriOfFile + "/"
	results := make([]string, 0, 1024) // pre-allocate

	// 3. WalkDir thay vì Walk
	err = filepath.WalkDir(media.FileDirectory, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == media.FileDirectory {
			return nil
		}

		// relative path
		rel, _ := filepath.Rel(media.FileDirectory, path)
		urlPath := strings.ReplaceAll(rel, "\\", "/")

		results = append(results, baseUrl+urlPath)
		return nil
	})

	if err != nil {

		return nil, err
	}

	return &results, nil

}

type FileManagerService struct {
	wx.Service
	FileManager services.FileManager
}

func (fm *FileManagerService) New() error {
	fm.FileManager = services.NewFileManagerLocal()
	return nil
}

func (media *Media) ListFiles2(ctx *struct {
	wx.Handler `route:"method:get"`
}, fileManager *FileManagerService) ([]string, error) {
	uriOfFile, err := wx.GetUriOfHandler[Media]("Files")
	if err != nil {
		return nil, err
	}
	// 2. Base URL
	baseUrl := ctx.GetAbsRootUri() + "/api" + uriOfFile + "/"
	lst, err := fileManager.FileManager.GetFileList()
	for i := 0; i < len(lst); i++ {
		lst[i] = baseUrl + lst[i]
	}
	return lst, err

}
