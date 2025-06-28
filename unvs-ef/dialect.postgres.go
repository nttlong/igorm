package unvsef

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type PostgresDialect struct {
	baseDialect
}

/*
Implementation of Dialect
Example: input "aaa","bbb" -> "aaa"."bbb"
*/
func (d *PostgresDialect) QuoteIdent(args ...string) string {
	return `"` + strings.Join(args, `"."`) + `"`
}

/*
Must be called after RefreshSchemaCache
*/
func (d *PostgresDialect) TableExists(dbName, name string) bool {
	if _, ok := d.schema[dbName]; !ok {
		return false
	}
	tbl, ok := d.schema[dbName][strings.ToLower(name)]
	return ok && tbl.Name != ""
}

/*
Must be called after RefreshSchemaCache
*/
func (d *PostgresDialect) ColumnExists(dbName, table, column string) bool {
	if _, ok := d.schema[strings.ToLower(dbName)]; !ok {
		return false
	}
	tbl, ok := d.schema[dbName][strings.ToLower(table)]
	if !ok {
		return false
	}
	_, colExists := tbl.Columns[strings.ToLower(column)]
	return colExists
}

/*
Must be called after RefreshSchemaCache
*/
func (d *PostgresDialect) SchemaMap(dbName string) map[string]TableSchema {
	if _, ok := d.schema[dbName]; !ok {
		return nil
	}
	return d.schema[dbName]
}

/*
Retrieve all tables, columns, unique constraints, index constraints, and foreign keys in the PostgreSQL database.
*/
func (d *PostgresDialect) RefreshSchemaCache(db *sql.DB, dbName string) error {
	schema := map[string]TableSchema{}

	// Load columns
	rows, err := db.Query(`
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
		tbl := schema[key]
		if tbl.Name == "" {
			tbl.Name = tableName
			tbl.Columns = map[string]ColumnSchema{}
		}
		tbl.Columns[strings.ToLower(columnName)] = ColumnSchema{
			Name:     columnName,
			Type:     dataType,
			Nullable: isNullable == "YES",
		}
		schema[key] = tbl
	}

	// Load unique constraints
	uniqueRows, err := db.Query(`
		SELECT 
			t.relname AS table_name,
			c.conname AS constraint_name
		FROM pg_constraint c
		JOIN pg_class t ON c.conrelid = t.oid
		WHERE c.contype = 'u'
	`)
	if err != nil {
		return err
	}
	defer uniqueRows.Close()

	for uniqueRows.Next() {
		var tableName, constraint string
		if err := uniqueRows.Scan(&tableName, &constraint); err != nil {
			return err
		}
		key := strings.ToLower(tableName)
		tbl := schema[key]
		tbl.UniqueConstraints = append(tbl.UniqueConstraints, constraint)
		schema[key] = tbl
	}

	// Load index constraints
	indexRows, err := db.Query(`
		SELECT 
			t.relname AS table_name,
			i.relname AS index_name
		FROM pg_index idx
		JOIN pg_class i ON idx.indexrelid = i.oid
		JOIN pg_class t ON idx.indrelid = t.oid
		WHERE NOT idx.indisprimary AND NOT idx.indisunique
	`)
	if err != nil {
		return err
	}
	defer indexRows.Close()

	for indexRows.Next() {
		var tableName, index string
		if err := indexRows.Scan(&tableName, &index); err != nil {
			return err
		}
		key := strings.ToLower(tableName)
		tbl := schema[key]
		tbl.IndexConstraints = append(tbl.IndexConstraints, index)
		schema[key] = tbl
	}

	// Load foreign key constraints
	fkRows, err := db.Query(`
		SELECT 
			t.relname AS table_name,
			c.conname AS constraint_name
		FROM pg_constraint c
		JOIN pg_class t ON c.conrelid = t.oid
		WHERE c.contype = 'f'
	`)
	if err != nil {
		return err
	}
	defer fkRows.Close()

	for fkRows.Next() {
		var tableName, fk string
		if err := fkRows.Scan(&tableName, &fk); err != nil {
			return err
		}
		key := strings.ToLower(tableName)
		tbl := schema[key]
		tbl.ForeignKeyConstraints = append(tbl.ForeignKeyConstraints, fk)
		schema[key] = tbl
	}

	d.schema[dbName] = schema
	return nil
}

func (d *PostgresDialect) GetSchema(db *sql.DB, dbName string) (map[string]TableSchema, error) {
	if _, ok := d.schema[dbName]; !ok {
		if err := d.RefreshSchemaCache(db, dbName); err != nil {
			d.schemaError = err
			return nil, err
		}
	}
	return d.schema[dbName], nil
}

/*
Generate SQL CREATE TABLE statement for a struct, including primary key, columns with constraints like NOT NULL or DEFAULT,
and CHECK constraints for columns with length restrictions.
*/
func (d *PostgresDialect) GenerateCreateTableSql(dbName string, typ reflect.Type) (string, error) {
	meta := utils.GetMetaInfo(typ)
	for tableName, fields := range meta {
		if d.TableExists(dbName, tableName) {
			return "", nil
		}
		var colDefs []string
		var checkConstraints []string
		for colName, field := range fields {
			colType := field.DBType
			if _colType, ok := d.mapGoTypeToDb[utils.ResolveFieldKind(field.Field)]; ok {
				colType = _colType
			}
			if field.DBType != "" {
				colType = field.DBType
			} else if field.Length != nil {
				colType = "citext" // Use citext for fields with length
				// Generate CHECK constraint for length
				checkName := utils.ToSnakeCase(tableName) + "_" + utils.ToSnakeCase(colName)
				checkConstraints = append(checkConstraints,
					fmt.Sprintf("CONSTRAINT %s CHECK (LENGTH(%s) <= %d)",
						utils.Quote(`"`, checkName),
						utils.Quote(`"`, colName),
						*field.Length))
			}

			colDef := []string{utils.Quote(`"`, colName), colType}
			if field.AutoIncrement {
				colDef = append(colDef, "SERIAL")
			}
			if field.Nullable {
				colDef = append(colDef, "NULL")
			} else {
				colDef = append(colDef, "NOT NULL")
			}
			if field.Default != "" {
				if mapDefaultValue, ok := d.mapDefaultValue[field.Default]; ok {
					colDef = append(colDef, "DEFAULT "+mapDefaultValue)
				}
			}
			colDefs = append(colDefs, strings.Join(colDef, " "))
		}
		pkConstraint, err := d.GetPkConstraint(typ)
		if err != nil {
			return "", err
		}
		// Combine column definitions, primary key constraint, and check constraints
		allConstraints := append([]string{strings.Join(colDefs, ", ")}, checkConstraints...)
		if pkConstraint != "" {
			allConstraints = append(allConstraints, pkConstraint)
		}
		stmt := fmt.Sprintf("CREATE TABLE %s (%s)",
			utils.Quote(`"`, tableName),
			strings.Join(allConstraints, ", "))
		return stmt, nil
	}
	return "", nil
}

/*
Generate primary key constraint for a table.
Example: CONSTRAINT primary____user_roles___user_id__role_id PRIMARY KEY (user_id, role_id)
*/
func (d *PostgresDialect) GetPkConstraint(typ reflect.Type) (string, error) {
	mapConstraints := utils.GetPkFromMetaByType(typ)
	tableName := utils.TableNameFromStruct(typ)

	cols := []string{}
	for _, constraints := range mapConstraints {
		for fieldName := range constraints {
			cols = append(cols, utils.Quote(`"`, fieldName))
		}
		constraintName := "primary____" + tableName + "___" + strings.Join(cols, "__")
		ret := fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", utils.Quote(`"`, constraintName), strings.Join(cols, ", "))
		return ret, nil
	}
	return "", fmt.Errorf("no primary key found for table %s", utils.TableNameFromStruct(typ))
}

/*
Generate SQL ALTER TABLE ADD COLUMN statements for columns missing in the database.
*/
func (d *PostgresDialect) GenerateAlterTableSql(dbName string, typ reflect.Type) ([]string, error) {
	tableName := utils.TableNameFromStruct(typ)
	meta := utils.GetMetaInfo(typ)
	var metaField map[string]FieldTag
	if _metaField, ok := meta[tableName]; ok {
		metaField = _metaField
	} else {
		return nil, fmt.Errorf("table %s not found in meta", tableName)
	}

	var alters []string
	for fieldName, field := range metaField {
		if d.ColumnExists(dbName, tableName, fieldName) {
			continue
		}
		colType := field.DBType
		if _colType, ok := d.mapGoTypeToDb[utils.ResolveFieldKind(field.Field)]; ok {
			colType = _colType
		}
		if field.DBType != "" {
			colType = field.DBType
		} else if field.Length != nil {
			colType = "citext"
		}

		colDef := []string{utils.Quote(`"`, fieldName), colType}
		if field.PrimaryKey {
			continue // primary key cannot be added
		}
		if field.AutoIncrement {
			colDef = append(colDef, "SERIAL")
		}
		if field.Default != "" {
			if mapDefaultValue, ok := d.mapDefaultValue[field.Default]; ok {
				colDef = append(colDef, "DEFAULT "+mapDefaultValue)
			} else if mapDefaultValue, ok := d.mapDefaultValue[utils.ResolveFieldKind(field.Field)]; ok {
				colDef = append(colDef, "DEFAULT "+mapDefaultValue)
			}
		} else if !field.Nullable {
			if mapDefaultValue, ok := d.mapDefaultValue[utils.ResolveFieldKind(field.Field)]; ok {
				colDef = append(colDef, "DEFAULT "+mapDefaultValue)
			}
		}

		stmt := fmt.Sprintf("ALTER TABLE %s ADD %s", utils.Quote(`"`, tableName), strings.Join(colDef, " "))
		// Add CHECK constraint for fields with length
		if field.Length != nil {
			checkName := utils.ToSnakeCase(tableName) + "_" + utils.ToSnakeCase(fieldName)
			checkStmt := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s CHECK (LENGTH(%s) <= %d)",
				utils.Quote(`"`, tableName),
				utils.Quote(`"`, checkName),
				utils.Quote(`"`, fieldName),
				*field.Length)
			alters = append(alters, stmt, checkStmt)
		} else {
			alters = append(alters, stmt)
		}
	}
	return alters, nil
}

/*
Generate SQL statements for unique constraints.
Example: ALTER TABLE "users" ADD CONSTRAINT "users__email" UNIQUE ("email")
*/
func (d *PostgresDialect) GenerateUniqueConstraintsSql(typ reflect.Type) map[string]string {
	tableName := utils.TableNameFromStruct(typ)
	tableName = utils.Quote(`"`, tableName)
	mapUniqueConstraints := utils.GetUniqueConstraintsFromMetaByType(typ)
	ret := map[string]string{}
	for constraintName, constraints := range mapUniqueConstraints {
		cols := []string{}
		for fieldName := range constraints {
			cols = append(cols, utils.Quote(`"`, fieldName))
		}
		key := constraintName
		constraintName = utils.Quote(`"`, constraintName)
		ret[key] = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (%s)", tableName, constraintName, strings.Join(cols, ", "))
	}
	return ret
}

/*
Generate SQL statements for index constraints.
Example: CREATE INDEX "users__email_idx" ON "users" ("email")
*/
func (d *PostgresDialect) GenerateIndexConstraintsSql(typ reflect.Type) map[string]string {
	tableName := utils.TableNameFromStruct(typ)
	tableName = utils.Quote(`"`, tableName)
	mapIndexConstraints := utils.GetIndexConstraintsFromMetaByType(typ)
	ret := map[string]string{}
	for constraintName, constraints := range mapIndexConstraints {
		cols := []string{}
		for fieldName := range constraints {
			cols = append(cols, utils.Quote(`"`, fieldName))
		}
		key := constraintName
		constraintName = utils.Quote(`"`, constraintName)
		ret[key] = fmt.Sprintf("CREATE INDEX %s ON %s (%s)", constraintName, tableName, strings.Join(cols, ", "))
	}
	return ret
}

func NewPostgresDialect(db *sql.DB) Dialect {
	return &PostgresDialect{
		baseDialect: baseDialect{
			schema: map[string]map[string]TableSchema{},
			mapGoTypeToDb: map[string]string{
				"string":    "CITEXT",
				"int":       "INTEGER",
				"int32":     "INTEGER",
				"int64":     "BIGINT",
				"uint":      "BIGINT",
				"uint32":    "BIGINT",
				"uint64":    "BIGINT",
				"int16":     "SMALLINT",
				"int8":      "SMALLINT",
				"uint8":     "SMALLINT",
				"bool":      "BOOLEAN",
				"float32":   "REAL",
				"float64":   "DOUBLE PRECISION",
				"time.Time": "TIMESTAMP",
			},
			mapDefaultValue: map[string]string{
				"string":    "''",
				"int":       "0",
				"int32":     "0",
				"int64":     "0",
				"uint":      "0",
				"uint32":    "0",
				"uint64":    "0",
				"int16":     "0",
				"int8":      "0",
				"uint8":     "0",
				"bool":      "FALSE",
				"float32":   "0",
				"float64":   "0",
				"time.Time": "CURRENT_TIMESTAMP",
				"true":      "TRUE",
				"false":     "FALSE",
				"now()":     "CURRENT_TIMESTAMP",
			},
		},
	}
}
