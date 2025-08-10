package example

import (
	"fmt"
	"mime/multipart"
	"vapi"

	"github.com/google/uuid"
)

type Media struct {
}

type AuthHandler struct {
	vapi.Handler
	// Auth *vapi.AuthClaims
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
	Files    []*multipart.FileHeader `json:"file"`
	FileName string                  `json:"file_name"`
	UploadId string                  `json:"upload_id"`
	Info     struct {
		FolderId string `json:"folder_id"`
	}
}) {
	fmt.Println(data.Info)
}
