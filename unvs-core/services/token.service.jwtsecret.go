package services

import (
	"dbx"
	"fmt"
	"reflect"
	"sync"

	coreModels "unvs.core/models"
)

var cacheGetJwtSecret sync.Map

func (p *TokenService) GetJwtSecret() ([]byte, error) {
	path := p.getPath()
	key := path + "://GetJwtSecret/v2" + p.TenantDb.TenantDbName
	if v, ok := cacheGetJwtSecret.Load(key); ok {
		return v.([]byte), nil
	}

	jwtSecret, err := p.getJwtSecret()
	if err != nil {
		return nil, err
	}
	cacheGetJwtSecret.Store(key, jwtSecret)
	return jwtSecret, nil
}
func (p *TokenService) getJwtSecret() ([]byte, error) {
	if p.EncryptionKey == "" {
		pkgPath := reflect.TypeOf(*p).PkgPath() + "/getJwtSecret"
		panic(fmt.Sprintf("encryption key is missing in %s", pkgPath))
	}

	jwtSecret, err := p.generateRandomSecret(255)
	if err != nil {
		return nil, err
	}
	encryptedJwtSecret, err := p.encryptBytes(p.EncryptionKey, []byte(jwtSecret))

	if err != nil {
		return nil, err
	}
	err = p.TenantDb.Insert(&coreModels.AppConfig{
		Name:      p.TenantDb.TenantDbName,
		Tenant:    p.TenantDb.TenantDbName,
		AppId:     dbx.NewUUID(),
		JwtSecret: *encryptedJwtSecret,
	})
	if err != nil {
		if dbxErr, ok := err.(*dbx.DBXError); ok {
			if dbxErr.Code == dbx.DBXErrorCodeDuplicate {
				if dbxErr.Fields[0] == "Tenant" || dbxErr.Fields[0] == "Name" {
					qr := dbx.Query[coreModels.AppConfig](p.TenantDb, p.Context).Where("Tenant =?", p.TenantDb.TenantDbName)
					qr.Select("JwtSecret")
					appConfig, err := qr.First()
					if err != nil {
						return nil, fmt.Errorf("failed to get jwt secret: %w", err)
					}
					if appConfig.JwtSecret == "" {
						return nil, fmt.Errorf("failed to get jwt secret: %w", err)
					}
					// Decrypt jwt secret

					return p.decryptBytes(p.EncryptionKey, appConfig.JwtSecret)

				}
			}
		}
		return nil, err
	}

	return []byte(jwtSecret), nil
}
