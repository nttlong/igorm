package vdb

import (
	"reflect"
	"strings"
	"sync"
	"vdb/migrate"
)

type makePostgresSqlInsertInit struct {
	once sync.Once
	val  string
}

func (d *postgresDialect) MakeSqlInsert(tableName string, columns []migrate.ColumnDef, value interface{}) (string, []interface{}) {
	key := d.Name() + "://" + tableName
	actual, _ := d.cacheMakeSqlInsert.LoadOrStore(key, &makePostgresSqlInsertInit{})
	init := actual.(*makePostgresSqlInsertInit)
	init.once.Do(func() {
		init.val = d.makeSqlInsert(tableName, columns)
	})
	dataVal := reflect.ValueOf(value)
	if dataVal.Kind() == reflect.Ptr {
		dataVal = dataVal.Elem()
	}
	args := []interface{}{}
	for _, col := range columns {
		if col.IsAuto {
			continue
		}
		fieldVal := dataVal.FieldByName(col.Field.Name)
		if fieldVal.IsValid() {
			args = append(args, fieldVal.Interface())
		} else {
			args = append(args, nil)
		}

	}

	return init.val, args

}
func (d *postgresDialect) makeSqlInsert(tableName string, columns []migrate.ColumnDef) string {

	sql := "INSERT INTO " + d.Quote(tableName) + " ("
	strFields := []string{}
	strValues := []string{}
	i := 1
	RETURNING_ID := ""
	for _, col := range columns {
		if col.IsAuto {
			RETURNING_ID = " RETURNING " + d.Quote(col.Name)
			continue
		}
		strFields = append(strFields, d.Quote(col.Name))
		strValues = append(strValues, d.ToParam(i))
		i++
	}

	sql += strings.Join(strFields, ", ") + ") VALUES (" + strings.Join(strValues, ", ") + ")"
	if RETURNING_ID != "" {
		sql += RETURNING_ID
	}
	return sql
}
