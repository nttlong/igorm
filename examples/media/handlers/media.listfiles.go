package handlers

import (
	"os"
	"path/filepath"
	"strings"
	"wx"
)

const rootPath = `D:\code\go\news2\igorm\examples\media\cmd\uploads`

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

func (media *Media) ListFiles(ctx *struct {
	wx.Handler `route:"method:get"`
}, dr *wx.Depend[DirectoryReader], userSvc *wx.Auth[UserInfo]) (*[]string, error) {

	// 1. Kiểm tra thư mục gốc
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		// http.Error(ctx.Res, "Directory not found.", http.StatusNotFound)
		return nil, err
	}

	// 2. Base URL
	baseUrl := "http://" + ctx.Req.Host + ctx.Req.URL.Path
	results := make([]string, 0, 1024) // pre-allocate

	// 3. WalkDir thay vì Walk
	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == rootPath {
			return nil
		}

		// relative path
		rel, _ := filepath.Rel(rootPath, path)
		urlPath := strings.ReplaceAll(rel, "\\", "/")

		results = append(results, baseUrl+"/"+urlPath)
		return nil
	})

	if err != nil {

		return nil, err
	}

	return &results, nil

}
func (media *Media) Hello(ctx *struct {
	wx.Handler `route:"method:get"`
}) (string, error) {

	return "Hello World", nil

}
