package tenantDB

import (
	"database/sql"
	"errors"
	"strings"
	"sync"

	_ "github.com/microsoft/go-mssqldb"
)

type TenantDB struct {
	*sql.DB
	dbName     string
	driverName string
	DbType     DB_DRIVER_TYPE

	Version     string
	hasDetected bool
}
type DB_DRIVER_TYPE string

const (
	DB_DRIVER_Postgres  DB_DRIVER_TYPE = "postgres"
	DB_DRIVER_MySQL     DB_DRIVER_TYPE = "mysql"
	DB_DRIVER_MariaDB   DB_DRIVER_TYPE = "mariadb"
	DB_DRIVER_MSSQL     DB_DRIVER_TYPE = "sqlserver"
	DB_DRIVER_SQLite    DB_DRIVER_TYPE = "sqlite"
	DB_DRIVER_Oracle    DB_DRIVER_TYPE = "oracle"
	DB_DRIVER_TiDB      DB_DRIVER_TYPE = "tidb"
	DB_DRIVER_Cockroach DB_DRIVER_TYPE = "cockroach"
	DB_DRIVER_Greenplum DB_DRIVER_TYPE = "greenplum"
	DB_DRIVER_Unknown   DB_DRIVER_TYPE = "unknown"
)

type DetectInfo struct {
}

type DbDetector struct {
	cacheDetectDatabaseType sync.Map
}

func (db *TenantDB) Detect() error {
	if db.hasDetected {
		return nil
	}

	var version string
	var dbName string

	queries := []struct {
		query string
	}{
		{"SELECT version();"},        // PostgreSQL, MySQL, Cockroach, Greenplum
		{"SELECT @@VERSION;"},        // SQL Server, Sybase
		{"SELECT sqlite_version();"}, // SQLite
		{"SELECT tidb_version();"},   // TiDB
		{"SELECT * FROM v$version"},  // Oracle
	}

	var dbType DB_DRIVER_TYPE = DB_DRIVER_Unknown

	for _, q := range queries {
		err := db.QueryRow(q.query).Scan(&version)
		if err == nil && version != "" {
			v := strings.ToLower(version)

			switch {
			case strings.Contains(v, "postgres"):
				dbType = DB_DRIVER_Postgres
				if strings.Contains(v, "greenplum") {
					dbType = DB_DRIVER_Greenplum
				}
				err = db.QueryRow("SELECT current_database();").Scan(&dbName)
			case strings.Contains(v, "cockroach"):
				dbType = DB_DRIVER_Cockroach
				err = db.QueryRow("SELECT current_database();").Scan(&dbName)
			case strings.Contains(v, "mysql"):
				dbType = DB_DRIVER_MySQL
				if strings.Contains(v, "mariadb") {
					dbType = DB_DRIVER_MariaDB
				}
				err = db.QueryRow("SELECT DATABASE();").Scan(&dbName)
			case strings.Contains(v, "mariadb"):
				dbType = DB_DRIVER_MariaDB
				err = db.QueryRow("SELECT DATABASE();").Scan(&dbName)
			case strings.Contains(v, "microsoft"), strings.Contains(v, "sql server"):
				dbType = DB_DRIVER_MSSQL
				err = db.QueryRow("SELECT DB_NAME();").Scan(&dbName)
			case strings.Contains(v, "sqlite"):
				dbType = DB_DRIVER_SQLite
				dbName = "main"
			case strings.Contains(v, "tidb"):
				dbType = DB_DRIVER_TiDB
				err = db.QueryRow("SELECT DATABASE();").Scan(&dbName)
			case strings.Contains(v, "oracle"):
				dbType = DB_DRIVER_Oracle
				err = db.QueryRow("SELECT SYS_CONTEXT('USERENV', 'DB_NAME') FROM dual").Scan(&dbName)
			}

			if err != nil {
				dbName = ""
			}
			db.DbType = dbType
			db.dbName = dbName
			db.Version = version
			db.hasDetected = true

			return nil
		}
	}

	return errors.New("unable to detect database type")
}
func Open(driverName string, dataSourceName string) (*TenantDB, error) {
	DB, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	return &TenantDB{
		DB:         DB,
		driverName: driverName,
	}, nil
}
