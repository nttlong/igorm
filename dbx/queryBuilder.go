package dbx

import (
	"context"
	"reflect"
)

type QrBuilder[T any] struct {
	dbx      *DBXTenant
	selector string
	where    string
	from     string
	args     []interface{}
	ctx      context.Context
}

func Query[T any](dbx *DBXTenant, ctx context.Context) *QrBuilder[T] {
	entityType, err := newEntityType(reflect.TypeFor[T]())
	if err != nil {
		panic(err)
	}
	return &QrBuilder[T]{
		dbx:      dbx,
		selector: "*",
		from:     entityType.TableName,
		ctx:      ctx,
	}

}

func (q QrBuilder[T]) First() (*T, error) {
	var zero T
	et := reflect.TypeOf(zero)
	entityType, err := newEntityType(et)
	if err != nil {
		return nil, err
	}
	sqlSelect := ""
	if q.where == "" {
		sqlSelect = "SELECT * FROM " + entityType.TableName + " LIMIT 1"
	} else {
		sqlSelect = "SELECT * FROM " + entityType.TableName + " WHERE " + q.where + " LIMIT 1"
	}
	rows, err := q.dbx.Query(sqlSelect, q.args...)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, nil
	}
	ret, err := fetchAllRows(rows.Rows, et)
	if err != nil {
		return nil, err
	}

	if len(ret.([]T)) == 0 {
		return nil, nil
	}
	retItem := ret.([]T)[0]
	return &retItem, nil
}

func Select[T any](dbx *DBXTenant, sql string, args ...interface{}) ([]T, error) {
	//fx := Select
	var zero T
	et := reflect.TypeOf(zero)

	rows, err := dbx.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	ret, err := fetchAllRows(rows.Rows, et)
	if err != nil {
		return nil, err
	}
	if len(ret.([]T)) == 0 {
		return nil, nil
	}
	retItem := ret.([]T)
	return retItem, nil

}
