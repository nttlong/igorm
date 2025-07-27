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
	if ret.GetDBName() != "" {
		migrator, err := NewMigrator(ret)
		if err != nil {
			return nil, err
		}
		err = migrator.DoMigrates()
		if err != nil {
			return nil, err
		}
	}

	return &TenantDB{TenantDB: ret, dsn: dns}, nil

}

type initCreateDBNoCache struct {
	once sync.Once
	dsn  string
	err  error
	db   *TenantDB
}

var cacheCreateDB = sync.Map{}

func (db *TenantDB) createDBNoCache(dbName string) (string, *TenantDB, error) {
	key := fmt.Sprintf("%s:%s", dbName, db.GetDbType())
	actual, _ := cacheCreateDB.LoadOrStore(key, &initCreateDBNoCache{})
	init := actual.(*initCreateDBNoCache)
	init.once.Do(func() {

		dialect := dialectFactory.create(db.GetDriverName())
		dsn, err := dialect.NewDataBase(db.DB, db.dsn, dbName)
		if err != nil {
			init.err = err
			return
		}
		init.dsn = dsn
		ret, err := Open(db.GetDriverName(), dsn)
		if err != nil {
			init.err = err
			return
		}
		init.db = ret

	})
	return init.dsn, init.db, init.err

}

func (db *TenantDB) CreateDB(dbName string) (*TenantDB, error) {
	_, tenantDb, err := db.createDBNoCache(dbName)
	if err != nil {
		return nil, err
	}
	return tenantDb, nil

	// dialect := dialectFactory.create(db.GetDriverName())
	// dsn, err := dialect.NewDataBase(db.DB, db.dsn, dbName)
	// if err != nil {
	// 	return nil, err
	// }
	// ret, err := Open(db.GetDriverName(), dsn)
	// if err != nil {
	// 	return nil, err
	// }

	// return ret, nil
}

//	type Model struct {
//		migrate.Entity
//	}

var ModelRegistry = migrate.ModelRegistry

func RegisterModel(models ...interface{}) {
	ModelRegistry.Add(models...)
}
func NewMigrator(db *tenantDB.TenantDB) (migrate.IMigrator, error) {
	return migrate.NewMigrator(db)
}
func NewMigrator2(db *tenantDB.TenantDB) (migrate.IMigrator, error) {
	return migrate.NewMigrator2(db)
}

type tenantDbManager struct {
	mapManagerDb map[string]bool
}

func (t *tenantDbManager) SetManagerDb(driver string, dbName string) {
	t.mapManagerDb[dbName+"://"+driver] = true
}
func (t *tenantDbManager) isManagerDb(driver string, dbName string) bool {
	if _, ok := t.mapManagerDb[dbName+"://"+driver]; ok {
		return true
	}
	return false
}
func SetManagerDb(driver string, dbName string) {
	tenantDbManagerInstance.SetManagerDb(driver, dbName)
}

type dbFunCall struct {
	expr string
	args []interface{}
}

func DbFunCall(expr string, args ...interface{}) dbFunCall {
	return dbFunCall{expr: expr, args: args}
}
func Expr(expr string, args ...interface{}) dbFunCall {
	return dbFunCall{expr: expr, args: args}
}
func (db *TenantDB) LikeValue(value string) string {
	dialect := dialectFactory.create(db.GetDriverName())
	return dialect.LikeValue(value)
}

var tenantDbManagerInstance = &tenantDbManager{mapManagerDb: make(map[string]bool)}

func init() {
	tenantDB.IsManagerDb = tenantDbManagerInstance.isManagerDb
}
