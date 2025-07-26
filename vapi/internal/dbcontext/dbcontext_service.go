package dbcontext

import (
	"vapi/internal/config"
	"vdb"
)

type DbContext struct {
	DB *vdb.TenantDB
}

// Khởi tạo DbContext từ ConfigService
func NewDbContext(cfg *config.ConfigService) (*DbContext, error) {
	dsn := cfg.Get().Database.Dsn
	driver := cfg.Get().Database.Driver

	db, err := vdb.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return &DbContext{DB: db}, nil
}
