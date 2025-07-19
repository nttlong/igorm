package vdb

import (
	"fmt"
	"reflect"
	"sync"
	"vdb/tenantDB"
)

func buildDeleteSql(db *TenantDB, typ reflect.Type, filter string) (string, error) {

	model := db.getModelFromCache(typ)
	source := model.dialect.Quote(model.tableName)
	if filter != "" {
		compiler, err := NewExprCompiler(db.TenantDB)
		if err != nil {
			return "", err
		}
		compiler.context.purpose = build_purpose_where
		compiler.context.tables = []string{model.tableName}
		compiler.context.alias = map[string]string{model.tableName: model.tableName}
		compiler.context.paramIndex = 1
		err = compiler.buildWhere(filter)
		if err != nil {
			return "", err
		}
		filter = compiler.content
	}
	sql := fmt.Sprintf("DELETE FROM %s", source)
	if filter != "" {
		sql += " WHERE " + filter
	}
	return sql, nil

}

type initBuildDeleteSql struct {
	once sync.Once
	val  string
	err  error
}

var buildDeleteSqlCache = sync.Map{}

func buildDeleteSqlWithCache(db *TenantDB, typ reflect.Type, filter string) (string, error) {
	key := db.GetDriverName() + ":" + typ.String() + ":" + filter
	actual, _ := buildDeleteSqlCache.LoadOrStore(key, &initBuildDeleteSql{})
	init := actual.(*initBuildDeleteSql)
	init.once.Do(func() {
		init.val, init.err = buildDeleteSql(db, typ, filter)
	})
	return init.val, init.err
}
func doDelete(db *TenantDB, entityData interface{}, filter string, args ...interface{}) (int64, error) {
	typ := reflect.TypeOf(entityData)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	sql, err := buildDeleteSqlWithCache(db, typ, filter)
	if err != nil {
		return 0, err
	}

	result, err := db.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	r, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return r, nil

}
func init() {
	tenantDB.OnDeleteEntity = func(db *tenantDB.TenantDB, entity interface{}, filter string, args ...interface{}) (int64, error) {
		return doDelete(&TenantDB{TenantDB: db}, entity, filter, args...)
	}

}
