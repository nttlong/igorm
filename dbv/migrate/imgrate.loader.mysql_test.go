package migrate

import (
	"dbv/tenantDB"
	"fmt"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestMySqlMigrate(t *testing.T) {
	// pg dsn
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/test"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)

	assert.NoError(t, err)

	migrator, err := NewMigrator(db)
	loader := migrator.GetLoader()
	pgLoader := loader.(*MigratorLoaderMysql)
	cols, err := pgLoader.LoadAllTable(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, cols)
	pks, err := pgLoader.LoadAllPrimaryKey(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, pks)
	uk, err := pgLoader.LoadAllUniIndex(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, uk)
	idx, err := pgLoader.LoadAllIndex(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, idx)
	fk, err := pgLoader.LoadForeignKey(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, fk)
	schema, err := pgLoader.LoadFullSchema(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, schema)
	mysqlMigrator := migrator.(*migratorMySql)
	tables, err := mysqlMigrator.GetSqlCreateTable(reflect.TypeOf(User{}))
	fmt.Println(tables)
	assert.NoError(t, err)
	assert.NotEmpty(t, tables)
}
func TestMysqlSqlAddColumns(t *testing.T) {
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/test"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)

	assert.NoError(t, err)

	migrator, err := NewMigrator(db)

	assert.NoError(t, err)

	assert.NoError(t, err)
	pgm := migrator.(*migratorMySql)
	sql, err := pgm.GetSqlAddColumn(reflect.TypeOf(User{}))
	assert.NoError(t, err)

	fmt.Println(sql)
	assert.NotEmpty(t, sql)
}
func TestMysqlAddIndex(t *testing.T) {
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/test"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)

	assert.NoError(t, err)

	migrator, err := NewMigrator(db)

	assert.NoError(t, err)
	pgm := migrator.(*migratorMySql)
	sql, err := pgm.GetSqlAddIndex(reflect.TypeOf(User{}))
	assert.NoError(t, err)

	fmt.Println(sql)
	assert.NotEmpty(t, sql)
}
func TestMysqlGetSqlAddUniqueIndex(t *testing.T) {
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/test"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)

	assert.NoError(t, err)

	migrator, err := NewMigrator(db)

	assert.NoError(t, err)
	pgm := migrator.(*migratorMySql)
	sql, err := pgm.GetSqlAddUniqueIndex(reflect.TypeOf(User{}))
	assert.NoError(t, err)

	fmt.Println(sql)
	assert.NotEmpty(t, sql)
}
func TestMysqlGetAddForeignKey(t *testing.T) {
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/test"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)

	assert.NoError(t, err)

	migrator, err := NewMigrator(db)

	assert.NoError(t, err)
	pgm := migrator.(*migratorMySql)
	sql, err := pgm.GetSqlAddForeignKey()
	assert.NoError(t, err)

	fmt.Println(sql)
	assert.NotEmpty(t, sql)
}
func BenchmarkMySqlLoadFullSchema(b *testing.B) {
	mySqlDsn := "root:123456@tcp(127.0.0.1:3306)/test"
	// create new migrate instance
	db, err := tenantDB.Open("mysql", mySqlDsn)

	assert.NoError(b, err)

	migrator, err := NewMigrator(db)

	assert.NoError(b, err)
	loader := migrator.GetLoader()
	pgLoader := loader.(*MigratorLoaderMysql)
	b.ResetTimer() // Reset timer để chỉ đo phần bên dưới
	for i := 0; i < b.N; i++ {

		schema, err := pgLoader.LoadFullSchema(db)
		assert.NoError(b, err)
		assert.NotEmpty(b, schema)

	}
}
