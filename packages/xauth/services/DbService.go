package services

import (
	"vdb"
	"wx"
	"xconfig"
)

type DbService struct {
	cfg *xconfig.Config
	Db  *vdb.TenantDB
}

func (s *DbService) New(cfg *wx.Global[ConfigService]) error {
	var err error
	config, err := cfg.Ins()
	if err != nil {
		return err
	}
	s.cfg = config.Data

	vdb.SetManagerDb(s.cfg.Database.Driver, s.cfg.Database.Postgres.Manager)
	s.Db, err = vdb.Open(s.cfg.Database.Driver, s.cfg.Database.Postgres.Dsn)
	if err != nil {
		return err
	}

	return nil
}
func (s *DbService) GetTenantDb(dbName string) (*vdb.TenantDB, error) {
	return s.Db.CreateDB(dbName)
}
