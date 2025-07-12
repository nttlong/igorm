package eorm

import (
	"eorm/migrate"
	"eorm/tenantDB"
)

func Open(driverName, dataSourceName string) (*tenantDB.TenantDB, error) {
	return tenantDB.Open(driverName, dataSourceName)
}

type Model struct {
	migrate.Entity
}

var ModelRegistry = migrate.ModelRegistry
