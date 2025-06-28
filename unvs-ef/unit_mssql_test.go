package unvsef

import (
	"database/sql"
	"testing"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/ibmdb/go_ibm_db"
	_ "github.com/lib/pq"

	_ "github.com/sijms/go-ora/v2"
	// ef "unvs.ef"
)

func createMssqlDb() *sql.DB {
	dsn := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		panic(err)
	}
	return db
}

func TestMssql(t *testing.T) {
	db := createMssqlDb()
	defer db.Close()
}
