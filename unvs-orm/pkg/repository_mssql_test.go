package orm_test

import (
	"testing"
	"time"
	"unsafe"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestRepository_MSSQL_NewEntityObject(t *testing.T) {

	repo := orm.Repository[OrderRepository]()
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
func BenchmarkRepository_NewEntityObject(b *testing.B) {

	repo := orm.Repository[OrderRepository]()
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
func TestRepository_MSSQL_NewEntityObjectThenInsert(t *testing.T) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=10000&encrypt=disable"
	db, err := orm.Open("mssql", sqlServerDns)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := orm.Repository[OrderRepository]()
	if repo.Err != nil {
		t.Fatal(repo.Err)
	}
	a := repo.Orders.New()
	a.Data.Note.Val = orm.Ptr("test note")
	a.Data.CreatedAt.Val = orm.Ptr(time.Now())
	a.Data.CreatedBy.Val = orm.Ptr("test user")
	a.Data.Version.Val = orm.Ptr(1)
	err = a.Insert()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, 0, a.Data.OrderId)
	assert.NotEqual(t, 0, a.Data.Version)
	assert.NotEqual(t, "", a.Data.Note)
}
func BenchmarkInsertRaw(b *testing.B) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=10000&encrypt=disable"
	db, err := orm.Open("mssql", sqlServerDns)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		b.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO orders (note, created_at, created_by, version) VALUES (?, ?, ?, ?)")
	if err != nil {
		b.Fatal(err)
	}
	defer stmt.Close()

	for j := 0; j < 10000; j++ {
		_, err := stmt.Exec("test note", time.Now(), "test user", 1)
		if err != nil {
			b.Fatal(err)
		}
	}
	tx.Commit()
}

func BenchmarkRepository_NewEntityObjectThenInsert(b *testing.B) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=10000&encrypt=disable"
	db, err := orm.Open("mssql", sqlServerDns)
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		b.Fatal(err)
	}

	repo := orm.Repository[OrderRepository]()
	if repo.Err != nil {
		b.Fatal(repo.Err)
	}

	a := repo.Orders.New()
	a.Use(tx)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()

		for k := 0; k < 10000; k++ {

			a.Data.Note.Val = orm.Ptr("test note")
			a.Data.CreatedAt.Val = orm.Ptr(time.Now())
			a.Data.CreatedBy.Val = orm.Ptr("test user")
			a.Data.Version.Val = orm.Ptr(1)
			err = a.InsertWithTransaction(tx)
			if err != nil {
				tx.Rollback()
				b.Fatal(err)
			}

		}
		err = tx.Commit()
		b.Log(err)

		elapsed := time.Since(start).Nanoseconds()
		avgPerOp := float64(elapsed) / 10000.0

		// Thêm metric tùy chỉnh
		b.ReportMetric(avgPerOp, "ns/op_for_new_10000_object")
	}
}

/*
goos: windows
goarch: amd64
pkg: unvs-orm/pkg
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkRepository_NewEntityObjectThenInsert-16    	       1	3864584300 ns/op	    386480 ns/op_for_new_10000_object	117343144 B/op	 1630538 allocs/op
PASS
ok  	unvs-orm/pkg	4.157s
*/
/*
goos: windows
goarch: amd64
pkg: unvs-orm/pkg
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
BenchmarkInsertRaw-16    	       1	1339304300 ns/op	63347864 B/op	 1000837 allocs/op
PASS
ok  	unvs-orm/pkg	1.423s
*/
//C:\Program Files\Go\bin\go.exe test -benchmem -run=^$ -bench ^BenchmarkRepository_NewEntityObject$ unvs-orm/pkg
