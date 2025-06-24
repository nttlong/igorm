package config

import (
	"dbx"
	"sync"
)

var onceConfig sync.Once
var cfg *dbx.Cfg

func CreateCfg() *dbx.Cfg {
	onceConfig.Do(func() {
		cfg = &dbx.Cfg{
			Driver:         AppConfigInstance.Database.Driver,
			Host:           AppConfigInstance.Database.Host,
			Port:           AppConfigInstance.Database.Port,
			User:           AppConfigInstance.Database.User,
			Password:       AppConfigInstance.Database.Password,
			SSL:            AppConfigInstance.Database.SSL,
			DbName:         AppConfigInstance.Database.Name,
			IsMultiTenancy: AppConfigInstance.Database.IsMultiTenancy,
		}

	})
	return cfg
}
func CreateDbx() *dbx.DBX {
	ret := dbx.NewDBX(*CreateCfg())
	return ret
}
func CreateTenantDbx(tenant string) (*dbx.DBXTenant, error) {
	db := CreateDbx()
	r, e := db.GetTenant(tenant)
	if e != nil {
		return nil, e
	}
	e = db.Open()
	if e != nil {
		return nil, e
	}
	return r, nil
}
