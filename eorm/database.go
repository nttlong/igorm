package eorm

import (
	"eorm/migrate"
	"eorm/tenantDB"

	_ "github.com/microsoft/go-mssqldb"
)

func Open(driverName, dns string) (*tenantDB.TenantDB, error) {

	return tenantDB.Open(driverName, dns)
}

//	type Model struct {
//		migrate.Entity
//	}

var ModelRegistry = migrate.ModelRegistry

func NewMigrator(db *tenantDB.TenantDB) (migrate.IMigrator, error) {
	return migrate.NewMigrator(db)
}
