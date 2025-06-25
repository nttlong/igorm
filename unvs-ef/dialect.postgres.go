package unvsef

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type PostgresDialect struct {
	DB     *sql.DB
	schema map[string]TableSchema
}

func (d *PostgresDialect) Func(name string, args ...Expr) Expr {
	return rawFunc{name, args}
}

func (d *PostgresDialect) QuoteIdent(table, column string) string {
	return fmt.Sprintf(`"%s"."%s"`, table, column)
}

func (d *PostgresDialect) TableExists(name string) bool {
	tbl, ok := d.schema[strings.ToLower(name)]
	return ok && tbl.Name != ""
}

func (d *PostgresDialect) ColumnExists(table, column string) bool {
	tbl, ok := d.schema[strings.ToLower(table)]
	if !ok {
		return false
	}
	_, colExists := tbl.Columns[strings.ToLower(column)]
	return colExists
}

func (d *PostgresDialect) SchemaMap() map[string]TableSchema {
	return d.schema
}

func (d *PostgresDialect) RefreshSchemaCache() error {
	d.schema = map[string]TableSchema{}

	// Load columns
	rows, err := d.DB.Query(`
		SELECT table_name, column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'public'
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

	// Load UNIQUE constraints
	uniqueRows, err := d.DB.Query(`
		SELECT tc.table_name, tc.constraint_name
		FROM information_schema.table_constraints tc
		WHERE tc.constraint_type = 'UNIQUE' AND tc.table_schema = 'public'
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

	// Load index constraints (excluding primary and unique)
	indexRows, err := d.DB.Query(`
		SELECT t.relname AS table_name, i.relname AS index_name
		FROM pg_class t
		JOIN pg_index ix ON t.oid = ix.indrelid
		JOIN pg_class i ON i.oid = ix.indexrelid
		WHERE NOT ix.indisprimary AND NOT ix.indisunique
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
func (d *PostgresDialect) GenerateCreateTableSQL(typ reflect.Type) (string, error) {
	// TODO: implement this
	panic("unimplemented")
}
