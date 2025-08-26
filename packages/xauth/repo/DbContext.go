package repo

import (
	"vdb"
	"wx"
	"xauth/config"
)

type DbContext struct {
	db *vdb.TenantDB
}

func (dbContext *DbContext) New(
	ctx *wx.HttpContext,
	cfgInject *wx.Global[config.ConfigService]) error {
	if dbContext.db != nil {
		return nil
	}
	var err error
	var cfgService config.ConfigService
	cfgService, err = cfgInject.Ins()

	if err != nil {
		return err
	}
	vdb.SetManagerDb("postgres", cfgService.Get().Database.Postgres.Manager)

	dbContext.db, err = vdb.Open("postgres", cfgService.Get().Database.Postgres.Dsn)
	if err != nil {
		return err
	}
	return nil
}
func (dbContext *DbContext) GetTenantDb(tenantName string) (*vdb.TenantDB, error) {
	return dbContext.db.CreateDB(tenantName)

}
