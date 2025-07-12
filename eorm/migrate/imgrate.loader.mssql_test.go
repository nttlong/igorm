package migrate

import (
	"eorm/tenantDB"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mssql_loader(t *testing.T) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=10000&encrypt=disable"
	db, err := tenantDB.Open("mssql", sqlServerDns)
	assert.NoError(t, err)

	migrator, err := MigratorLoader(db)
	assert.NoError(t, err)
	tables, err := migrator.LoadAllTable(db.DB)
	assert.NoError(t, err)
	assert.NotEmpty(t, tables)
}
