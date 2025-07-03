package orm_test

import (
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestRepository_MSSQL(t *testing.T) {
	for i := 0; i < 10; i++ {
		mssqlDns := "server=localhost;database=master;user id=sa;password=123456;app name=test"
		db, err := orm.Open("mssql", mssqlDns)
		assert.NoError(t, err)
		defer db.Close()

		ret := orm.Repository[OrderRepository](db)
		assert.NoError(t, ret.Err)
		assert.Equal(t, ret.Orders.TenantDb, ret.TenantDb)
		assert.Equal(t, ret.OrderItems.TenantDb, ret.TenantDb)
	}

}
func BenchmarkRepository_MSSQL(b *testing.B) {
	mssqlDns := "server=localhost;database=master;user id=sa;password=123456;app name=test"
	db, err := orm.Open("mssql", mssqlDns)
	if err != nil {
		b.Error(err)
		return
	}
	defer db.Close()

	for i := 0; i < b.N; i++ {
		orm.Repository[OrderRepository](db)
	}
}
