package eorm

import (
	"strings"

	"github.com/go-sql-driver/mysql"
)

func (d *mySqlDialect) ParseError1048(err *mysql.MySQLError) *DialectError {
	col := ""
	if strings.Contains(err.Message, "Column '") {
		col = strings.Split(err.Message, "'")[1]
		col = strings.Split(col, "'")[0]
	}
	ret := &DialectError{
		Err:       err,
		ErrorType: DIALECT_DB_ERROR_TYPE_REQUIRED,
		DbCols:    []string{col},

		ErrorMessage: "require",
	}
	ret.Reload()
	return ret
}
