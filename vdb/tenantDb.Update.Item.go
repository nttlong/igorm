package vdb

import (
	"context"
	"reflect"
	"strings"
)

func (db *TenantDB) UpdateWithContext(context context.Context, item interface{}) UpdateResult {
	typ := reflect.TypeOf(item)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()

	}
	// compiler, err := NewExprCompiler(db.TenantDB)
	// if err != nil {
	// 	return DeleteResult{Error: err}
	// }

	model := db.getModelFromCache(typ)
	dialect := model.dialect
	val := reflect.ValueOf(item).Elem()
	sql := "UPDATE " + dialect.Quote(model.tableName) + " SET "
	args := make([]interface{}, 0)
	strPlaceHoldesr := []string{}
	where := ""

	for i, col := range *model.cols {
		if col.PKName != "" || col.IsAuto {
			if col.PKName != "" {
				where = dialect.Quote(col.Name) + " = " + dialect.ToParam(i+1)
				args = append(args, val.FieldByIndex(col.IndexOfField).Interface())
			}
			continue

		}
		strPlaceHoldesr = append(strPlaceHoldesr, col.Name+" = "+dialect.ToParam(i+1))
		args = append(args, val.FieldByIndex(col.IndexOfField).Interface())
	}
	sql += strings.Join(strPlaceHoldesr, ",")
	if where != "" {
		sql += " WHERE " + where
	}
	r, err := db.Exec(sql, args...)
	if err != nil {
		return UpdateResult{RowsAffected: 0, Sql: sql, Error: err}
	}
	n, err := r.RowsAffected()
	if err != nil {
		return UpdateResult{RowsAffected: 0, Sql: sql, Error: err}
	}
	return UpdateResult{RowsAffected: n, Sql: sql, Error: nil}

}
func (db *TenantDB) Update(item interface{}) UpdateResult {
	return db.UpdateWithContext(context.Background(), item)

}
