package vdb

import (
	"regexp"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func (d *mySqlDialect) ParseError1062(err *mysql.MySQLError) *DialectError {

	if strings.Contains(err.Message, "Duplicate entry ") {
		re := regexp.MustCompile(`for key '([^']+)'`)
		match := re.FindStringSubmatch(err.Message)
		if len(match) > 1 {
			//"for key 'users.UQ_users__username'"
			constraintName := match[1]
			if strings.Contains(constraintName, ".") {
				//"users.UQ_users__username"
				constraintName = strings.Split(constraintName, ".")[1]
			}
			ret := &DialectError{
				Err:            err,
				ErrorType:      DIALECT_DB_ERROR_TYPE_DUPLICATE,
				ErrorMessage:   "duplicate",
				ConstraintName: constraintName,
			}
			ret.Reload()
			return ret
		} else {
			ret := &DialectError{
				Err:          err,
				ErrorType:    DIALECT_DB_ERROR_TYPE_DUPLICATE,
				ErrorMessage: "duplicate",
			}
			ret.Reload()
			return ret
		}

	}
	return nil

}
