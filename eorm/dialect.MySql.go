package eorm

import (
	"database/sql"
	"eorm/migrate"
	"fmt"
	"strings"
)

type mySqlDialect struct {
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
func (d *mySqlDialect) SqlFunction(delegator *DialectDelegateFunction) (string, error) {
	switch delegator.FuncName {
	case "NOW":
		delegator.HandledByDialect = true
		return "NOW()", nil
	default:

		return "", nil
	}
}
func (d *mySqlDialect) MakeSqlInsert(tableName string, columns []migrate.ColumnDef, data interface{}) (string, []interface{}) {
	panic(fmt.Errorf("not implemented, see file eorm/dialect.mssql.go"))
}
func (d *mySqlDialect) ParseError(err error) DialectError {
	panic(fmt.Errorf("not implemented, see file eorm/dialect.mssql.go"))
}
