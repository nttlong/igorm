package dbv

import (
	"database/sql"
	"dbv/migrate"
	"fmt"
	"strings"
)

// Certain functions are transformed during compilation based on the SQL dialect.
// The dialect decides whether a function needs to be adapted for the target database driver.
// If required, the function will be rewritten here.
type DialectDelegateFunction struct {
	FuncName         string
	Args             []string
	HandledByDialect bool // ✅ Indicates if this function is allowed to be delegated to the dialect
}
type DIALECT_DB_ERROR_TYPE = int

const (
	DIALECT_DB_ERROR_TYPE_UNKNOWN DIALECT_DB_ERROR_TYPE = iota
	DIALECT_DB_ERROR_TYPE_DUPLICATE
	DIALECT_DB_ERROR_TYPE_REFERENCES // ✅ refrences_violation
	DIALECT_DB_ERROR_TYPE_REQUIRED
)

func ErrorMessage(t DIALECT_DB_ERROR_TYPE) string {
	switch t {
	case DIALECT_DB_ERROR_TYPE_UNKNOWN:
		return "unknown"
	case DIALECT_DB_ERROR_TYPE_DUPLICATE:
		return "duplicate"
	case DIALECT_DB_ERROR_TYPE_REFERENCES:
		return "references"
	case DIALECT_DB_ERROR_TYPE_REQUIRED:
		return "required"
	default:
		return "unknown"
	}
}
func (e DialectError) Error() string {
	return fmt.Sprintf("code=%s, %s: %s cols %v tables %v, entity fields %v", e.Code, ErrorMessage(e.ErrorType), e.ErrorMessage, strings.Join(e.DbCols, ","), strings.Join(e.Tables, ","), strings.Join(e.Fields, ","))
}
func (e DialectError) Unwrap() error {
	return e.Err
}

type DialectError struct {
	Err          error
	ErrorType    DIALECT_DB_ERROR_TYPE
	Code         string
	ErrorMessage string
	DbCols       []string
	Fields       []string
	Tables       []string
}

func (e *DialectError) Reload() {
	e.Code = "ERR" + fmt.Sprintf("%04d", e.ErrorType)
	e.ErrorMessage = ErrorMessage(e.ErrorType)
}

type Dialect interface {
	ParseError(err error) error
	Name() string
	Quote(str ...string) string
	GetTableAndColumnsDictionary(db *sql.DB) (map[string]string, error)
	ToText(value string) string
	SqlFunction(delegator *DialectDelegateFunction) (string, error)
	MakeSqlInsert(tableName string, columns []migrate.ColumnDef, data interface{}) (string, []interface{})
}
