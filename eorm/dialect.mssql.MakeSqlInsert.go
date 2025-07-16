package dbv

import (
	"dbv/migrate"
	"reflect"
	"strings"
	"sync"
)

type makeMssqlSqlInsertInit struct {
	once sync.Once
	val  string
}

func (d *mssqlDialect) MakeSqlInsert(tableName string, columns []migrate.ColumnDef, value interface{}) (string, []interface{}) {
	key := d.Name() + "://" + tableName
	actual, _ := d.cacheMakeSqlInsert.LoadOrStore(key, &makeMssqlSqlInsertInit{})
	init := actual.(*makeMssqlSqlInsertInit)
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
func (d *mssqlDialect) makeSqlInsert(tableName string, columns []migrate.ColumnDef) string {
	sql := "INSERT INTO " + d.Quote(tableName) + " ("
	strFields := []string{}
	strValues := []string{}
	insertedFieldName := ""
	for _, col := range columns {
		if col.IsAuto {
			insertedFieldName = col.Name
			continue
		}
		strFields = append(strFields, d.Quote(col.Name))
		strValues = append(strValues, "?")
	}
	if insertedFieldName != "" {
		sql += strings.Join(strFields, ", ") + ") OUTPUT INSERTED." + d.Quote(insertedFieldName) + " VALUES (" + strings.Join(strValues, ", ") + ")"
	} else {
		sql += strings.Join(strFields, ", ") + ") VALUES (" + strings.Join(strValues, ", ") + ")"
	}

	return sql
}
