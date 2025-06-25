package unvsef

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/ibmdb/go_ibm_db"
	_ "github.com/lib/pq"

	_ "github.com/sijms/go-ora/v2"
	"github.com/stretchr/testify/assert"
	// ef "unvs.ef"
)

func TestSqlServer(t *testing.T) {
	dsn := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	d := &SqlServerDialect{
		DB: db,
	}

	d.RefreshSchemaCache()
	t.Log(d.schema)
}
func TestPostgres(t *testing.T) {
	dsn := "user=postgres password=123456 host=localhost port=5432 dbname=fx001 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	d := &PostgresDialect{
		DB: db,
	}

	err = d.RefreshSchemaCache()
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)
	t.Log(d.schema)
}
func TestMySql(t *testing.T) {
	dsn := "root:123456@tcp(localhost:3306)/root?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	d := &MySQLDialect{
		DB: db,
	}

	err = d.RefreshSchemaCache()
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)
	t.Log(d.schema)
}

type SampleModel struct {
	NullField DbField[*string] `db:"length(50);index(idx_test1)"`
	_         DbField[any]     `db:"table(custom_users_test3)"`
	Id        DbField[uint64]  `db:"primaryKey;autoIncrement"`
	Id2       DbField[uint64]  `db:"primaryKey"`
	Name      DbField[string]  `db:"length(50);index"`
	Code      DbField[string]  `db:"length(50);unique"`

	Test2 DbField[string]   `db:"length(50);index(idx_test1)"`
	Test3 DbField[*float64] `db:"type:decimal(10,2)"`
	Test4 DbField[bool]     `db:"default:true"`
	Test5 DbField[*bool]
	Test6 DbField[time.Time]
	Test7 DbField[*time.Time]
}

func TestGetMetaInfo(t *testing.T) {
	ret := utils.GetMetaInfo(reflect.TypeOf(SampleModel{}))
	for k, v := range ret {
		t.Log(k, v)
	}
}
func TestGetSQLCreate(t *testing.T) {
	ret := utils.GetMetaInfo(reflect.TypeOf(SampleModel{}))
	for k, v := range ret {
		t.Log(k, v)
	}
}
func TestSQLServerGenerateMakeTableSQL(t *testing.T) {
	n := utils.TableNameFromStruct(reflect.TypeOf(SampleModel{}))
	assert.Equal(t, "custom_users", n)
	dsn := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	d := NewSqlServerDialect(db)

	err = d.RefreshSchemaCache()
	if err != nil {
		t.Fatal(err)
	}
	sql, err := d.GenerateCreateTableSQL(reflect.TypeOf(SampleModel{}))
	assert.NoError(t, err)
	t.Log(sql)
	r, err := db.Exec(sql)
	assert.NoError(t, err)
	t.Log(r)
	sqls, err := d.GenerateAlterTableSQL(reflect.TypeOf(SampleModel{}))
	assert.NoError(t, err)
	for _, sql := range sqls {
		t.Log(sql)
		r, err := db.Exec(sql)
		assert.NoError(t, err)
		t.Log(r)
	}
	t.Log(sqls)

}
func TestToSnakeCase(t *testing.T) {
	testSample := map[string]string{
		"Id":        "id",
		"ID":        "id",
		"Name":      "name",
		"Code":      "code",
		"Test1":     "test1",
		"Test2":     "test2",
		"UserID":    "user_id",
		"UserId":    "user_id",
		"UserName":  "user_name",
		"UserName1": "user_name1",
	}

	for k, v := range testSample {
		r := utils.ToSnakeCase(k)
		assert.Equal(t, v, r)
	}
}
func TestResolveFieldKind(t *testing.T) {
	type TestStruct struct {
		_  DbField[any]    `db:"table(custom_users)"`
		Id DbField[uint64] `db:"primaryKey;autoIncrement"`
	}
	field := reflect.TypeOf(TestStruct{}).Field(0)
	kinde := utils.ResolveFieldKind(field)
	fmt.Println(kinde)
	// Output: uint64

	assert.Equal(t, "reflect.Uint64", kinde)

}
