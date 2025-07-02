package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	ef "unvs.ef"
)

type Order struct {
	ef.Entity[Order]
	OrderId   ef.FieldNumber[uint64] `db:"primaryKey;autoIncrement"`
	Version   ef.FieldNumber[int]    `db:"primaryKey"`
	Note      ef.FieldString         `db:"length(200)"`
	CreatedAt ef.FieldDateTime
	UpdatedAt *ef.FieldDateTime
	CreatedBy ef.FieldString  `db:"length(100)"`
	UpdatedBy *ef.FieldString `db:"length(100)"`
}
type OrderItem struct {
	*ef.Entity[OrderItem]
	Id        ef.FieldNumber[uint64] `db:"primaryKey;autoIncrement"`
	OrderId   ef.FieldNumber[uint64] `db:"index(order_ref_idx)"`
	Version   ef.FieldNumber[int]    `db:"index(order_ref_idx)"`
	Product   ef.FieldString         `db:"length(100)"`
	Quantity  ef.FieldNumber[int]
	CreatedAt ef.FieldDateTime
	UpdatedAt *ef.FieldDateTime
	CreatedBy ef.FieldString  `db:"length(100)"`
	UpdatedBy *ef.FieldString `db:"length(100)"`
}
type OrderRepository struct {
	*ef.TenantDb
	Orders     *Order
	OrderItems *OrderItem
}

func (r *OrderRepository) Init() {
	r.NewRelationship().
		From(r.Orders.OrderId, r.Orders.Version).
		To(r.OrderItems.OrderId, r.OrderItems.Version)
}
func Test1(t *testing.T) {
	// db := &ef.DbField{}
	// v := reflect.ValueOf(db).Elem()
	// imp := v.FieldByName("imp")
	// if imp.IsNil() {
	// 	imp = reflect.New(imp.Type().Elem()).Elem()
	// }
	// fmt.Println(imp.IsValid())
	// imp.FieldByName("Key").SetString("orders")
	// fmt.Println(db.GetKey())
}
func TestDynaType(t *testing.T) {
	type FieldValue interface {
		int | int64 | uint64 | bool | float64
	}

	type NumField[T FieldValue] struct {
		Number T
	}
	fx := NumField[int64]{}
	val := reflect.ValueOf(&fx)
	fmt.Println(val.Type().String())
}
func TestDML(t *testing.T) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", sqlServerDns)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := ef.Repo[OrderRepository](db, true) // create repos
	if repo.Err != nil {
		fmt.Print(repo.Err)
	}
	assert.NoError(t, repo.Err)
	obj := repo.OrderItems.New()
	obj.OrderId.Set(ef.Lit(uint64(1)))
	obj.Version.Set(ef.Lit(0))

	obj.Product.Set(ef.Lit("product1"))
	obj.Quantity.Set(ef.Lit(0))
	obj.CreatedAt.Set(nil)

}
func BenchmarkDML(b *testing.B) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", sqlServerDns)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := ef.Repo[OrderRepository](db, true) // create repos
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			assert.NoError(b, repo.Err)
			obj := repo.OrderItems.New()       //<-- create new object
			obj.OrderId.Set(ef.Lit(uint64(1))) //<-- set valaues
			obj.Version.Set(ef.Lit(0))

			obj.Product.Set(ef.Lit("product1"))
			obj.Quantity.Set(ef.Lit(0))
			obj.CreatedAt.Set(nil)
		}
	}

}
func TestCompilerJoinExpr(t *testing.T) {
	//repo.Orders.OrderId.Eq(repo.OrderItems.OrderId)
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", sqlServerDns)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := ef.Repo[OrderRepository](db, false) // create repos
	expr := repo.Orders.OrderId.Eq(repo.OrderItems.OrderId)
	ret := expr.InnerJoin()
	fmt.Println(ret)
	ret = repo.Orders.Note.Len().Add(10).Eq(repo.OrderItems.Quantity.Add(100)).InnerJoin().InnerJoin(
		// repo.Orders.OrderId.Eq(repo.OrderItems.OrderId),
		// repo.Orders.Version.Le(repo.OrderItems.Version),
		repo.OrderItems.Quantity.Eq(1000),
	)
	sql, args := ret.GetSqlJoin(repo.Dialect)
	fmt.Println(sql)
	fmt.Println(args)
}
func TestCreateRepository(t *testing.T) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", sqlServerDns)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := ef.Repo[OrderRepository](db, true) // create repos
	if repo.Err != nil {
		fmt.Print(repo.Err)
	}
	assert.NoError(t, repo.Err)
	cmd := repo.From(repo.Orders).Select(
		repo.Orders.OrderId,
		repo.Orders.Version).OrderBy(
		repo.Orders.OrderId.Desc(),
		repo.Orders.Version.Desc(),
	).Limit(100)

	sql, args := cmd.ToSQL(repo.Dialect)
	fmt.Println(sql)
	fmt.Println(args)
	// rows, err := cmd.ExecTo(struct {
	// 	OrderId int64
	// 	Version int
	// }{})
	// if err != nil {
	// 	t.Fatal(err)
	// }

	r, err := ef.GetRows[struct {
		OrderId int64
		Version int
	}](cmd)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(r))
	fmt.Println(len(r))

}

func FetchOrdersDirectly(rows *sql.Rows) ([]OrderData, error) {
	defer rows.Close()
	results := resultPool.Get().([]OrderData)
	defer resultPool.Put(results[:0])

	for rows.Next() {
		var order OrderData
		var updatedAt sql.NullTime
		var updatedBy sql.NullString

		err := rows.Scan(&order.OrderId, &order.CreatedAt, &order.CreatedBy, &order.Note, &updatedAt, &updatedBy, &order.Version)
		if err != nil {
			return nil, err
		}

		if updatedAt.Valid {
			order.UpdatedAt = &updatedAt.Time
		}
		if updatedBy.Valid {
			order.UpdatedBy = &updatedBy.String
		}

		results = append(results, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

var resultPool = sync.Pool{New: func() interface{} { return make([]OrderData, 0, 10000) }}

// FetchOrdersDirectly đọc dữ liệu từ *sql.Rows và điền vào slice []OrderData
func FetchOrdersDirectlyWithPool(rows *sql.Rows) ([]OrderData, error) {
	defer rows.Close() // Đảm bảo rows luôn được đóng

	var orders []OrderData
	for rows.Next() {
		var order OrderData
		var updatedAt sql.NullTime   // Biến tạm cho time.Time nullable
		var updatedBy sql.NullString // Biến tạm cho string nullable

		err := rows.Scan(
			&order.OrderId,
			&order.CreatedAt,
			&order.CreatedBy,
			&order.Note,
			&updatedAt, // Scan vào biến tạm nullable
			&updatedBy, // Scan vào biến tạm nullable
			&order.Version,
		)
		if err != nil {
			return nil, err
		}

		// Xử lý các trường nullable
		if updatedAt.Valid {
			order.UpdatedAt = &updatedAt.Time
		} else {
			order.UpdatedAt = nil
		}
		if updatedBy.Valid {
			order.UpdatedBy = &updatedBy.String
		} else {
			order.UpdatedBy = nil
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
func BenchmarkRawSqlAndFetchAllRowsWithFetchOrdersDirectlyWithPoolAndSetCnnPool(b *testing.B) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", sqlServerDns)
	if err != nil {
		b.Fatal(err)
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute)
	defer db.Close()
	b.ResetTimer()
	avg := int64(0)
	for i := 0; i < b.N; i++ {
		start := time.Now()
		sqlExec := `SELECT [order].[order_id] AS [OrderId], [order].[created_at] AS [CreatedAt], [order].[created_by] AS [CreatedBy], [order].[note] AS [Note], [order].[updated_at] AS [UpdatedAt], [order].[updated_by] AS [UpdatedBy], [order].[version] AS [Version] FROM [order] ORDER BY [order].[order_id] ASC, [order].[version] ASC OFFSET 0 ROWS FETCH NEXT 10000 ROWS ONLY`
		rows, err := db.Query(sqlExec)
		if err != nil {
			b.Fatal(err)
		}
		start = time.Now()
		_, err = FetchOrdersDirectlyWithPool(rows)
		n := time.Since(start).Nanoseconds()
		avg += n
		// fmt.Printf("Tong thoi gian tu luc query den fetch toan bo 1000 dong la: \t\t%d \t\tnns\n", time.Since(start).Nanoseconds())
		if err != nil {
			b.Fatal(err)
		}
	}
	b.Log("Avg time: ", avg/int64(b.N))
}
func BenchmarkRawSqlAndFetchAllRowsWithFetchOrdersDirectlyWithSetCnnPool(b *testing.B) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", sqlServerDns)

	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(100) //<-- da dat SetMaxOpenConns
	db.SetMaxIdleConns(10)  //<-- da dat SetMaxOpenConns
	defer db.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()
		sqlExec := `SELECT [order].[order_id] AS [OrderId], [order].[created_at] AS [CreatedAt], [order].[created_by] AS [CreatedBy], [order].[note] AS [Note], [order].[updated_at] AS [UpdatedAt], [order].[updated_by] AS [UpdatedBy], [order].[version] AS [Version] FROM [order] ORDER BY [order].[order_id] ASC, [order].[version] ASC OFFSET 0 ROWS FETCH NEXT 10000 ROWS ONLY`
		rows, err := db.Query(sqlExec) //<-- exec chu fetch du lieu
		if err != nil {
			b.Fatal(err)
		}

		//ef.FetchAllRows(rows, reflect.TypeOf(OrderData{})) //<-- lan nay co doc du lieu
		_, err = FetchOrdersDirectly(rows)
		fmt.Printf("Tong thoi gian tu luc query den fetch toan bo 1000 dong la: \t\t%d \t\tnns\n", time.Since(start).Nanoseconds())
		if err != nil {
			b.Fatal(err)
		}

		// fmt.Println(cmd)
	}

}
func BenchmarkRawSqlAndFetchAllRowsWithFetchOrdersDirectly(b *testing.B) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", sqlServerDns)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		start := time.Now()
		sqlExec := `SELECT [order].[order_id] AS [OrderId], [order].[created_at] AS [CreatedAt], [order].[created_by] AS [CreatedBy], [order].[note] AS [Note], [order].[updated_at] AS [UpdatedAt], [order].[updated_by] AS [UpdatedBy], [order].[version] AS [Version] FROM [order] ORDER BY [order].[order_id] ASC, [order].[version] ASC OFFSET 0 ROWS FETCH NEXT 10000 ROWS ONLY`
		rows, err := db.Query(sqlExec) //<-- exec chu fetch du lieu
		if err != nil {
			b.Fatal(err)
		}

		//ef.FetchAllRows(rows, reflect.TypeOf(OrderData{})) //<-- lan nay co doc du lieu
		_, err = FetchOrdersDirectly(rows)
		fmt.Printf("Tong thoi gian tu luc query den fetch toan bo 1000 dong la: \t\t%d \t\tnns\n", time.Since(start).Nanoseconds())
		if err != nil {
			b.Fatal(err)
		}

		// fmt.Println(cmd)
	}

}

func TestInsert(t *testing.T) {
	type OrderData struct {
		OrderId   int64
		Version   int
		Note      string
		CreatedAt time.Time
		CreatedBy string
	}
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", sqlServerDns)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := ef.Repo[OrderRepository](db, true)
	avg := int64(0)
	for i := 0; i < 1000; i++ {

		sqlCmd := repo.InsertInto(repo.Orders).Values(&OrderData{
			CreatedBy: "admin",
			CreatedAt: time.Now(),

			Version: i,
			Note:    "test",
		})
		start := time.Now()

		result, err := sqlCmd.Exec()
		n := time.Since(start).Milliseconds()
		avg += n

		fmt.Printf("exec time: \t\t%d \t\tnns\n", n)
		fmt.Println("-------------------------")
		if err != nil {
			fmt.Print(err)
			t.Fatal(err)
		}
		t.Log(result)
	}

	fmt.Printf("arg avg time: %d nns\n", avg/1000)

}
func BenchmarkInsert(b *testing.B) {
	type OrderData struct {
		OrderId   int64
		Version   int
		Note      string
		CreatedAt time.Time
		CreatedBy string
	}
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa"
	db, err := sql.Open("sqlserver", sqlServerDns)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := ef.Repo[OrderRepository](db, true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		sqlCmd := repo.InsertInto(repo.Orders).Values(&OrderData{
			CreatedBy: "admin",
			CreatedAt: time.Now(),

			Version: i,
			Note:    "test",
		})
		// sqlCmd.ToSQL()
		result, err := sqlCmd.ExecWithContext(b.Context())

		if err != nil {
			fmt.Print(err)

		}
		_ = result

	}
}
