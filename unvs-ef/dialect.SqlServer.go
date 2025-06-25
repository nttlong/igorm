package unvsef

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type SqlServerDialect struct {
	DB              *sql.DB
	schema          map[string]TableSchema
	mapGoTypeToDb   map[string]string
	mapDefaultValue map[string]string
}

func (d *SqlServerDialect) Func(name string, args ...Expr) Expr {
	return rawFunc{name, args}
}

func (d *SqlServerDialect) QuoteIdent(table, column string) string {
	return fmt.Sprintf(`[%s].[%s]`, table, column)
}

func (d *SqlServerDialect) TableExists(name string) bool {
	tbl, ok := d.schema[strings.ToLower(name)]
	return ok && tbl.Name != ""
}

func (d *SqlServerDialect) ColumnExists(table, column string) bool {
	tbl, ok := d.schema[strings.ToLower(table)]
	if !ok {
		return false
	}
	_, colExists := tbl.Columns[strings.ToLower(column)]
	return colExists
}

func (d *SqlServerDialect) SchemaMap() map[string]TableSchema {
	return d.schema
}

func (d *SqlServerDialect) RefreshSchemaCache() error {
	d.schema = map[string]TableSchema{}

	// Load columns
	rows, err := d.DB.Query(`
		SELECT TABLE_NAME, COLUMN_NAME, DATA_TYPE, IS_NULLABLE
		FROM INFORMATION_SCHEMA.COLUMNS
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
		SELECT tc.TABLE_NAME, tc.CONSTRAINT_NAME
		FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS tc
		WHERE tc.CONSTRAINT_TYPE = 'UNIQUE'
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
		tbl := d.schema[key]
		tbl.UniqueConstraints = append(tbl.UniqueConstraints, constraint)
		d.schema[key] = tbl
	}

	// Load index constraints
	indexRows, err := d.DB.Query(`
		SELECT t.name AS table_name, i.name AS index_name
		FROM sys.indexes i
		JOIN sys.tables t ON t.object_id = i.object_id
		WHERE i.is_primary_key = 0 AND i.is_unique_constraint = 0 AND i.name IS NOT NULL
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
		tbl := d.schema[key]
		tbl.IndexConstraints = append(tbl.IndexConstraints, index)
		d.schema[key] = tbl
	}

	return nil
}

// generate sql create table for struct
func (d *SqlServerDialect) GenerateCreateTableSQL(typ reflect.Type) (string, error) {
	meta := utils.GetMetaInfo(typ)
	for tableName, fields := range meta {
		if d.TableExists(tableName) {
			return "", fmt.Errorf("table %s already exists", tableName)
		}

		var colDefs []string
		for colName, field := range fields {
			colType := field.DBType
			if _colType, ok := d.mapGoTypeToDb[utils.ResolveFieldKind(field.Field)]; ok {
				colType = _colType
			}

			if field.DBType != "" {
				colType = field.DBType
			} else if field.Length != nil {
				colType = fmt.Sprintf("NVARCHAR(%d)", *field.Length)
			}

			colDef := []string{utils.Quote("[]", colName), colType}

			if field.AutoIncrement {
				colDef = append(colDef, "IDENTITY(1,1)")
			}
			if field.Unique {
				colDef = append(colDef, "UNIQUE")
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
		stmt := fmt.Sprintf("CREATE TABLE %s (%s %s)", utils.Quote("[]", tableName), strings.Join(colDefs, ", "), pkConstraint)
		return stmt, nil
	}
	return "", nil
}
func (d SqlServerDialect) GetPkConstraint(typ reflect.Type) (string, error) {
	mapConstraints := utils.GetPkFromMetaByType(typ)

	cols := []string{}
	constraintName := ""
	for _constraintName, constraints := range mapConstraints {
		constraintName = _constraintName
		for fieldName := range constraints {
			cols = append(cols, fieldName)
		}
		//CONSTRAINT PK_MyTable PRIMARY KEY (Col1, Col2)
		//CONSTRAINT PK_MyTable PRIMARY KEY (Col1)
		ret := fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", constraintName, strings.Join(cols, ", "))
		return ret, nil
	}
	return "", fmt.Errorf("no primary key found for table %s", utils.TableNameFromStruct(typ))
}
func (d SqlServerDialect) GenerateAlterTableSQL(typ reflect.Type) ([]string, error) {

	tableName := utils.TableNameFromStruct(typ)

	meta := utils.GetMetaInfo(typ)
	var metaField map[string]FieldTag
	if _metaField, ok := meta[tableName]; ok {
		metaField = _metaField

	} else {
		return nil, fmt.Errorf("table %s not found in meta", tableName)
	}

	//schema := d.SchemaMap()[tableName]
	var alters []string

	for fieldName, field := range metaField {

		if d.ColumnExists(tableName, fieldName) {
			continue
		}
		colType := field.DBType
		if _colType, ok := d.mapGoTypeToDb[utils.ResolveFieldKind(field.Field)]; ok {
			colType = _colType
		}
		if field.DBType != "" {
			colType = field.DBType
		} else if field.Length != nil {
			colType = fmt.Sprintf("VARCHAR(%d)", *field.Length)
		}

		colDef := []string{utils.Quote("[]", fieldName), colType}
		if field.PrimaryKey {
			continue // primary key cannot be added
		}
		if field.AutoIncrement {
			colDef = append(colDef, "IDENTITY(1,1)")
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

		stmt := fmt.Sprintf("ALTER TABLE %s ADD %s", tableName, strings.Join(colDef, " "))
		alters = append(alters, stmt)
	}

	return alters, nil
}
func NewSqlServerDialect(db *sql.DB) *SqlServerDialect {
	return &SqlServerDialect{
		DB:     db,
		schema: map[string]TableSchema{},
		mapGoTypeToDb: map[string]string{
			"string":    "NVARCHAR(MAX)",
			"int":       "INT",
			"int32":     "INT",
			"int64":     "BIGINT",
			"uint":      "BIGINT",
			"uint32":    "BIGINT",
			"uint64":    "BIGINT",
			"int16":     "SMALLINT",
			"int8":      "TINYINT",
			"uint8":     "TINYINT",
			"bool":      "BIT",
			"float32":   "REAL",
			"float64":   "FLOAT",
			"time.Time": "DATETIME2",
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
			"bool":      "0",
			"float32":   "0",
			"float64":   "0",
			"time.Time": "now()",
			"true":      "1",
			"false":     "0",
			"now()":     "GETDATE()",
		},
	}
}
