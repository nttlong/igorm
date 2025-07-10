package eorm

import "database/sql"

// Certain functions are transformed during compilation based on the SQL dialect.
// The dialect decides whether a function needs to be adapted for the target database driver.
// If required, the function will be rewritten here.
type DialectDelegateFunction struct {
	FuncName         string
	Args             []string
	HandledByDialect bool // ✅ Indicates if this function is allowed to be delegated to the dialect
}
type Dialect interface {
	Name() string
	Quote(str ...string) string
	GetTableAndColumnsDictionary(db *sql.DB) (map[string]string, error)
	ToText(value string) string
	SqlFunction(delegator *DialectDelegateFunction) (string, error)
}
