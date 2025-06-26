package unvsef

import (
	"database/sql"
	"reflect"
)

type MySqlDialect struct {
	baseDialect
}

func (d *MySqlDialect) Func(name string, args ...Expr) Expr {
	panic("unimplemented")

}

/*
	depends on bd driver type the function will be implement in

dialect.<driver name>.go
*/
func (d *MySqlDialect) QuoteIdent(table, column string) string {
	panic("unimplemented")

}

// Schema management methods
/*
	Check is table existing in Db, call after call RefreshSchemaCache
	Purpose: for Database migration only
*/
func (d *MySqlDialect) TableExists(dbName, name string) bool {
	panic("unimplemented")

}

/*
Check is a column of a table  existing in Db, call after call RefreshSchemaCache
Purpose: for Database migration only
*/
func (d *MySqlDialect) ColumnExists(dbName, table string, column string) bool {
	panic("unimplemented")

}

/*
Gathering all info in a specific database
The info is including:

	table and columns
	all indexes
	all unique key
	all foreign key

# Method also store that info in cache

# Purpose: for Database migration only
*/
func (d *MySqlDialect) RefreshSchemaCache(db *sql.DB, dbName string) error {
	panic("unimplemented")

}
func (d *MySqlDialect) GetSchema(db *sql.DB, dbName string) (map[string]TableSchema, error)

/*
Get schema info from cache. Call after RefreshSchemaCache
*/
func (d *MySqlDialect) SchemaMap(dbName string) map[string]TableSchema

/*
Based on the metadata information, the system will generate SQL statements to create Unique Constraints appropriately.

# Purpose:

	for Database migration only

# Note:

	1- Meta information is obtained by calling RefreshSchemaCache
	2- All constraint names are obtained by combining the table name, a double underscore ("__"), and the constraint name

# return map [constraint name] [ sql create constraint]
*/
func (d *MySqlDialect) GenerateUniqueConstraintsSql(typ reflect.Type) map[string]string

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
func (d *MySqlDialect) GenerateIndexConstraintsSql(typ reflect.Type) map[string]string

/*
# Purpose:

	for Database migration only

# Note:

	This function purely generates SQL statements to create a table,
	including only the primary key and columns,
	with additional table-level constraints such as Required or Nullable combined with Default.
*/
func (d *MySqlDialect) GenerateCreateTableSql(dbName string, typ reflect.Type) (string, error)

/*
Generate all SQL ALTER TABLE ADD COLUMN statements for columns that exist in the model but are missing in the database.

# Purpose:

	for Database migration only

# Note:

	This method refer to GetPkConstraint
*/
func (d *MySqlDialect) GenerateAlterTableSql(dbName string, typ reflect.Type) ([]string, error) {
	panic("unimplemented")

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
	panic("unimplemented")

}
