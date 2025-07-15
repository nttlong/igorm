package migrate

import (
	"database/sql"
	"eorm/tenantDB"
	"fmt"
)

/*
this function will all primary key in MySql database and return a map of table name and column info
return map[<Primary key constraint name>]ColumnsInfo, error
struct ColumnsInfo  below:

	type ColumnsInfo struct {
		TableName string
		Columns   []ColumnInfo
	}
	type ColumnInfo struct {

			Name string //Db field name

			DbType string //Db field type

			Nullable bool

			Length int
		}
		tenantDB.TenantDB is sql.DB
*/
func (m *MigratorLoaderMysql) LoadAllPrimaryKey(db *tenantDB.TenantDB) (map[string]ColumnsInfo, error) {
	query := `
		SELECT
			kcu.CONSTRAINT_NAME,
			kcu.TABLE_NAME,
			kcu.COLUMN_NAME,
			c.DATA_TYPE,
			c.IS_NULLABLE,
			c.CHARACTER_MAXIMUM_LENGTH
		FROM
			INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
			JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu
				ON tc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
				AND tc.TABLE_SCHEMA = kcu.TABLE_SCHEMA
				AND tc.TABLE_NAME = kcu.TABLE_NAME
			JOIN INFORMATION_SCHEMA.COLUMNS c
				ON c.TABLE_SCHEMA = kcu.TABLE_SCHEMA
				AND c.TABLE_NAME = kcu.TABLE_NAME
				AND c.COLUMN_NAME = kcu.COLUMN_NAME
		WHERE
			tc.CONSTRAINT_TYPE = 'PRIMARY KEY'
			AND tc.TABLE_SCHEMA = DATABASE()
		ORDER BY
			kcu.CONSTRAINT_NAME, kcu.ORDINAL_POSITION
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query primary keys: %w", err)
	}
	defer rows.Close()

	result := make(map[string]ColumnsInfo)

	for rows.Next() {
		var constraintName, tableName, columnName, dataType, isNullable string
		var charMaxLength sql.NullInt64

		if err := rows.Scan(&constraintName, &tableName, &columnName, &dataType, &isNullable, &charMaxLength); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		col := ColumnInfo{
			Name:     columnName,
			DbType:   dataType,
			Nullable: isNullable == "YES",
			Length:   0,
		}
		if charMaxLength.Valid {
			col.Length = int(charMaxLength.Int64)
		}
		fakeConstraintName := fmt.Sprintf("%s_%s", constraintName, tableName)
		if _, exists := result[fakeConstraintName]; !exists {
			result[fakeConstraintName] = ColumnsInfo{
				TableName: tableName,
				Columns:   []ColumnInfo{col},
			}
		} else {
			cols := result[fakeConstraintName].Columns
			cols = append(cols, col)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return result, nil
}
