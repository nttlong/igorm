package eorm

import (
	"eorm/tenantDB"
)

func Open(driverName, dataSourceName string) (*tenantDB.TenantDB, error) {
	return tenantDB.Open(driverName, dataSourceName)
}
