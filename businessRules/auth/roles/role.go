package auth

import (
	"context"

	authErr "unvs.br.auth/errors"

	"dbx"
	"dynacall"
	_ "dynacall"

	userUtils "unvs.br.auth/users"
)

type RoleService struct {
	dynacall.Caller
	Tenant      string
	TenantDb    *dbx.DBXTenant
	Context     context.Context
	AccessToken string
	JwtSecret   []byte
	tokenInfo   *userUtils.OAuth2Token
}

func (s *RoleService) validateAccessToken() error {
	// TODO: validate access token
	accessToken := s.AccessToken
	if accessToken == "" {
		e := &authErr.AuthError{
			Code:    authErr.ErrInvalidToken,
			Message: "access token is missing",
		}
		return e
	}
	tokenInfo, err := userUtils.DecodeAccessToken(s.JwtSecret, accessToken)
	if err != nil {
		return err
	}
	s.tokenInfo = tokenInfo

	return nil
}
func init() {
	dynacall.RegisterCaller(&RoleService{
		Caller: dynacall.Caller{
			Path: "auth",
		},
	})
}
