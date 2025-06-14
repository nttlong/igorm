package dbx

import (
	"context"
	"database/sql"
	"sync"

	mssql "github.com/microsoft/go-mssqldb"
)

// cache duplicate error to avoid multiple query to get column name
// key is table name + constraint name
// value is list of column names ex: "cold1,col2,col3"
var errorMssqlErrorDuplicateCache = sync.Map{}

var parseErrorByMssqlErrorCache = sync.Map{} //cache error to avoid multiple query to get column name

type mssqlErrorParser struct {
}

func (p *mssqlErrorParser) ParseError(ctx context.Context, db *sql.DB, err error) *DBXError {
	if err == nil {
		return nil
	}
	if mssqlErr, ok := err.(mssql.Error); ok {
		switch mssqlErr.Number {
		case 8152:
			return p.parseError8152(ctx, db, mssqlErr)
		case 2601:
			return &DBXError{Code: DBXErrorCodeDuplicate, Message: "duplicate error"}
		case 2627:

			return p.parseError2627(ctx, db, mssqlErr)
		case 547:
			return &DBXError{Code: DBXErrorCodeReferenceConstraint, Message: "reference constraint error"}
		case 50000:
			return &DBXError{Code: DBXErrorCodeInvalidSize, Message: "invalid size error"}
		default:
			return &DBXError{Code: DBXErrorCodeUnknown, Message: "unknown error"}
		}
	}
	return &DBXError{Code: DBXErrorCodeUnknown, Message: "unknown error"}
}

var MssqlErrorParser *mssqlErrorParser

func init() {
	MssqlErrorParser = &mssqlErrorParser{}
}
