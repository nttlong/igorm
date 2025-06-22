package dbx

import (
	"context"
	"database/sql"
	"strings"

	"github.com/lib/pq"
)

type postgresErrorParser struct {
}

func (p *postgresErrorParser) ParseError(ctx context.Context, db *sql.DB, err error) *DBXError {
	if pgErr, ok := err.(*pq.Error); ok {
		//"23505"
		//"duplicate key value violates unique constraint \"AppConfig_Name_uk\""
		//"Key (\"Name\")=(fx001) already exists."
		// fmt.Println(pgErr.Code)
		// fmt.Println(pgErr.Message)
		if pgErr.Code == "23505" {
			strFields := strings.Split(pgErr.Detail, "=")[0]
			strFields = strings.Split(strFields, "(")[1]
			strFields = strings.Split(strFields, ")")[0]
			strFields = strings.ReplaceAll(strFields, "\"", "")

			strValues := strings.Split(pgErr.Detail, "=")[1]
			strValues = strings.Split(strValues, "(")[1]
			strValues = strings.Split(strValues, ")")[0]
			return &DBXError{
				Code:           DBXErrorCodeDuplicate,
				Message:        "duplicate key error",
				TableName:      pgErr.Table,
				ConstraintName: pgErr.Constraint,
				Fields:         strings.Split(strFields, ","),
				Values:         strings.Split(strValues, ","),
			}

		}
		if pgErr.Code == "23514" {

			dbField := ""
			if strings.Contains(pgErr.Constraint, "_check_length") {
				items := strings.Split(pgErr.Constraint, "_")
				if len(items) > 2 {
					if items[0] == pgErr.Table {
						dbField = items[1]
						entity := Entities.GetEntityTypeByTableName(pgErr.Table)
						if entity != nil {
							entityField := entity.GetFieldByName(dbField)
							if entityField != nil {
								return &DBXError{
									Code:           DBXErrorCodeInvalidSize,
									Message:        "invalid size error",
									TableName:      pgErr.Table,
									ConstraintName: pgErr.Constraint,
									Fields:         []string{dbField},
									MaxSize:        entityField.MaxLen,
								}
							}
						}
					}
				}
			}
			return &DBXError{
				Code:           DBXErrorCodeInvalidSize,
				Message:        "invalid size error",
				TableName:      pgErr.Table,
				ConstraintName: pgErr.Constraint,
			}
		}

	}
	return &DBXError{
		Code:    DBXErrorCodeUnknown,
		Message: err.Error(),
	}
}

var PostgresErrorParser postgresErrorParser

func init() {
	PostgresErrorParser = postgresErrorParser{}
}
