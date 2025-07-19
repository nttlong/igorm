package vdb

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	_ "github.com/microsoft/go-mssqldb"
)

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
func (d *mssqlDialect) ToParam(index int) string {
	return fmt.Sprintf("@p%d", index)
}
func (d *mssqlDialect) SqlFunction(delegator *DialectDelegateFunction) (string, error) {
	//delegator.Approved = true
	delegator.FuncName = strings.ToUpper(delegator.FuncName)
	return "", nil
}
