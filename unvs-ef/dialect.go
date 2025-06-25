package unvsef

import (
	"fmt"
	"reflect"
	"strings"
)

// TableSchema represents a table and its columns in the database.
type TableSchema struct {
	UniqueConstraints []string
	IndexConstraints  []string
	Name              string
	Columns           map[string]ColumnSchema
}

// ColumnSchema holds metadata about a column.
type ColumnSchema struct {
	Name          string
	Type          string
	Nullable      bool
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	Length        *int
	Index         bool
	Comment       string
}

// Dialect allows different SQL dialects (e.g., PostgreSQL, MSSQL, MySQL) to be supported.
type Dialect interface {
	Func(name string, args ...Expr) Expr
	QuoteIdent(table, column string) string

	// Schema management methods
	TableExists(name string) bool                  // Kiểm tra xem bảng có tồn tại trong cơ sở dữ liệu không (Check if a table exists in the database)
	ColumnExists(table string, column string) bool // Kiểm tra cột có tồn tại trong bảng không (Check if a column exists in a table)
	RefreshSchemaCache() error                     // Tải lại toàn bộ schema từ DB vào bộ nhớ đệm (Reload schema metadata into cache)
	SchemaMap() map[string]TableSchema             // Trả về toàn bộ bảng đã cached cùng cột của chúng (Return schema cache with all known tables/columns)
	UniqueConstraints(ttyp reflect.Type) []string  // Danh sách constraint UNIQUE trên bảng (List of UNIQUE constraints on a table)
	IndexConstraints(typ reflect.Type) []string    // Danh sách constraint INDEX trên bảng (List of INDEX constraints on a table)

	// Generate CREATE TABLE SQL
	GenerateCreateTableSQL(typ reflect.Type) (string, error)
	// Generate ALTER TABLE ADD COLUMN statements if fields are missing
	GenerateAlterTableSQL(typ reflect.Type) ([]string, error)
	GetPkConstraint(typ reflect.Type) (string, error) // Lấy ra constraint PRIMARY KEY của bảng (Get PRIMARY KEY constraint of a table)
}

// rawFunc allows wrapping generic functions like FUNC(arg1, arg2, ...)
type rawFunc struct {
	name string
	args []Expr
}

func (f rawFunc) ToSQL(d Dialect) (string, []interface{}) {
	parts := []string{}
	args := []interface{}{}
	for _, a := range f.args {
		sql, aargs := a.ToSQL(d)
		parts = append(parts, sql)
		args = append(args, aargs...)
	}
	return fmt.Sprintf("%s(%s)", f.name, strings.Join(parts, ", ")), args
}

// --------------------- Literal Expression ---------------------

func (l Literal[T]) ToSQL(d Dialect) (string, []interface{}) {
	return "?", []interface{}{l.Value}
}
