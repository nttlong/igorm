package dbcontext

import (
	"vapi/internal/config"
	"vdb"
)

type DbContext struct {
	DB *vdb.TenantDB
}

func (db *DbContext) New(driver string, dsn string) error {
	var err error
	db.DB, err = vdb.Open(driver, dsn)
	if err != nil {
		return err
	}

	return nil
}

// Khởi tạo DbContext từ ConfigService
func NewDbContext(cfg *config.ConfigService) (*DbContext, error) {
	dsn := cfg.Get().Database.Dsn
	driver := cfg.Get().Database.Driver
	vdb.SetManagerDb(driver, cfg.Get().Database.Manager)

	db, err := vdb.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return &DbContext{DB: db}, nil
}
