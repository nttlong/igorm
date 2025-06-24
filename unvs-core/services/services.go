package services

// import (
// 	"context"
// 	"dbx"

// 	cacher "unvs.core/cacher"
// )

// type factory struct {
// }
// type TokenServiceParams struct {
// 	Cache         cacher.Cache
// 	Context       context.Context
// 	TenantDb      *dbx.DBXTenant
// 	EncryptionKey string
// 	Language      string
// }
// type PasswordServiceParams struct {
// 	Cache   cacher.Cache
// 	Context context.Context
// }
// type AuthServiceParams struct {
// 	Cache         cacher.Cache
// 	Context       context.Context
// 	EncryptionKey string
// 	Language      string

// 	TenantDb *dbx.DBXTenant
// }

// func (f *factory) GetTokenService(params TokenServiceParams) *TokenService {

// 	return &TokenService{
// 		Cache:         params.Cache,
// 		Context:       params.Context,
// 		TenantDb:      params.TenantDb,
// 		EncryptionKey: params.EncryptionKey,
// 		Language:      params.Language,
// 	}
// }
// func (f *factory) GetPasswordService(params PasswordServiceParams) *PasswordService {
// 	return &PasswordService{
// 		Cache:   params.Cache,
// 		Context: params.Context,
// 	}
// }
// func (f *factory) GetAuthService(params AuthServiceParams) *AuthenticateService {
// 	return &AuthenticateService{
// 		TokenService: TokenService{
// 			Cache:         params.Cache,
// 			Context:       params.Context,
// 			TenantDb:      params.TenantDb,
// 			EncryptionKey: params.EncryptionKey,
// 			Language:      params.Language,
// 		},
// 		PasswordService: PasswordService{
// 			Cache:   params.Cache,
// 			Context: params.Context,
// 		},
// 		Context: params.Context,
// 		Cache:   params.Cache,
// 	}
// }

// var Factory = &factory{}
