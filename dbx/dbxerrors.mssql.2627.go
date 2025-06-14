package dbx

import (
	"context"
	"database/sql"
	"strings"

	mssql "github.com/microsoft/go-mssqldb"
)

func (p *mssqlErrorParser) parseError2627(ctx context.Context, db *sql.DB, err mssql.Error) *DBXError {
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
