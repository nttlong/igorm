package orm_test

import (
	"testing"
	"time"
	"unsafe"
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
func TestRepository_MSSQL_New(t *testing.T) {
	mssqlDns := "server=localhost;database=master;user id=sa;password=123456;app name=test"
	db, err := orm.Open("mssql", mssqlDns)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := orm.Repository[OrderRepository](db)
	if repo.Err != nil {
		t.Fatal(repo.Err)
	}
	a := repo.Orders.New()

	b := repo.Orders.New()
	assert.NotEqual(t,
		uintptr(unsafe.Pointer(&a)),
		uintptr(unsafe.Pointer(&b)),
	)

	// Thêm metric tùy chỉnh

}
func BenchmarkRepository_MSSQL(b *testing.B) {
	mssqlDns := "server=localhost;database=master;user id=sa;password=123456;app name=test"
	db, err := orm.Open("mssql", mssqlDns)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	repo := orm.Repository[OrderRepository](db)
	if repo.Err != nil {
		b.Fatal(repo.Err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		start := time.Now()

		for k := 0; k < 10000; k++ {

			_ = repo.Orders.New() // Tạo mới object (nên tránh tính thời gian khác)
		}

		elapsed := time.Since(start).Nanoseconds()
		avgPerOp := float64(elapsed) / 10000.0

		// Thêm metric tùy chỉnh
		b.ReportMetric(avgPerOp, "ns/op_for_new_10000_object")
	}
}

/* unsafe
goos: windows
goarch: amd64
pkg: unvs-orm/pkg
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkRepository_MSSQL-16    	      91	  11077929 ns/op	      1394 ns/op_for_new_10000_object	12960637 B/op	   80009 allocs/op
PASS
ok  	unvs-orm/pkg	5.662s
*/

/*  safe
goos: windows
goarch: amd64
pkg: unvs-orm/pkg
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkRepository_MSSQL-16    	      91	  12228938 ns/op	      1875 ns/op_for_new_10000_object	12960361 B/op	   80009 allocs/op
PASS
ok  	unvs-orm/pkg	6.374s
*/
