package vdb

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

type postgresDialect struct {
	cacheMakeSqlInsert sync.Map
}

func (d *postgresDialect) LikeValue(val string) string {

	return replaceStarWithCache("postgres", val, '*', '%')
}
func (d *postgresDialect) Name() string {
	return "postgres"
}
func (d *postgresDialect) Quote(name ...string) string {
	return "\"" + strings.Join(name, "\".\"") + "\""
}
func (d *postgresDialect) GetTableAndColumnsDictionary(db *sql.DB) (map[string]string, error) {
	panic(fmt.Errorf("not implemented, see file eorm/dialect.mssql.go"))
}
func (d *postgresDialect) ToText(value string) string {
	return fmt.Sprintf("'%s'::citext", value)
}
func (d *postgresDialect) ToParam(index int) string {
	return fmt.Sprintf("$%d", index)
}
func (d *postgresDialect) SqlFunction(delegator *DialectDelegateFunction) (string, error) {
	switch delegator.FuncName {
	case "LEN":
		delegator.FuncName = "LENGTH"
		delegator.HandledByDialect = true
		return "LENGTH" + "(" + strings.Join(delegator.Args, ", ") + ")", nil

	default:

		return "", nil
	}
}

func (d *postgresDialect) ParseError(err error) error {
	panic(fmt.Errorf("not implemented, see file eorm/dialect.msPostgressql.go"))
}
