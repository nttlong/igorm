package tenant

import (
	"sync"
	"vdb"
)

type TenantService struct {
	Db      *vdb.TenantDB
	manager string
}

func (svc *TenantService) New(db *vdb.TenantDB, manager string) error {
	svc.Db = db
	svc.manager = manager
	vdb.SetManagerDb(db.GetDriverName(), manager)
	return nil
}
func NewTenantService(db *vdb.TenantDB, manager string) *TenantService {
	svc := &TenantService{
		Db:      db,
		manager: manager,
	}
	vdb.SetManagerDb(db.GetDriverName(), manager)
	return svc

}

type initTenantServiceTenant struct {
	Db *vdb.TenantDB

	manager string
	err     error
	once    sync.Once
}

var initTenantServiceTenantCache sync.Map

func (svc *TenantService) Tenant(tenantName string) (*vdb.TenantDB, error) {
	key := svc.Db.GetDriverName() + ":" + tenantName
	actual, _ := initTenantServiceTenantCache.LoadOrStore(key, &initTenantServiceTenant{})
	item := actual.(*initTenantServiceTenant)
	item.once.Do(func() {
		tenant, err := svc.Db.CreateDB(tenantName)

		item.Db = tenant
		item.err = err
	})
	return item.Db, item.err

}
