package config

import (
	"dbx"
)

func createCfg() *dbx.Cfg {
	return &dbx.Cfg{
		Driver:   "mssql",
		Host:     "localhost",
		Port:     0,
		User:     "sa",
		Password: "123456",
		SSL:      false,
	}
}
func createDbx() *dbx.DBX {
	ret := dbx.NewDBX(*createCfg())
	return ret
}
func CreateTenantDbx(tenant string) (*dbx.DBXTenant, error) {
	db := createDbx()
	r, e := db.GetTenant(tenant)
	if e != nil {
		return nil, e
	}
	db.Open()
	return r, nil
}
