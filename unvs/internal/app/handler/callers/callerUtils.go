package caller

import (
	"sync"

	"dbx"
)

func callerUtilsSetJWTInfo() {

}
func getCfg() dbx.Cfg {
	ret := dbx.Cfg{
		Driver: "mssql",
		Host:   "localhost",

		User:     "sa",
		Password: "123456",
		SSL:      false,
	}
	return ret
}

var db dbx.DBX
var once sync.Once

func getDbx() dbx.DBX {
	once.Do(func() {
		cfg := getCfg()
		_db := dbx.NewDBX(cfg)
		db = *_db
	})

	return db
}

var cacheDbTenants = sync.Map{}

func getTenantDb(tenant string) (*dbx.DBXTenant, error) {
	// check if tenant is in cache
	if val, ok := cacheDbTenants.Load(tenant); ok {
		ret := val.(dbx.DBXTenant)
		return &ret, nil
	}
	// create new dbx for tenant
	fx := getDbx()
	fx.Open()
	defer fx.Close()

	ret, err := fx.GetTenant(tenant)
	if err != nil {
		return nil, err
	}
	ret.Open()
	cacheDbTenants.Store(tenant, *ret)
	return ret, nil

}
