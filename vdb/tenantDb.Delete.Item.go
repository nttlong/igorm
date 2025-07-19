package vdb

import (
	"reflect"
	"sync"
)

type initGetSqlDeleteItem struct {
	once sync.Once
	val  string
}

var getSqlDeleteItemCache sync.Map

func getSqlDeleteItem(db *TenantDB, typ reflect.Type) string {

	model := db.getModelFromCache(typ)

	sql := "DELETE FROM " + model.dialect.Quote(model.tableName) + " WHERE "
	for i, col := range model.keyCols {
		if i > 0 {
			sql += " AND "
		}
		sql += col.Name + " = " + model.dialect.ToParam(i+1)
	}
	return sql

}
func getSqlDeleteItemCached(db *TenantDB, typ reflect.Type) string {
	key := db.GetDriverName() + "://" + db.GetDBName() + "://" + typ.String()
	actual, _ := getSqlDeleteItemCache.LoadOrStore(key, &initGetSqlDeleteItem{})
	init := actual.(*initGetSqlDeleteItem)
	init.once.Do(func() {
		init.val = getSqlDeleteItem(db, typ)
	})
	return init.val
}
func (db *TenantDB) Delete(item interface{}) DeleteResult {
	typ := reflect.TypeOf(item)
	val := reflect.ValueOf(item)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	model := db.getModelFromCache(typ)
	argsKey := []interface{}{}
	for _, col := range model.keyCols {
		argsKey = append(argsKey, val.FieldByIndex(col.IndexOfField).Interface())
	}
	sql := getSqlDeleteItemCached(db, typ)

	result, err := db.Exec(sql, argsKey...)
	if err != nil {
		return DeleteResult{Error: err}
	}
	r, err := result.RowsAffected()
	if err != nil {
		return DeleteResult{Error: err}
	}
	return DeleteResult{RowsAffected: r}

}
