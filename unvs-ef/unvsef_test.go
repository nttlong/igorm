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

type SysUser struct {
	Id          DbField[uint64]  `db:"primaryKey;autoIncrement"`
	Code        DbField[string]  `db:"length(50);primaryKey"`
	Email       DbField[string]  `db:"length(50);unique"`
	Description DbField[*string] `db:"length(200)"`
}

func TestSQl(t *testing.T) {
	n := utils.TableNameFromStruct(reflect.TypeOf(SampleModel{}))
	assert.Equal(t, "custom_users", n)
	dsn := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	d := NewSqlServerDialect(db)
	u := Entity[SysUser]()
	query := NewQuery().
		Select(d.Func("LEN", u.Email)).
		From(FromStruct[SysUser]())
	sql, _ := query.ToSQL(d)
	fmt.Println(sql)
}
func TestInsert(t *testing.T) {
	dsn := "user=postgres password=123456 host=localhost port=5432 dbname=fx001 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	d := NewSqlServerDialect(db)
	if err != nil {
		t.Fatal(err)
	}
	test := "ccccc"
	u := &SysUser{}
	u.Id.Set(1)
	u.Code.Set("U001")
	u.Email.Set("abc@test.com")
	u.Description.Set(&test)

	sql, args := Insert(u).ToSQL(d)
	fmt.Println(sql)
	fmt.Println(args)

}

type Article struct {
	Id DbField[uint64] `db:"primaryKey;autoIncrement"`

	Title   DbField[string] `db:"FTS(title_idx);length(50)"`
	Content DbField[string] `db:"FTS(content_idx);length(50)"`
}
type Comment struct {
	Id        DbField[uint64] `db:"primaryKey;autoIncrement"`
	ArticleId DbField[uint64] `db:"index"`
	Content   DbField[string] `db:"FTS(content_idx)"`
}

type Repository struct {
	TenantDb
	Articles *Article
	Comments *Comment
}

func TestBuildEntity(t *testing.T) {
	dsn := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	r, err := Repo[Repository](db, true)

	assert.NoError(t, err)
	err = r.DoMigrate()
	assert.NoError(t, err)
	qr := r.Query().From(FromStruct[Article]()).Select(r.Articles.Id, r.Articles.Title)
	sql, params := qr.SQLCommand()
	fmt.Println(sql)
	fmt.Println(params)
	// NewQuery().From(FromStruct[Article]()).Select(r.Articles.Id, r.Articles.Title)

	fmt.Println(r)
}
func TestFTS(t *testing.T) {
	dsn := "user=postgres password=123456 host=localhost port=5432 dbname=fx001 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	a := Entity[Article]()
	Funcs.Dialect = NewSqlServerDialect(db)

	query := NewQuery().
		Select(a.Title, a.Content).
		From("articles").
		Where(Funcs.FullTextContains(a.Content, "machine learning"))

	sql, args := query.ToSQL(Funcs.Dialect)
	fmt.Println(sql)
	fmt.Println(args)
}
func TestSqlServer(t *testing.T) {
	dsn := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", dsn)

	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	d := &SqlServerDialect{}

	d.RefreshSchemaCache(db, "aaa")
	assert.NoError(t, err)
	t.Log(d.schema)
}
func TestPostgres(t *testing.T) {
	dsn := "user=postgres password=123456 host=localhost port=5432 dbname=fx001 sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	d := &PostgresDialect{}

	err = d.RefreshSchemaCache(db, "fx001")
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
	d := &MySqlDialect{}

	err = d.RefreshSchemaCache(db, "root")
	if err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)
	t.Log(d.schema)
}

type SampleModel struct {
	_         DbField[any]     `db:"table(table_test001)"`
	NullField DbField[*string] `db:"length(50);index(idx_test1)"`

	Id   DbField[uint64] `db:"primaryKey;autoIncrement"`
	Id2  DbField[uint64] `db:"primaryKey"`
	Name DbField[string] `db:"length(50);index"`
	Code DbField[string] `db:"length(50);unique"`
	Unk1 DbField[string] `db:"unique(uk1);length(50)"`
	Unk2 DbField[string] `db:"unique(uk2);length(50)"`

	Test2 DbField[string]   `db:"length(50);index(idx_test1)"`
	Test3 DbField[*float64] `db:"type:decimal(10,2)"`
	Test4 DbField[bool]     `db:"default:true"`
	Test5 DbField[*bool]
	Test6 DbField[time.Time]
	Test7 DbField[*time.Time]
	Test8 DbField[time.Time] `db:"default:now()"`
	Unk3  DbField[string]    `db:"unique(uk_test);length(50)"`
	Unk4  DbField[string]    `db:"unique(uk_test);length(50)"`
	Unk5  DbField[bool]      `db:"unique(uk_test)"`
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

	err = d.RefreshSchemaCache(db, "aaa")
	if err != nil {
		t.Fatal(err)
	}
	sql, err := d.GenerateCreateTableSql("aaa", reflect.TypeOf(SampleModel{}))
	assert.NoError(t, err)
	t.Log(sql)
	r, err := db.Exec(sql)
	assert.NoError(t, err)
	t.Log(r)
	sqls, err := d.GenerateAlterTableSql("aaa", reflect.TypeOf(SampleModel{}))
	assert.NoError(t, err)
	for _, sql := range sqls {
		t.Log(sql)
		r, err := db.Exec(sql)
		assert.NoError(t, err)
		t.Log(r)
	}
	sqlMap := d.GenerateUniqueConstraintsSql(reflect.TypeOf(SampleModel{}))
	for _, sql := range sqlMap {
		t.Log(sql)
		r, err := db.Exec(sql)
		assert.NoError(t, err)
		t.Log(r)
	}
	t.Log(sqls)

}
func TestMigrate(t *testing.T) {
	dsn := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	mssql := NewSqlServerDialect(db)
	et := reflect.TypeOf(SampleModel{})
	sqls, err := utils.GetScriptMigrate(db, "aaa", mssql, et)
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
