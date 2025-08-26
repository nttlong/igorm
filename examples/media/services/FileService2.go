package services

import (
	"os"
	"path/filepath"
	"strings"
)

type FileManager interface {
	GetRootDirPath() string
	GetFileList() ([]string, error)
}
type FileManagerLocal struct {
	rootDirPath string
}

func (fm *FileManagerLocal) GetRootDirPath() string {
	return fm.rootDirPath
}
func (fm *FileManagerLocal) GetFileList() ([]string, error) {
	var err error
	// 1. Kiểm tra thư mục gốc
	if _, err := os.Stat(fm.GetRootDirPath()); os.IsNotExist(err) {
		// http.Error(ctx.Res, "Directory not found.", http.StatusNotFound)
		return nil, err
	}

	// 2. Base URL

	results := make([]string, 0, 1024) // pre-allocate

	// 3. WalkDir thay vì Walk
	err = filepath.WalkDir(fm.GetRootDirPath(), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == fm.GetRootDirPath() {
			return nil
		}

		// relative path
		rel, _ := filepath.Rel(fm.GetRootDirPath(), path)
		urlPath := strings.ReplaceAll(rel, "\\", "/")

		results = append(results, urlPath)
		return nil
	})

	if err != nil {

		return nil, err
	}

	return results, nil
}
func NewFileManagerLocal() FileManager {
	return &FileManagerLocal{
		rootDirPath: "./uploads",
	}
}
