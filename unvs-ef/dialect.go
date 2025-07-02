package unvsef

import (
	"database/sql"
	"reflect"
	"sync"
)

// TableSchema represents a table and its columns in the database.

type TableSchema struct {
	UniqueConstraints     []string
	IndexConstraints      []string
	ForeignKeyConstraints []string
	Name                  string
	Columns               map[string]ColumnSchema
}
type baseDialect struct {
	schema          map[string]map[string]TableSchema
	mapGoTypeToDb   map[string]string
	mapDefaultValue map[string]string
	schemaOnce      sync.Once
	schemaError     error
}

// ColumnSchema holds metadata about a column.
// metadata obtains by exec SQL from db
type ColumnSchema struct {
	Name          string // real column name in db
	Type          string // real column name in db
	Nullable      bool
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	Length        *int
	Index         bool
	Comment       string
}

/*
#Dialect allows different SQL dialects (e.g., PostgreSQL, MSSQL, MySQL) to be supported.
*/
type Dialect interface {
	GetParamPlaceholder() string
	/* depends on bd driver type the function will be implement in
	dialect.<driver name>.go
	*/
	BuildSqlInsert(TableName string, AutoKeyField string, fields ...string) string
	MakeLimitOffset(limit *int, offset *int) string
	/*
		Example:
			QuoteIdent("table_name")
			=> [table_name] if dialect is MSSQL
			QuoteIdent("table_name","column_name")
			=> [table_name].[column_name] if dialect is MSSQL
	*/
	QuoteIdent(args ...string) string

	// Schema management methods
	/*
		Check is table existing in Db, call after call RefreshSchemaCache
		Purpose: for Database migration only
	*/
	TableExists(dbName, name string) bool // Kiểm tra xem bảng có tồn tại trong cơ sở dữ liệu không (Check if a table exists in the database)
	/*
		Check is a column of a table  existing in Db, call after call RefreshSchemaCache
		Purpose: for Database migration only
	*/
	ColumnExists(dbName, table string, column string) bool // Kiểm tra cột có tồn tại trong bảng không (Check if a column exists in a table)
	/*
		Gathering all info in a specific database
		The info is including:
			table and columns
			all indexes
			all unique key
			all foreign key
		Method also store that info in cache

		# Purpose: for Database migration only
	*/
	RefreshSchemaCache(db *sql.DB, dbName string) error // Tải lại toàn bộ schema từ DB vào bộ nhớ đệm (Reload schema metadata into cache)
	GetSchema(db *sql.DB, dbName string) (map[string]TableSchema, error)
	/*
		Get schema info from cache. Call after RefreshSchemaCache
	*/
	SchemaMap(dbName string) map[string]TableSchema
	/*
		Based on the metadata information, the system will generate SQL statements to create Unique Constraints appropriately.

		# Purpose:
			for Database migration only

		# Note:

			1- Meta information is obtained by calling RefreshSchemaCache
			2- All constraint names are obtained by combining the table name, a double underscore ("__"), and the constraint name
		# return map [constraint name] [ sql create constraint]


	*/
	GenerateUniqueConstraintsSql(typ reflect.Type) map[string]string
	/*
		Based on the metadata information, the system will generate SQL statements to create
		Index Constraints appropriately.
		# Purpose:
			for Database migration only
		# Note:

			1- Meta information is obtained by calling RefreshSchemaCache
			2- All constraint names are obtained by combining the table name, a double underscore ("__"), and the constraint name
		# return map [constraint name] [ sql create constraint]

	*/
	GenerateIndexConstraintsSql(typ reflect.Type) map[string]string

	/*

		# Purpose:
			for Database migration only
		# Note:
			This function purely generates SQL statements to create a table,
			including only the primary key and columns,
			with additional table-level constraints such as Required or Nullable combined with Default.

	*/
	GenerateCreateTableSql(dbName string, typ reflect.Type) (string, error)
	/*
		Generate all SQL ALTER TABLE ADD COLUMN statements for columns that exist in the model but are missing in the database.

		# Purpose:
			for Database migration only
		# Note:
			This method refer to GetPkConstraint

	*/
	GenerateAlterTableSql(dbName string, typ reflect.Type) ([]string, error)
	/*
		Generate all Primary Key Constraint . The method will be called by GenerateCreateTableSQL

		Example:
			struct UserRole Struct {
				UserId DbFlied[int] `db:primaryKey`
				RoleId DbFlied[int] `db:primaryKey`
			}
		# Requirement: Constraint name is "primary____user_roles___user_id__role_id"
		Constraint name Primary key is the combination of "primary", four underscores, table name triple underscore and list of filed name join by double underscore

	*/
	GetPkConstraint(typ reflect.Type) (string, error)
}
type BaseDialect struct {
	schema          map[string]TableSchema
	schemaError     error
	mapGoTypeToDb   map[string]string
	mapDefaultValue map[string]string
	schemaOnce      sync.Once
}

// rawFunc allows wrapping generic functions like FUNC(arg1, arg2, ...)
