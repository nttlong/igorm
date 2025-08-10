package example

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"vapi"

	"github.com/google/uuid"
)

type Media struct {
}

type AuthHandler struct {
	vapi.Handler
	//Auth *vapi.AuthClaims
}
type UploadResult struct {
	UploadId string
}

func (m *Media) Register(
	ctx *struct {
		AuthHandler
	},
	data struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	},
) (*UploadResult, error) {
	return &UploadResult{
		UploadId: uuid.New().String(),
	}, nil
}
func (m *Media) Upload(ctx *AuthHandler, data struct {
	Files []*multipart.FileHeader `json:"file"`

	Info struct {
		FolderId string `json:"folder_id"`
	}
}, auth *vapi.AuthClaims) ([]string, error) {
	if data.Files == nil {
		return nil, fmt.Errorf("file is required")
	}
	ret := []string{}
	for _, file := range data.Files {
		uploadDir := "./uploads/"

		// Tạo thư mục nếu chưa tồn tại
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("không tạo được thư mục upload: %w", err)
		}

		f, err := file.Open() // file là *multipart.FileHeader
		if err != nil {
			return nil, err
		}
		defer f.Close()

		// Tạo file đích
		out, err := os.Create(filepath.Join(uploadDir, file.Filename))
		if err != nil {
			return nil, err
		}
		defer out.Close()

		// Copy dữ liệu
		if _, err = io.Copy(out, f); err != nil {
			return nil, err
		}

	}
	return ret, nil
}
func (m *Media) File(ctx struct {
	AuthHandler `route:"method:get;uri:@/{FileName}"`
	FileName    string
}) error {
	fileName := "./uploads/" + ctx.FileName
	return ctx.StreamingFile(fileName)

}
