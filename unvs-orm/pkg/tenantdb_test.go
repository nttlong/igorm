package orm_test

import (
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestTenantDB(t *testing.T) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=10000&encrypt=disable"
	db, err := orm.Open("mssql", sqlServerDns)

	if err != nil {
		t.Error(err)
		return
	}
	defer db.Close()
	tenantDb, err := orm.NewTenantDb(db)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, "aaa", tenantDb.DbName)
}
func BenchmarkTenantDB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=10000&encrypt=disable"
		db, err := orm.Open("mssql", sqlServerDns)

		if err != nil {
			b.Error(err)
			return
		}
		defer db.Close()
		tenantDb, err := orm.NewTenantDb(db)
		if err != nil {
			b.Error(err)
			return
		}
		b.Log(tenantDb)
	}

}
