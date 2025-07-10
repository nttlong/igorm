package eorm

import (
	"database/sql"
	"fmt"
)

type dbDictionaryItem struct {
	tablesAndColumns map[string]string
}
type DbDictionaryReceiver struct {
	schemas map[string]map[string]dbDictionaryItem
}

func (d *DbDictionaryReceiver) GetTableAndColumnsDictionary(dbName string, driverName string) map[string]string {
	return d.schemas[driverName][dbName].tablesAndColumns
}
func (d *DbDictionaryReceiver) LoadTableAndColumnsDictionary(db *sql.DB, dbName string) error {
	panic(fmt.Errorf("not implemented in file eorm/db.dictionaries.go, line 19"))
}
