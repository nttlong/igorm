package eorm

import "database/sql"

type Dialect interface {
	Name() string
	Quote(str ...string) string
	GetTableAndColumnsDictionary(db *sql.DB) (map[string]string, error)
	ToText(value string) string
}
