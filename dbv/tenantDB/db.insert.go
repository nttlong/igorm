package tenantDB

func (db *TenantDB) Insert(data ...interface{}) error {
	for _, d := range data {
		err := OnDbInsertFunc(db, d)
		if err != nil {
			return err
		}
	}
	return nil
}
func (tx *TenantTx) Insert(data ...interface{}) error {
	for _, d := range data {
		err := OnTxDbInsertFunc(tx, d)
		if err != nil {
			return err
		}
	}
	return nil
}

type OnDbInsertFuncType func(db *TenantDB, data interface{}) error
type OnTxDbInsertFuncType func(tx *TenantTx, data interface{}) error

var OnDbInsertFunc OnDbInsertFuncType
var OnTxDbInsertFunc OnTxDbInsertFuncType
