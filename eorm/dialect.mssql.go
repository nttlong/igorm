package eorm

import (
	"database/sql"
	"fmt"
	"strings"
)

type mssqlDialect struct {
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
