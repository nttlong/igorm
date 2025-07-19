package vdb

import (
	"fmt"
	"sync"
	"vdb/migrate"
	"vdb/tenantDB"

	_ "github.com/microsoft/go-mssqldb"
)

type TenantDB struct {
	*tenantDB.TenantDB
	dsn string
}

func Open(driverName, dns string) (*TenantDB, error) {
	ret, err := tenantDB.Open(driverName, dns)
	if err != nil {
		return nil, err
	}
	migrator, err := NewMigrator(ret)
	if err != nil {
		return nil, err
	}
	err = migrator.DoMigrates()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &TenantDB{TenantDB: ret, dsn: dns}, nil

}

type initCreateDBNoCache struct {
	once sync.Once
	dsn  string
	err  error
}

var cacheCreateDB = sync.Map{}

func (db *TenantDB) createDBNoCache(dbName string) (string, error) {
	key := fmt.Sprintf("%s:%s", dbName, db.GetDbType())
	actual, _ := cacheCreateDB.LoadOrStore(key, &initCreateDBNoCache{})
	init := actual.(*initCreateDBNoCache)
	init.once.Do(func() {

		dialect := dialectFactory.create(db.GetDriverName())
		dsn, err := dialect.NewDataBase(db.DB, db.dsn, dbName)
		init.dsn = dsn
		init.err = err
	})
	return init.dsn, init.err

}
func (db *TenantDB) CreateDB(dbName string) (*TenantDB, error) {
	dialect := dialectFactory.create(db.GetDriverName())
	dsn, err := dialect.NewDataBase(db.DB, db.dsn, dbName)
	if err != nil {
		return nil, err
	}
	ret, err := Open(db.GetDriverName(), dsn)
	return ret, err
}

//	type Model struct {
//		migrate.Entity
//	}

var ModelRegistry = migrate.ModelRegistry

func NewMigrator(db *tenantDB.TenantDB) (migrate.IMigrator, error) {
	return migrate.NewMigrator(db)
}
func NewMigrator2(db *tenantDB.TenantDB) (migrate.IMigrator, error) {
	return migrate.NewMigrator2(db)
}
