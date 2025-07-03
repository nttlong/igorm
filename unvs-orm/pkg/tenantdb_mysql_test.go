package orm_test

import (
	"testing"
	orm "unvs-orm"
)

func TestTenantDBMySql(t *testing.T) {
	mysqlDns := "root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := orm.Open("mysql", mysqlDns)

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
func BenchmarkTenantDBMySql(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mysqlDns := "root:123456@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := orm.Open("mysql", mysqlDns)

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
