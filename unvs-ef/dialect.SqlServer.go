package unvsef

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type SqlServerDialect struct {
	baseDialect
}

func (d *SqlServerDialect) Func(name string, args ...Expr) Expr {
	return rawFunc{name, args}
}

/*
Implementation of Dialect
Example: input "aaa","bbb"->[aaa].[bbb]
*/
func (d *SqlServerDialect) QuoteIdent(table, column string) string {
	return fmt.Sprintf(`[%s].[%s]`, table, column)
}

/*
Must be call after RefreshSchemaCache
*/
func (d *SqlServerDialect) TableExists(dbName, name string) bool {
	if _, ok := d.schema[dbName]; !ok {
		return false
	}
	tbl, ok := d.schema[dbName][strings.ToLower(name)]
	return ok && tbl.Name != ""
}

/*
Must be call after RefreshSchemaCache
*/
func (d *SqlServerDialect) ColumnExists(dbName, table, column string) bool {
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
Must be call after RefreshSchemaCache
*/
func (d *SqlServerDialect) SchemaMap(dbName string) map[string]TableSchema {
	if _, ok := d.schema[dbName]; !ok {
		return nil
	}
	return d.schema[dbName]
}

/*
The method must perform the following tasks:

	Retrieve all tables and columns in the SQL Server database.

	Retrieve all unique constraints in the SQL Server database.

	Retrieve all index constraints in the SQL Server database.

#Note:

	A unique constraint refers to a constraint created by the statement CREATE UNIQUE NONCLUSTERED.
*/
func (d *SqlServerDialect) RefreshSchemaCache(db *sql.DB, dbName string) error {

	schema := map[string]TableSchema{}

	// Load columns
	rows, err := db.Query(`
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
    t.name AS TableName,
    i.name AS IndexName
	FROM sys.indexes i
	JOIN sys.tables t ON i.object_id = t.object_id
	WHERE i.type_desc = 'NONCLUSTERED' and is_unique_constraint=1
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
		tbl := schema[key]
		tbl.IndexConstraints = append(tbl.IndexConstraints, index)
		schema[key] = tbl
	}
	d.schema[dbName] = schema

	return nil
}
func (d *SqlServerDialect) GetSchema(db *sql.DB, dbName string) (map[string]TableSchema, error) {
	if _, ok := d.schema[dbName]; !ok {
		if err := d.RefreshSchemaCache(db, dbName); err != nil {
			d.schemaError = err
			return nil, err
		}
	}

	return d.schema[dbName], nil

}

// generate sql create table for struct
func (d *SqlServerDialect) GenerateCreateTableSql(dbName string, typ reflect.Type) (string, error) {

	meta := utils.GetMetaInfo(typ)
	for tableName, fields := range meta {

		if d.TableExists(dbName, tableName) {
			return "", nil
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
			// if field.Unique {
			// 	colDef = append(colDef, "UNIQUE") kh check unique constraint, unique index se duoc tao sau
			// }
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
func (d *SqlServerDialect) GetPkConstraint(typ reflect.Type) (string, error) {
	mapConstraints := utils.GetPkFromMetaByType(typ)
	tableName := utils.TableNameFromStruct(typ)

	cols := []string{}

	for _, constraints := range mapConstraints {

		for fieldName := range constraints {
			cols = append(cols, fieldName)
		}
		//CONSTRAINT PK_MyTable PRIMARY KEY (Col1, Col2)
		//CONSTRAINT PK_MyTable PRIMARY KEY (Col1)
		constraintName := "primary____" + tableName + "___" + strings.Join(cols, "__")
		ret := fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", constraintName, strings.Join(cols, ", "))
		return ret, nil
	}
	return "", fmt.Errorf("no primary key found for table %s", utils.TableNameFromStruct(typ))
}
func (d *SqlServerDialect) GenerateAlterTableSql(dbName string, typ reflect.Type) ([]string, error) {

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
func (d *SqlServerDialect) GenerateUniqueConstraintsSql(typ reflect.Type) map[string]string {
	tableName := utils.TableNameFromStruct(typ)
	tableName = utils.Quote("[]", tableName)
	mapUniqueConstraints := utils.GetUniqueConstraintsFromMetaByType(typ)
	var ret = map[string]string{}
	for constraintName, constraints := range mapUniqueConstraints {
		//CREATE UNIQUE NONCLUSTERED INDEX idx_email ON [users] ([email]) WHERE [email] IS NOT NULL;
		cols := []string{}
		wheres := []string{}

		for fieldName := range constraints {
			cols = append(cols, utils.Quote("[]", fieldName))

			wheres = append(wheres, fmt.Sprintf("%s IS NOT NULL", utils.Quote("[]", fieldName)))

		}
		key := constraintName
		constraintName = utils.Quote("[]", constraintName)

		//CREATE UNIQUE NONCLUSTERED INDEX idx_email ON [users] ([email]) WHERE [email] IS NOT NULL;
		ret[key] = fmt.Sprintf("CREATE UNIQUE NONCLUSTERED INDEX %s ON %s (%s) WHERE %s", constraintName, tableName, strings.Join(cols, ", "), strings.Join(wheres, " AND "))
		// ret = append(ret, fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s UNIQUE (%s)", tableName, constraintName, constraint.Field.Name))
	}
	return ret

}
func (d *SqlServerDialect) GenerateIndexConstraintsSql(typ reflect.Type) map[string]string {
	tableName := utils.TableNameFromStruct(typ)
	tableName = utils.Quote("[]", tableName)
	mapIndexConstraints := utils.GetIndexConstraintsFromMetaByType(typ)
	var ret = map[string]string{}
	for constraintName, constraints := range mapIndexConstraints {
		cols := []string{}
		for fieldName := range constraints {
			cols = append(cols, utils.Quote("[]", fieldName))
		}
		key := constraintName
		constraintName = utils.Quote("[]", constraintName)
		ret[key] = fmt.Sprintf("CREATE NONCLUSTERED INDEX %s ON %s (%s)", constraintName, tableName, strings.Join(cols, ", "))
	}
	return ret

}

func NewSqlServerDialect(db *sql.DB) Dialect {
	return &SqlServerDialect{

		baseDialect: baseDialect{
			schema: map[string]map[string]TableSchema{},
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
		},
	}
}
