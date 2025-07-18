package dbv

import (
	"dbv/migrate"
	"dbv/tenantDB"

	_ "github.com/microsoft/go-mssqldb"
)

type TenantDB struct {
	*tenantDB.TenantDB
}

func Open(driverName, dns string) (*TenantDB, error) {
	ret, err := tenantDB.Open(driverName, dns)
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
	return &TenantDB{TenantDB: ret}, nil

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
