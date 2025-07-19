package vdb

import (
	"database/sql"
	"strings"
)

func (d *postgresDialect) NewDataBase(db *sql.DB, sampleDsn string, dbName string) (string, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`
	err := db.QueryRow(query, dbName).Scan(&exists)
	if err != nil {
		return "", err
	}
	if !exists {
		_, err := db.Exec(`CREATE DATABASE a001`)
		if err != nil {
			return "", err
		}
	}
	items := strings.Split(sampleDsn, "?")
	if len(items) > 1 {
		dsn := items[0] + "/" + dbName + "?" + items[1]
		return dsn, nil
	} else {
		dsn := items[0] + "/" + dbName
		return dsn, nil
	}

}
