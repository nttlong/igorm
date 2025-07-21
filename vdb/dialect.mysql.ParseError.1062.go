package vdb

import (
	"regexp"

	"github.com/go-sql-driver/mysql"
)

func (d *mySqlDialect) ParseError1452(err *mysql.MySQLError) *DialectError {

	re := regexp.MustCompile("CONSTRAINT `([^`]+)`")
	match := re.FindStringSubmatch(err.Message)
	if len(match) > 1 {

		ret := &DialectError{
			ErrorType:      DIALECT_DB_ERROR_TYPE_REFERENCES,
			ConstraintName: match[1],
		}
		ret.Reload()
		return ret
	} else {
		return nil
	}

	return nil

}
