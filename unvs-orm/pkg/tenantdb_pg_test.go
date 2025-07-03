package orm_test

import (
	"testing"
	orm "unvs-orm"
)

func TestTenantDBPg(t *testing.T) {
	pgDns := "postgres://postgres:123456@localhost:5432/test?sslmode=disable"
	db, err := orm.Open("pgx", pgDns)

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
func BenchmarkTenantDBPg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pgDns := "postgres://postgres:123456@localhost:5432/test?sslmode=disable"
		db, err := orm.Open("postgres", pgDns)

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
