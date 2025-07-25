package vdb

import (
	"reflect"
	"strings"
	"sync"
	"vdb/migrate"
)

type makeMySqlSqlInsertInit struct {
	once sync.Once
	val  string
}

func (d *mySqlDialect) MakeSqlInsert(tableName string, columns []migrate.ColumnDef, value interface{}) (string, []interface{}) {
	key := d.Name() + "://" + tableName
	actual, _ := d.cacheMakeSqlInsert.LoadOrStore(key, &makeMySqlSqlInsertInit{})
	init := actual.(*makeMySqlSqlInsertInit)
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
func (d *mySqlDialect) makeSqlInsert(tableName string, columns []migrate.ColumnDef) string {
	sql := "INSERT INTO " + d.Quote(tableName) + " ("
	strFields := []string{}
	strValues := []string{}

	for _, col := range columns {
		if col.IsAuto {
			// MySQL: bỏ qua trường tự tăng, nhưng không dùng OUTPUT như SQL Server
			continue
		}
		strFields = append(strFields, d.Quote(col.Name))
		strValues = append(strValues, "?")
	}

	sql += strings.Join(strFields, ", ") + ") VALUES (" + strings.Join(strValues, ", ") + ")"
	return sql
}
