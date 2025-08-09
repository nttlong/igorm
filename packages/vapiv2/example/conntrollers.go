package example

import (
	"vapi"
)

type Media struct {
}

type AuthHandler struct {
	vapi.Handler
	Auth *vapi.AuthClaims
}
type UploadResult struct {
}

func (m *Media) Register(
	ctx *struct {
		AuthHandler
	},
) (*UploadResult, error) {
	panic("implement me")
}
