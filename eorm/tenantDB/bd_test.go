package tenantDB

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTenantDB(t *testing.T) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=10000&encrypt=disable"
	db, err := Open("mssql", sqlServerDns)

	assert.NoError(t, err)
	defer db.Close()
	err = db.Detect()
	assert.NoError(t, err)
	// do something with the db
}
