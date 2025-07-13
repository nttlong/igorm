package tenantDB

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"
	"sync"

	_ "github.com/microsoft/go-mssqldb"
)

type TenantDB struct {
	*sql.DB
	info *tenantDBInfo
}
type tenantDBInfo struct {
	dbName string

	driverName string
	DbType     DB_DRIVER_TYPE

	Version     string
	hasDetected bool
	key         string
}
type TenantTx struct {
	*sql.Tx
	info *tenantDBInfo
}

func (tx *TenantTx) GetDriverName() string {
	return tx.info.driverName
}
func (tx *TenantTx) GetDBName() string {
	return tx.info.dbName
}
func (tx *TenantTx) GetDbType() DB_DRIVER_TYPE {
	return tx.info.DbType
}
func (db *TenantDB) Begin() (*TenantTx, error) {
	db.Detect()
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &TenantTx{
		Tx:   tx,
		info: db.info,
	}, nil

}
func (db *TenantDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*TenantTx, error) {
	db.Detect()
	tx, err := db.DB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &TenantTx{
		Tx:   tx,
		info: db.info,
	}, nil
}

func (db *TenantDB) GetDriverName() string {
	return db.info.driverName
}
func (db *TenantDB) GetDBName() string {
	return db.info.dbName
}
func (db *TenantDB) GetDbType() DB_DRIVER_TYPE {
	return db.info.DbType
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

type DbDetector struct {
	cacheDetectDatabaseType sync.Map
}
type dbDetectInit struct {
	once sync.Once
	val  tenantDBInfo
	err  error
}

var cacheDbDetector sync.Map

func (info *tenantDBInfo) Detect(db *sql.DB) (*tenantDBInfo, error) {

	key := info.driverName + ":" + info.key
	actual, _ := cacheDbDetector.LoadOrStore(key, &dbDetectInit{})
	di := actual.(*dbDetectInit)
	di.once.Do(func() {
		di.val = tenantDBInfo{
			driverName: info.driverName,
			key:        info.key,
		}
		di.err = di.val.detect(db)
	})
	if di.err != nil {
		return nil, di.err
	} else {

		return &di.val, nil
	}
}

func (info *tenantDBInfo) detect(db *sql.DB) error {

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
			info.DbType = dbType
			info.dbName = dbName
			info.Version = version

			return nil
		}
	}

	return errors.New("unable to detect database type")
}
func Open(driverName, dsn string) (*TenantDB, error) {

	DB, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256([]byte(dsn))
	// Truncate nếu cần, ví dụ lấy 16 byte đầu (32 hex chars)
	key := hex.EncodeToString(hash[:16])
	info := &tenantDBInfo{
		driverName: driverName,
		key:        key,
	}

	info, err = info.Detect(DB)
	ret := &TenantDB{
		DB:   DB,
		info: info,
	}
	if err != nil {
		return nil, err
	}
	return ret, nil
}
func (db *TenantDB) Detect() error {
	info, err := db.info.Detect(db.DB)
	if err != nil {
		return err
	}
	db.info = info
	return nil
}
