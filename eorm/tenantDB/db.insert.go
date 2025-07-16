package tenantDB

func (db *TenantDB) Insert(data interface{}) error {
	return OnDbInsertFunc(db, data)
}

type OnDbInsertFuncType func(db *TenantDB, data interface{}) error

var OnDbInsertFunc OnDbInsertFuncType
