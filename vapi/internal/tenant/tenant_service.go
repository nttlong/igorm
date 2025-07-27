package tenant

import (
	"vdb"
)

type TenantService struct {
	db      *vdb.TenantDB
	manager string
}

func (svc *TenantService) New(db *vdb.TenantDB, manager string) error {
	svc.db = db
	svc.manager = manager
	vdb.SetManagerDb(db.GetDriverName(), manager)
	return nil
}
func NewTenantService(db *vdb.TenantDB, manager string) *TenantService {
	svc := &TenantService{
		db:      db,
		manager: manager,
	}
	vdb.SetManagerDb(db.GetDriverName(), manager)
	return svc

}
func (svc *TenantService) Tenant(tenantName string) (*vdb.TenantDB, error) {
	tenant, err := svc.db.CreateDB(tenantName)

	if err != nil {
		return nil, err
	}
	return tenant, nil
}
