package dbx

import (
	"context"
	"database/sql"

	mssql "github.com/microsoft/go-mssqldb"
)

func (p *mssqlErrorParser) parseError8152(ctx context.Context, db *sql.DB, err mssql.Error) *DBXError {
	return &DBXError{
		Code:           DBXErrorCodeInvalidSize,
		Message:        err.Message,
		TableName:      "",
		ConstraintName: "",
		Fields:         nil,
		Values:         nil,
		MaxSize:        0,
	}
}
