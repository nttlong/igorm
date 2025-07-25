package migrate

import (
	"database/sql"
	"fmt"
	"vdb/tenantDB"
)

/*
This function load all table information from the mysql database and return a map[string]map[string]ColumnInfo
return map[<table name>]<column name>: ColumnInfo (table name and column name are name of table and column in mysql)

columnInfo is struct

		type ColumnInfo struct {

		Name string //Db field name

		DbType string //field name

		Nullable bool

		Length int
	}

@db is a pointer to the TenantDB object tenantDB.TenantDB is sql.DB
*/
func (m *MigratorLoaderMysql) LoadAllTable(db *tenantDB.TenantDB) (map[string]map[string]ColumnInfo, error) {
	query := `
		SELECT TABLE_NAME, COLUMN_NAME, DATA_TYPE, IS_NULLABLE, CHARACTER_MAXIMUM_LENGTH
		FROM INFORMATION_SCHEMA.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()

	result := make(map[string]map[string]ColumnInfo)

	for rows.Next() {
		var tableName, columnName, dataType, isNullable string
		var charMaxLength sql.NullInt64

		err := rows.Scan(&tableName, &columnName, &dataType, &isNullable, &charMaxLength)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if _, exists := result[tableName]; !exists {
			result[tableName] = make(map[string]ColumnInfo)
		}

		colInfo := ColumnInfo{
			Name:     columnName,
			DbType:   dataType,
			Nullable: isNullable == "YES",
			Length:   0,
		}
		if charMaxLength.Valid {
			colInfo.Length = int(charMaxLength.Int64)
		}

		result[tableName][columnName] = colInfo
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return result, nil
}
