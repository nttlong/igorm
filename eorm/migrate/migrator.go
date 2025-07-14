package migrate

import (
	"eorm/tenantDB"
	"fmt"
	"reflect"
	"sync"

	_ "github.com/lib/pq"
	"golang.org/x/sync/singleflight"
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

var (
	cacheNewMigrator sync.Map
	migratorGroup    singleflight.Group
)

func NewMigrator(db *tenantDB.TenantDB) (IMigrator, error) {
	err := db.Detect()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%s:%s", db.GetDBName(), db.GetDbType())

	// 1. Check cache trước
	if v, ok := cacheNewMigrator.Load(key); ok {
		return v.(IMigrator), nil
	}

	// 2. Dùng singleflight để tránh gọi trùng
	v, err, _ := migratorGroup.Do(key, func() (interface{}, error) {
		// Check cache lần nữa trong group (để tránh race)
		if v, ok := cacheNewMigrator.Load(key); ok {
			return v, nil
		}

		var ret IMigrator
		loader, err := MigratorLoader(db)
		if err != nil {
			return nil, err
		}
		switch db.GetDbType() {
		case tenantDB.DB_DRIVER_MSSQL:
			ret = &migratorMssql{
				db:     db,
				loader: loader,
			}
		case tenantDB.DB_DRIVER_Postgres:
			ret = &migratorPostgres{
				db:     db,
				loader: loader,
			}
		default:
			return nil, fmt.Errorf("unsupported database type: %s", db.GetDbType())
		}

		cacheNewMigrator.Store(key, ret)
		return ret, nil
	})

	if err != nil {
		return nil, err
	}

	return v.(IMigrator), nil
}
