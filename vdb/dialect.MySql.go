package vdb

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

type mySqlDialect struct {
	cacheMakeSqlInsert sync.Map
}

func (d *mySqlDialect) Quote(name ...string) string {
	return "`" + strings.Join(name, "`.`") + "`"
}
func (d *mySqlDialect) Name() string {
	return "mysql"
}
func (d *mySqlDialect) GetTableAndColumnsDictionary(db *sql.DB) (map[string]string, error) {
	panic(fmt.Errorf("not implemented, see file eorm/dialect.mssql.go"))
}
func (d *mySqlDialect) ToText(value string) string {
	return fmt.Sprintf("'%s'", value)
}
func (d *mySqlDialect) ToParam(index int) string {
	return fmt.Sprintf(":%d", index)
}
func (d *mySqlDialect) SqlFunction(delegator *DialectDelegateFunction) (string, error) {
	switch delegator.FuncName {
	case "NOW":
		delegator.HandledByDialect = true
		return "NOW()", nil
	default:

		return "", nil
	}
}
