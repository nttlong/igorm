package migrate

import (
	"dbv/tenantDB"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mssql_loader(t *testing.T) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=10000&encrypt=disable"
	db, err := tenantDB.Open("mssql", sqlServerDns)

	assert.NoError(t, err)

	migrator, err := NewMigrator(db)
	assert.NoError(t, err)
	tables, err := migrator.GetSqlCreateTable(reflect.TypeOf(User{}))
	assert.NoError(t, err)
	assert.NotEmpty(t, tables)

}
func TestLoadFK(t *testing.T) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=10000&encrypt=disable"
	db, err := tenantDB.Open("mssql", sqlServerDns)
	assert.NoError(t, err)

	loader := &MigratorLoaderMssql{}
	lst, err := loader.LoadForeignKey(db)
	assert.NoError(t, err)
	assert.NotEmpty(t, lst)

}
