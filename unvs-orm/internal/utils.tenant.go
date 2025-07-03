package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func (u *utilsPackage) getTenantDb(db *sql.DB, typ reflect.Type) (*TenantDb, error) {
	dbType, dbName, err := u.DetectDatabaseType(db)
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%s_%s", dbType, dbName)
	//check from cache
	if val, ok := u.cacheGetTenantDb.Load(key); ok {
		return val.(*TenantDb), nil

	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.Anonymous {
			continue
		}
		var baseType reflect.Type
		if field.Type.Kind() == reflect.Ptr {
			baseType = reflect.TypeOf(&TenantDb{})
		} else {
			baseType = reflect.TypeOf(TenantDb{})
		}
		if baseType.String() == field.Type.String() {
			_dbSchema, err := u.NewTenantDb(db)
			if err != nil {
				return nil, err
			}
			u.cacheGetTenantDb.Store(key, _dbSchema)
			return _dbSchema, nil

		} else {
			_dbSchema, err := u.getTenantDb(db, field.Type)
			if err != nil {
				return nil, err
			} else if _dbSchema != nil {
				u.cacheGetTenantDb.Store(key, _dbSchema)
				return _dbSchema, nil
			} else {
				continue
			}
		}
	}
	return nil, nil
}
func (u *utilsPackage) DetectDatabaseType(db *sql.DB) (DBType, string, error) {
	var version string
	var comment string

	queries := []struct {
		query string
	}{
		{"SELECT 'version_comment' ,version();"},       // PostgreSQL,  Cockroach, Greenplum
		{"SELECT 'version_comment', @@VERSION;"},       // SQL Server, Sybase
		{"SELECT 'version_comment',sqlite_version();"}, // SQLite
		{"SELECT 'version_comment',tidb_version();"},   // TiDB
		{"SELECT * FROM v$version"},                    // Oracle
		{"SHOW VARIABLES LIKE 'version_comment'"},      //MySQL
	}

	for _, q := range queries {
		err := db.QueryRow(q.query).Scan(&comment, &version)
		if err == nil && version != "" {
			v := strings.ToLower(version)

			switch {
			case strings.Contains(v, "postgres"):
				if strings.Contains(v, "greenplum") {
					return DBGreenplum, version, nil
				}
				return DBPostgres, version, nil
			case strings.Contains(v, "cockroach"):
				return DBCockroach, version, nil
			case strings.Contains(v, "mysql"):
				if strings.Contains(v, "mariadb") {
					return DBMariaDB, version, nil
				}
				return DBMySQL, version, nil
			case strings.Contains(v, "mariadb"):
				return DBMariaDB, version, nil
			case strings.Contains(v, "microsoft") || strings.Contains(v, "sql server"):
				return DBMSSQL, version, nil
			case strings.Contains(v, "sqlite"):
				return DBSQLite, version, nil
			case strings.Contains(v, "tidb"):
				return DBTiDB, version, nil
			case strings.Contains(v, "oracle"):
				return DBOracle, version, nil
			}
		}
	}

	return DBUnknown, version, errors.New("unable to detect database type")
}
func (u *utilsPackage) NewTenantDb(db *sql.DB) (*TenantDb, error) {

	key := fmt.Sprintf("NewTenantDb %v", db)
	if val, ok := u.cacheNewTenantDb.Load(key); ok {
		return val.(*TenantDb), nil
	}
	ret := &TenantDb{
		DB: db,
	}

	dbDetect, dbTypeName, err := u.DetectDatabaseType(db)

	if err != nil {
		return nil, err
	}
	dbName, err := u.GetCurrentDatabaseName(db, dbDetect)
	if err != nil {
		return nil, err
	}
	ret.DbName = dbName
	if dbDetect == DBMSSQL {
		ret.Dialect = NewMssqlDialect()
		ret.DBType = DBMSSQL
		ret.DBTypeName = dbTypeName
	} else if dbDetect == DBMySQL {
		ret.Dialect = NewMysqlDialect()
		ret.DBType = DBMySQL
		ret.DBTypeName = dbTypeName
	} else if dbDetect == DBPostgres {
		ret.Dialect = NewPostgresDialect()
		ret.DBType = DBPostgres
		ret.DBTypeName = dbTypeName
	} else {
		return nil, fmt.Errorf("Unsupported database type '%s'", dbTypeName)

	}
	u.cacheNewTenantDb.Store(key, ret)
	return ret, nil
}
func (u *utilsPackage) GetDbName(db *sql.DB) (string, error) {
	dbType, _, err := u.DetectDatabaseType(db)
	if err != nil {
		return "", err
	}
	return u.GetCurrentDatabaseName(db, dbType)
}
func (u *utilsPackage) GetCurrentDatabaseName(db *sql.DB, dbType DBType) (string, error) {
	var query string
	var dbName string

	switch dbType {
	case DBPostgres, DBGreenplum, DBCockroach:
		query = "SELECT current_database();"
	case DBMySQL, DBMariaDB, DBTiDB:
		query = "SELECT DATABASE();"
	case DBMSSQL:
		query = "SELECT DB_NAME();"
	case DBSQLite:
		query = "PRAGMA database_list;" // SQLite đặc biệt hơn, xem dưới
	case DBOracle:
		query = "SELECT SYS_CONTEXT('USERENV','DB_NAME') FROM dual;"
	default:
		return "", fmt.Errorf("unsupported db type: %s", dbType)
	}

	if dbType == DBSQLite {
		type sqliteEntry struct {
			Seq  int
			Name string
			File string
		}
		rows, err := db.Query(query)
		if err != nil {
			return "", err
		}
		defer rows.Close()

		for rows.Next() {
			var entry sqliteEntry
			if err := rows.Scan(&entry.Seq, &entry.Name, &entry.File); err != nil {
				return "", err
			}
			if entry.Seq == 0 {
				return entry.Name, nil // thường là "main"
			}
		}
		return "", fmt.Errorf("no database found in sqlite PRAGMA list")
	}

	// Các DB bình thường chỉ cần query trả 1 giá trị
	err := db.QueryRow(query).Scan(&dbName)
	if err != nil {
		return "", err
	}

	return dbName, nil
}
