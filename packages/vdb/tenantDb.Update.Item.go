package vdb

import (
	"context"
	"reflect"
	"strings"
	"sync"
)

type initMakeUpdateSqlFromType struct {
	once sync.Once
	val  *initMakeUpdateSqlFromTypeItem
	err  error
}
type initMakeUpdateSqlFromTypeItem struct {
	sql           string
	fieldIndex    [][]int
	keyFieldIndex [][]int
	err           error
}

var makeUpdateSqlFromTypeWithCacheData = sync.Map{}

func makeUpdateSqlFromTypeWithCache(db *TenantDB, typ reflect.Type) (*initMakeUpdateSqlFromTypeItem, error) {
	key := db.GetDriverName() + ":" + typ.String()
	actual, _ := makeUpdateSqlFromTypeWithCacheData.LoadOrStore(key, &initMakeUpdateSqlFromType{})
	init := actual.(*initMakeUpdateSqlFromType)
	init.once.Do(func() {
		init.val, init.err = makeUpdateSqlFromType(db, typ)
	})
	return init.val, init.err

}
func makeUpdateSqlFromType(db *TenantDB, typ reflect.Type) (*initMakeUpdateSqlFromTypeItem, error) {
	ret := initMakeUpdateSqlFromTypeItem{
		sql:           "",
		fieldIndex:    nil,
		keyFieldIndex: nil,
	}

	model := db.getModelFromCache(typ)
	if model.err != nil {
		return nil, model.err
	}
	dialect := model.dialect

	sql := "UPDATE " + dialect.Quote(model.tableName) + " SET "
	conditional := []string{}

	strPlaceHoldesr := []string{}

	for i, col := range *model.cols {
		if col.PKName != "" || col.IsAuto {
			if col.PKName != "" {
				conditional = append(conditional, dialect.Quote(col.Name)+" = "+dialect.ToParam(i+1))
				ret.keyFieldIndex = append(ret.keyFieldIndex, col.IndexOfField)

			}
			continue

		}
		strPlaceHoldesr = append(strPlaceHoldesr, col.Name+" = "+dialect.ToParam(i+1))
		ret.fieldIndex = append(ret.fieldIndex, col.IndexOfField)

	}
	sql += strings.Join(strPlaceHoldesr, ",")
	if len(conditional) > 0 {
		sql += " WHERE " + strings.Join(conditional, " AND ")
	}
	ret.sql = sql
	return &ret, nil
}
func (db *TenantDB) UpdateWithContext(context context.Context, item interface{}) UpdateResult {
	typ := reflect.TypeOf(item)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()

	}
	info, err := makeUpdateSqlFromTypeWithCache(db, typ)
	if err != nil {
		return UpdateResult{RowsAffected: 0, Sql: "", Error: err}
	}
	val := reflect.ValueOf(item).Elem()
	args := make([]interface{}, 0)
	for _, index := range info.fieldIndex {
		args = append(args, val.FieldByIndex(index).Interface())
	}
	for _, index := range info.keyFieldIndex {
		args = append(args, val.FieldByIndex(index).Interface())
	}
	r, err := db.ExecContext(context, info.sql, args...)
	if err != nil {
		return UpdateResult{RowsAffected: 0, Sql: info.sql, Error: err}
	}
	n, err := r.RowsAffected()
	if err != nil {
		return UpdateResult{RowsAffected: 0, Sql: info.sql, Error: err}
	}
	return UpdateResult{RowsAffected: n, Sql: info.sql, Error: nil}
	// compiler, err := NewExprCompiler(db.TenantDB)
	// if err != nil {
	// 	return DeleteResult{Error: err}
	// }

	// model := db.getModelFromCache(typ)
	// dialect := model.dialect
	// val := reflect.ValueOf(item).Elem()
	// sql := "UPDATE " + dialect.Quote(model.tableName) + " SET "
	// args := make([]interface{}, 0)
	// strPlaceHoldesr := []string{}
	// where := ""

	// for i, col := range *model.cols {
	// 	if col.PKName != "" || col.IsAuto {
	// 		if col.PKName != "" {
	// 			where = dialect.Quote(col.Name) + " = " + dialect.ToParam(i+1)
	// 			args = append(args, val.FieldByIndex(col.IndexOfField).Interface())
	// 		}
	// 		continue

	// 	}
	// 	strPlaceHoldesr = append(strPlaceHoldesr, col.Name+" = "+dialect.ToParam(i+1))
	// 	args = append(args, val.FieldByIndex(col.IndexOfField).Interface())
	// }
	// sql += strings.Join(strPlaceHoldesr, ",")
	// if where != "" {
	// 	sql += " WHERE " + where
	// }
	// r, err := db.Exec(sql, args...)
	// if err != nil {
	// 	return UpdateResult{RowsAffected: 0, Sql: sql, Error: err}
	// }
	// n, err := r.RowsAffected()
	// if err != nil {
	// 	return UpdateResult{RowsAffected: 0, Sql: sql, Error: err}
	// }
	// return UpdateResult{RowsAffected: n, Sql: sql, Error: nil}

}
func (db *TenantDB) Update(item interface{}) UpdateResult {
	return db.UpdateWithContext(context.Background(), item)

}
