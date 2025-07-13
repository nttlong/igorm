package eorm

import (
	"eorm/migrate"
	"strings"

	_ "github.com/microsoft/go-mssqldb"

	mssql "github.com/microsoft/go-mssqldb"
)

func (d *mssqlDialect) Error2627(err mssql.Error) DialectError {

	if strings.Contains(err.Message, "'") {
		constraint := strings.Split(err.Message, "'")[1]
		constraint = strings.Split(constraint, "'")[0]

		result := migrate.FindUKConstraint(constraint)
		if result != nil {
			cols := []string{}
			fields := []string{}
			for _, col := range result.Columns {
				cols = append(cols, col.Name)
				fields = append(fields, col.Field.Name)
			}
			ret := DialectError{
				Err:          err,
				ErrorType:    DIALECT_DB_ERROR_TYPE_DUPLICATE,
				ErrorMessage: err.Message,
				DbCols:       cols,
				Fields:       fields,
				Tables:       []string{result.TableName},
			}
			ret.Reload()
			return ret
		}

	}
	// errorMsg := err.Message
	panic("not implemented")
}
