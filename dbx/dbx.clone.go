package dbx

import (
	"fmt"
	"sync"
)

var (
	cacheDBXTenantClones = sync.Map{}
)

func (db *DBXTenant) Clone(instanceName string) *DBXTenant {
	key := fmt.Sprintf("%s_%s", db.TenantDbName, instanceName)
	if v, ok := cacheDBXTenantClones.Load(key); ok {
		return v.(*DBXTenant)
	}
	// Tạo clone mới
	// Đọc cache với lock đọc (RWMutex)

	// Tạo clone mới
	ret := &DBXTenant{
		DBX: DBX{
			cfg:      db.cfg,
			dns:      db.dns,
			executor: db.executor,
			compiler: db.compiler,
			isOpen:   db.isOpen,
		},
		TenantDbName: db.TenantDbName,
	}
	ret.Open()

	// Lưu vào cache
	cacheDBXTenantClones.Store(key, ret)

	return ret
}
