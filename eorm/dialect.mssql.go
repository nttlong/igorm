package eorm

import (
	"database/sql"
	"eorm/migrate"
	"fmt"
	"reflect"
	"strings"
	"sync"

	_ "github.com/microsoft/go-mssqldb"

	mssql "github.com/microsoft/go-mssqldb"
)

type makeSqlInsertInit struct {
	once sync.Once
	val  string
}
type mssqlDialect struct {
	cacheMakeSqlInsert sync.Map
}

func (d *mssqlDialect) Quote(name ...string) string {
	return "[" + strings.Join(name, "].[") + "]"
}
func (d *mssqlDialect) Name() string {
	return "mssql"
}
func (d *mssqlDialect) GetTableAndColumnsDictionary(db *sql.DB) (map[string]string, error) {
	panic(fmt.Errorf("not implemented, see file eorm/dialect.mssql.go"))
}
func (d *mssqlDialect) ToText(value string) string {
	return fmt.Sprintf("N'%s'", value)
}
func (d *mssqlDialect) SqlFunction(delegator *DialectDelegateFunction) (string, error) {
	//delegator.Approved = true
	delegator.FuncName = strings.ToUpper(delegator.FuncName)
	return "", nil
}
func (d *mssqlDialect) MakeSqlInsert(tableName string, columns []migrate.ColumnDef, value interface{}) (string, []interface{}) {
	key := d.Name() + "://" + tableName
	actual, _ := d.cacheMakeSqlInsert.LoadOrStore(key, &makeSqlInsertInit{})
	init := actual.(*makeSqlInsertInit)
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
func (d *mssqlDialect) ParseError(err error) DialectError {
	//go-mssqldb.Error
	if mssqlErr, ok := err.(mssql.Error); ok {
		return d.Error2627(mssqlErr)
	}

	panic(fmt.Errorf("not supported error type: %T, err: %v in file eorm/dialect.mssql.go", err, err))
}
