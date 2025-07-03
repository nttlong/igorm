package orm_test

import (
	"testing"
	orm "unvs-orm"
)

func TestTenantDB(t *testing.T) {
	mssqlDns := "server=localhost;database=master;user id=sa;password=123456;app name=test"
	db, err := orm.Open("mssql", mssqlDns)

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
	t.Log(tenantDb)
}
func BenchmarkTenantDB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mssqlDns := "server=localhost;database=master;user id=sa;password=123456;app name=test"
		db, err := orm.Open("mssql", mssqlDns)

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
