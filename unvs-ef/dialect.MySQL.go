package unvsef

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type MySQLDialect struct {
	DB     *sql.DB
	schema map[string]TableSchema
}

func (d *MySQLDialect) Func(name string, args ...Expr) Expr {
	return rawFunc{name, args}
}

func (d *MySQLDialect) QuoteIdent(table, column string) string {
	return fmt.Sprintf("`%s`.`%s`", table, column)
}

func (d *MySQLDialect) TableExists(name string) bool {
	tbl, ok := d.schema[strings.ToLower(name)]
	return ok && tbl.Name != ""
}

func (d *MySQLDialect) ColumnExists(table, column string) bool {
	tbl, ok := d.schema[strings.ToLower(table)]
	if !ok {
		return false
	}
	_, colExists := tbl.Columns[strings.ToLower(column)]
	return colExists
}

func (d *MySQLDialect) SchemaMap() map[string]TableSchema {
	return d.schema
}

func (d *MySQLDialect) RefreshSchemaCache() error {
	d.schema = map[string]TableSchema{}

	// Load columns
	rows, err := d.DB.Query(`
		SELECT TABLE_NAME, COLUMN_NAME, DATA_TYPE, IS_NULLABLE
		FROM information_schema.columns
		WHERE table_schema = DATABASE()
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, columnName, dataType, isNullable string
		if err := rows.Scan(&tableName, &columnName, &dataType, &isNullable); err != nil {
			return err
		}

		key := strings.ToLower(tableName)
		tbl := d.schema[key]
		if tbl.Name == "" {
			tbl.Name = tableName
			tbl.Columns = map[string]ColumnSchema{}
		}
		tbl.Columns[strings.ToLower(columnName)] = ColumnSchema{
			Name:     columnName,
			Type:     dataType,
			Nullable: isNullable == "YES",
		}
		d.schema[key] = tbl
	}

	// Load unique constraints
	uniqueRows, err := d.DB.Query(`
		SELECT TABLE_NAME, CONSTRAINT_NAME
		FROM information_schema.table_constraints
		WHERE CONSTRAINT_TYPE = 'UNIQUE' AND table_schema = DATABASE()
	`)
	if err != nil {
		return err
	}
	defer uniqueRows.Close()

	for uniqueRows.Next() {
		var tableName, constraintName string
		if err := uniqueRows.Scan(&tableName, &constraintName); err != nil {
			return err
		}
		key := strings.ToLower(tableName)
		tbl := d.schema[key]
		tbl.UniqueConstraints = append(tbl.UniqueConstraints, constraintName)
		d.schema[key] = tbl
	}

	// Load index constraints (excluding PRIMARY and UNIQUE)
	indexRows, err := d.DB.Query(`
		SELECT TABLE_NAME, INDEX_NAME
		FROM information_schema.statistics
		WHERE NON_UNIQUE = 1 AND table_schema = DATABASE()
	`)
	if err != nil {
		return err
	}
	defer indexRows.Close()

	for indexRows.Next() {
		var tableName, indexName string
		if err := indexRows.Scan(&tableName, &indexName); err != nil {
			return err
		}
		key := strings.ToLower(tableName)
		tbl := d.schema[key]
		tbl.IndexConstraints = append(tbl.IndexConstraints, indexName)
		d.schema[key] = tbl
	}

	return nil
}
func (d *MySQLDialect) GenerateCreateTableSQL(typ reflect.Type) (string, error) {
	// TODO: implement this
	panic("unimplemented")
}
