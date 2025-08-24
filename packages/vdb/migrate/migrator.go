package migrate

import (
	"fmt"
	"reflect"
	"sync"
	"vdb/tenantDB"

	_ "github.com/lib/pq"
)

type IMigrator interface {
	GetLoader() IMigratorLoader
	Quote(names ...string) string
	GetSqlInstallDb() ([]string, error)
	GetColumnDataTypeMapping() map[reflect.Type]string
	GetGetDefaultValueByFromDbTag() map[string]string
	GetSqlCreateTable(entityType reflect.Type) (string, error)
	GetSqlAddColumn(entityType reflect.Type) (string, error)
	GetSqlAddIndex(entityType reflect.Type) (string, error)
	GetSqlAddUniqueIndex(entityType reflect.Type) (string, error)
	GetSqlMigrate(entityType reflect.Type) ([]string, error)
	GetSqlAddForeignKey() ([]string, error)
	GetFullScript() ([]string, error)
	DoMigrate(entityType reflect.Type) error
	DoMigrates() error
}

type migratorInit struct {
	once sync.Once
	val  IMigrator
	err  error
}

// var cacheNewMigrator sync.Map

func NewMigrator2(db *tenantDB.TenantDB) (IMigrator, error) {
	err := db.Detect()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%s:%s", db.GetDBName(), db.GetDbType())

	// Load hoặc khởi tạo mới đối tượng quản lý init
	actual, _ := cacheNewMigrator.LoadOrStore(key, &migratorInit{})

	mi := actual.(*migratorInit)
	mi.once.Do(func() {
		switch db.GetDbType() {
		case tenantDB.DB_DRIVER_MSSQL:
			var loader IMigratorLoader
			loader, mi.err = MigratorLoader(db)
			if mi.err != nil {
				return
			}
			mi.val = &migratorMssql{
				db:     db,
				loader: loader,
			}
		default:
			mi.err = fmt.Errorf("unsupported database type: %s", db.GetDbType())
		}
	})

	return mi.val, mi.err
}

type initNewMigrator struct {
	once sync.Once
	val  IMigrator
	err  error
}

var cacheNewMigrator sync.Map

func NewMigrator(db *tenantDB.TenantDB) (IMigrator, error) {
	err := db.Detect()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%s:%s", db.GetDBName(), db.GetDbType())
	actual, _ := cacheNewMigrator.LoadOrStore(key, &initNewMigrator{})
	mi := actual.(*initNewMigrator)
	mi.once.Do(func() {

		loader, err := MigratorLoader(db)
		if err != nil {
			mi.err = err
		}
		switch db.GetDbType() {
		case tenantDB.DB_DRIVER_MSSQL:

			mi.val = &migratorMssql{
				db:     db,
				loader: loader,
			}
		case tenantDB.DB_DRIVER_Postgres:
			mi.val = &migratorPostgres{
				db:     db,
				loader: loader,
			}
		case tenantDB.DB_DRIVER_MySQL:
			mi.val = &migratorMySql{
				db:     db,
				loader: loader,
			}
		default:
			mi.err = fmt.Errorf("unsupported database type: %s", db.GetDbType())
		}
	})
	return mi.val, mi.err
}
