package migrate

import (
	"eorm/tenantDB"
	"fmt"
	"reflect"
	"sync"
)

type IMigrator interface {
	Quote(names ...string) string
	GetColumnDataTypeMapping() map[reflect.Type]string
	GetGetDefaultValueByFromDbTag() map[string]string
	GetSqlCreateTable(entityType reflect.Type) (string, error)
	GetSqlAddColumn(entityType reflect.Type) (string, error)
	GetSqlAddIndex(entityType reflect.Type) (string, error)
	GetSqlAddUniqueIndex(entityType reflect.Type) (string, error)
	GetSqlMigrate(entityType reflect.Type) ([]string, error)
	DoMigrate(entityType reflect.Type) error
	DoMigrates() error
}

type migratorInit struct {
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
