package repo

import (
	"vdb"
	"xauth/config"
)

type DbContext struct {
	db *vdb.TenantDB
}

func NewDbContext(cfg config.ConfigService) (*DbContext, error) {
	ret := &DbContext{}
	db, err := vdb.Open("postgres", cfg.Get().Database.Postgres.Dsn)
	if err != nil {
		return nil, err
	}
	ret.db = db
	return ret, nil

}
func (dbContext *DbContext) GetTenantDb(tenantName string) (*vdb.TenantDB, error) {
	return dbContext.db.CreateDB(tenantName)

}
