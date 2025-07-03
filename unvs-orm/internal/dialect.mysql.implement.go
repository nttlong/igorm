package internal

import (
	"database/sql"
	"reflect"
)

func (d *MySqlDialect) GetParamPlaceholder() string {
	return d.paramPlaceholder
}

/*
	depends on bd driver type the function will be implement in

dialect.<driver name>.go
*/
func (d *MySqlDialect) BuildSqlInsert(TableName string, AutoKeyField string, fields ...string) string {
	panic("not implemented")
}
func (d *MySqlDialect) MakeLimitOffset(limit *int, offset *int) string {
	panic("not implemented")
}

/*
Example:

	QuoteIdent("table_name")
	=> [table_name] if dialect is MSSQL
	QuoteIdent("table_name","column_name")
	=> [table_name].[column_name] if dialect is MSSQL
*/
func (d *MySqlDialect) QuoteIdent(args ...string) string {
	panic("not implemented")
}

// Schema management methods
/*
	Check is table existing in Db, call after call RefreshSchemaCache
	Purpose: for Database migration only
*/
func (d *MySqlDialect) TableExists(dbName, name string) bool {
	panic("not implemented")
}

/*
Check is a column of a table  existing in Db, call after call RefreshSchemaCache
Purpose: for Database migration only
*/
func (d *MySqlDialect) ColumnExists(dbName, table string, column string) bool {
	panic("not implemented")
}

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
func (d *MySqlDialect) RefreshSchemaCache(db *sql.DB, dbName string) error {
	panic("not implemented")
}
func (d *MySqlDialect) GetSchema(db *sql.DB, dbName string) (map[string]TableSchema, error) {
	panic("not implemented")
}

/*
Get schema info from cache. Call after RefreshSchemaCache
*/
func (d *MySqlDialect) SchemaMap(dbName string) map[string]TableSchema {
	panic("not implemented")
}

/*
Based on the metadata information, the system will generate SQL statements to create Unique Constraints appropriately.

# Purpose:

	for Database migration only

# Note:

	1- Meta information is obtained by calling RefreshSchemaCache
	2- All constraint names are obtained by combining the table name, a double underscore ("__"), and the constraint name

# return map [constraint name] [ sql create constraint]
*/
func (d *MySqlDialect) GenerateUniqueConstraintsSql(typ reflect.Type) map[string]string {
	panic("not implemented")
}

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
func (d *MySqlDialect) GenerateIndexConstraintsSql(typ reflect.Type) map[string]string {
	panic("not implemented")
}

/*
# Purpose:

	for Database migration only

# Note:

	This function purely generates SQL statements to create a table,
	including only the primary key and columns,
	with additional table-level constraints such as Required or Nullable combined with Default.
*/
func (d *MySqlDialect) GenerateCreateTableSql(dbName string, typ reflect.Type) (string, error) {
	panic("not implemented")
}

/*
Generate all SQL ALTER TABLE ADD COLUMN statements for columns that exist in the model but are missing in the database.

# Purpose:

	for Database migration only

# Note:

	This method refer to GetPkConstraint
*/
func (d *MySqlDialect) GenerateAlterTableSql(dbName string, typ reflect.Type) ([]string, error) {
	panic("not implemented")
}

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
func (d *MySqlDialect) GetPkConstraint(typ reflect.Type) (string, error) {
	panic("not implemented")
}
