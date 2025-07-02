package internal

import (
	"database/sql"
	"reflect"
)

func (u *utilsPackage) GetScriptMigrate(db *sql.DB, dbName string, dialect Dialect, typ ...reflect.Type) ([]string, error) {
	dbSchema, err := utils.extractSchema(db, dbName, dialect)
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, t := range typ {
		sqlCmds, err := dialect.GenerateCreateTableSql(dbName, t)
		if err != nil {
			return nil, err
		}
		if sqlCmds != "" {
			ret = append(ret, sqlCmds)
		} else { //<--- table is existing in Database, just add columns
			sqlAddCols, err := dialect.GenerateAlterTableSql(dbName, t)
			if err != nil {
				return nil, err
			}
			ret = append(ret, sqlAddCols...)
		}
		sqlUniqueConstraints := dialect.GenerateUniqueConstraintsSql(t)

		if err != nil {
			return nil, err
		}

		for constraintName, sql := range sqlUniqueConstraints {
			if !dbSchema.unique[constraintName] {
				ret = append(ret, sql)
			}
		}
		sqlIndexConstraints := dialect.GenerateIndexConstraintsSql(t)
		if err != nil {
			return nil, err
		}
		for constraintName, sql := range sqlIndexConstraints {
			if !dbSchema.index[constraintName] {
				ret = append(ret, sql)
			}

		}

	}
	return ret, nil
}
