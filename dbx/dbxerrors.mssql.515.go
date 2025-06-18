package dbx

import (
	"context"
	"database/sql"
	"strings"

	mssql "github.com/microsoft/go-mssqldb"
)

func (p *mssqlErrorParser) parseError515(ctx context.Context, db *sql.DB, err mssql.Error) *DBXError {

	if strings.Contains(err.Message, "Cannot insert the value NULL into column '") {

		fieldName := strings.Split(err.Message, "Cannot insert the value NULL into column '")[1]
		fieldName = strings.Split(fieldName, "'")[0]

		return &DBXError{
			Code:           DBXErrorCodeMissingRequiredField,
			Message:        "Required field is missing",
			TableName:      "",
			ConstraintName: "",
			Fields:         []string{fieldName},
			Values:         []string{""},
			MaxSize:        0,
		}
	}
	return &DBXError{
		Code:           DBXErrorCodeUnknown,
		Message:        err.Message,
		TableName:      "",
		ConstraintName: "",
	}

}
