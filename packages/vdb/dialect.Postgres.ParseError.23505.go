package vdb

import (
	"vdb/migrate"

	"github.com/lib/pq"
)

func (d *postgresDialect) ParseError23505(dbSchame *migrate.DbSchema, err *pq.Error) error {
	ukContraint := err.Constraint
	if colsInfo, ok := dbSchame.UniqueKeys[ukContraint]; ok {
		dbCols := []string{}
		for _, col := range colsInfo.Columns {
			dbCols = append(dbCols, col.Name)
		}
		return &DialectError{
			Err:            err,
			ErrorType:      DIALECT_DB_ERROR_TYPE_DUPLICATE,
			DbCols:         dbCols,
			Table:          colsInfo.TableName,
			Fields:         dbCols,
			ErrorMessage:   "duplicate",
			ConstraintName: ukContraint,
		}
	}

	return err
}
