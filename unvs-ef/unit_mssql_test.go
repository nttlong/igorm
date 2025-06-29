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

func GetMysqlDb() *sql.DB {
	dsn := "root:123456@tcp(localhost:3306)/aaa"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return db
}
func createMssqlDb() *sql.DB {
	dsn := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		panic(err)
	}
	return db
}
func CreateMssqlRepo() (*Repository, error) {
	db := createMssqlDb()

	defer db.Close()
	ret, err := buildRepositoryFromStruct[Repository](db, true)
	return ret, err
}
func TestMssql(t *testing.T) {
	ret, err := CreateMssqlRepo()
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
}
