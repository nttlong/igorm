package dbx

import (
	"context"
	"database/sql"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func (p *mysqlErrorParser) ParseError1062(ctx context.Context, db *sql.DB, err mysql.MySQLError) error {

	//"Duplicate entry 'tenant3' for key 'tenants.Name_uk'"
	items := strings.Split(err.Message, "'")
	if len(items) != 5 {
		return &DBXError{
			Code:    DBXErrorCodeDuplicate,
			Message: err.Message,
		}
	}
	strConstraint := items[3]
	if strings.Contains(strConstraint, ".") {
		tableName := strings.Split(strConstraint, ".")[0]
		constraintName := strings.Split(strConstraint, ".")[1]
		if strings.Contains(constraintName, "_") {
			fieldName := strings.Split(constraintName, "_")[0]
			return &DBXError{
				Code:           DBXErrorCodeDuplicate,
				TableName:      tableName,
				Fields:         []string{fieldName},
				ConstraintName: constraintName,
				Values:         []string{items[1]},
			}
		}
	}
	return &DBXError{
		Code:    DBXErrorCodeDuplicate,
		Message: err.Message,
	}
}
