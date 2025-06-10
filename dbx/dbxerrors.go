package dbx

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	mssql "github.com/microsoft/go-mssqldb"
)

type DBXErrorCode int

const (
	DBXErrorCodeUnknown DBXErrorCode = iota
	DBXErrorCodeDuplicate
	DBXErrorCodeInvalidSize
	DBXErrorCodeReferenceConstraint
)

func (e DBXErrorCode) String() string {
	switch e {
	case DBXErrorCodeUnknown:
		return "unknown error"
	case DBXErrorCodeDuplicate:
		return "duplicate error"
	case DBXErrorCodeInvalidSize:
		return "invalid size error"
	case DBXErrorCodeReferenceConstraint:
		return "reference constraint error"
	default:
		return "unknown error"
	}
}

type DBXError struct {
	// Error code the value is one of DBXErrorCode
	Code DBXErrorCode `json:"code"`
	// Error message
	Message string `json:"message"`
	// table name
	TableName string `json:"tableName"`
	//constraint name
	ConstraintName string `json:"constraintName"`
	// list of column names caused the error
	Fields []string `json:"fields"`
	// values of columns caused the error
	Values []string `json:"values"`
}
type DBXMigrationError struct {
	Message   string `json:"message"`
	Err       error  `json:"error"`
	DBName    string `json:"dbName"`
	TableName string `json:"tableName"`
	Code      string `json:"code"`
	Sql       string `json:"sql"`
}

func (e *DBXError) Error() string {
	return e.Message
}
func (e DBXMigrationError) Error() string {
	return e.Message
}

// cache duplicate error to avoid multiple query to get column name
// key is table name + constraint name
// value is list of column names ex: "cold1,col2,col3"
var errorMssqlErrorDuplicateCache = sync.Map{}

func parseErrorByMssqlErrorDuplicate(ctx context.Context, db *sql.DB, err mssql.Error) *DBXError {
	ret := &DBXError{Code: DBXErrorCodeDuplicate, Message: "duplicate error"}
	//"Violation of UNIQUE KEY constraint 'User_Username_uk'. Cannot insert duplicate key in object 'dbo.User'. The duplicate key value is (testuser)."
	errMsg := err.Message
	//get contraint name
	constraintName := strings.Split(errMsg, "constraint '")[1]
	constraintName = strings.Split(constraintName, "'")[0]
	tableName := strings.Split(errMsg, "in object '")[1]
	tableName = strings.Split(tableName, "'")[0]
	tableName = strings.Split(tableName, ".")[1]
	value := strings.Split(errMsg, "The duplicate key value is (")[1]
	value = strings.Split(value, ")")[0]
	ret.Message = "duplicate error"
	ret.TableName = tableName
	ret.ConstraintName = constraintName
	ret.Values = strings.Split(value, ",")
	cacheKey := tableName + constraintName
	//check cache
	if strFields, ok := errorMssqlErrorDuplicateCache.Load(cacheKey); ok {
		if strFields != nil {
			ret.Fields = strings.Split(strFields.(string), ",")
			return ret
		}
	}
	ret.Fields = dbxEntityCache.get_uk(constraintName)

	return ret
}

var parseErrorByMssqlErrorCache = sync.Map{} //cache error to avoid multiple query to get column name

func parseErrorByMssqlError(ctx context.Context, db *sql.DB, err error) *DBXError {
	if mssqlErr, ok := err.(mssql.Error); ok {
		key := fmt.Sprintf("%s:%s", mssqlErr.Number, mssqlErr.Message)
		if v, ok := parseErrorByMssqlErrorCache.Load(key); ok {
			return v.(*DBXError)
		}
		ret := parseErrorByMssqlErrorNoCache(ctx, db, err)
		parseErrorByMssqlErrorCache.Store(key, ret)
		return ret
	} else {
		return parseErrorByMssqlErrorNoCache(ctx, db, err)
	}

}
func parseErrorByMssqlErrorNoCache(ctx context.Context, db *sql.DB, err error) *DBXError {
	if err == nil {
		return nil
	}
	if mssqlErr, ok := err.(mssql.Error); ok {
		switch mssqlErr.Number {
		case 2601:
			return &DBXError{Code: DBXErrorCodeDuplicate, Message: "duplicate error"}
		case 2627:

			return parseErrorByMssqlErrorDuplicate(ctx, db, mssqlErr)
		case 547:
			return &DBXError{Code: DBXErrorCodeReferenceConstraint, Message: "reference constraint error"}
		case 50000:
			return &DBXError{Code: DBXErrorCodeInvalidSize, Message: "invalid size error"}
		default:
			return &DBXError{Code: DBXErrorCodeUnknown, Message: "unknown error"}
		}
	}
	return nil
}
