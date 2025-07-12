package migrate

import (
	"eorm/tenantDB"
	"reflect"
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

func NewMigrator(db *tenantDB.TenantDB) (IMigrator, error) {
	err := db.Detect()
	if err != nil {
		return nil, err
	}
	switch db.DbType {
	case tenantDB.DB_DRIVER_MSSQL:
		loader, err := MigratorLoader(db)
		if err != nil {
			return nil, err
		}
		return &migratorMssql{
			db:     db,
			loader: loader,
		}, nil
	default:
		panic("unsupported database type")
	}
}
