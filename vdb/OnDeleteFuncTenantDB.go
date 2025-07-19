package vdb

import (
	"reflect"
	"sync"
	"vdb/tenantDB"
)

func onDeleteFuncTenantDBNoCache(db *tenantDB.TenantDB, typ reflect.Type, filter string, args ...interface{}) (string, error) {
	dialect := dialectFactory.Create(db.GetDriverName())
	repoType := inserterObj.getEntityInfo(typ)
	tableName := dialect.Quote(repoType.tableName)
	compiler, err := NewExprCompiler(db)
	if err != nil {
		return "", err
	}

	compiler.context.purpose = build_purpose_where
	compiler.context.tables = []string{repoType.tableName}
	compiler.context.alias = map[string]string{repoType.tableName: repoType.tableName}
	err = compiler.buildWhere(filter)
	if err != nil {
		return "", err
	}

	target := compiler.content
	sql := "delete from " + tableName + " where " + target
	_, err = db.Exec(sql, args...)
	if err != nil {
		return "", err
	}
	return sql, nil

}

type itemOnDeleteFuncTenantDB struct {
	val string
	err error
}
type initOnDeleteFuncTenantDB struct {
	once sync.Once
	val  itemOnDeleteFuncTenantDB
}

var onDeleteFuncTenantDBCache sync.Map

func onDeleteFuncTenantDB(db *tenantDB.TenantDB, data interface{}, filter string, args ...interface{}) (int64, error) {
	typ := reflect.TypeOf(data)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	key := db.GetDriverName() + ":" + typ.String() + "//" + filter
	actual, _ := onDeleteFuncTenantDBCache.LoadOrStore(key, &initOnDeleteFuncTenantDB{})
	init := actual.(*initOnDeleteFuncTenantDB)
	init.once.Do(func() {
		init.val.val, init.val.err = onDeleteFuncTenantDBNoCache(db, typ, filter, args...)
	})
	if init.val.err != nil {
		return 0, init.val.err
	}
	r, err := db.Exec(init.val.val, args...)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func init() {
	tenantDB.OnDeleteFuncTenantDB = onDeleteFuncTenantDB
}
