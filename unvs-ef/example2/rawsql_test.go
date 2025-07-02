package examples_test

import (
	"database/sql"
	"sync"
	"testing"
	"time"
)

type OrderData struct {
	OrderId   int64
	Version   int
	CreatedAt time.Time
	CreatedBy string
	Note      string
	UpdatedAt *time.Time
	UpdatedBy *string
}

var resultPool = sync.Pool{New: func() interface{} { return make([]OrderData, 10000) }}

func FetchOrdersDirectlyWithPool(rows *sql.Rows) ([]OrderData, error) {
	defer rows.Close()
	results := resultPool.Get().([]OrderData)
	defer resultPool.Put(results)
	var count int

	for rows.Next() {
		if count >= 10000 {
			break
		}
		err := rows.Scan(&results[count].OrderId, &results[count].CreatedAt, &results[count].CreatedBy, &results[count].Note, &results[count].UpdatedAt, &results[count].UpdatedBy, &results[count].Version)
		if err != nil {
			return nil, err
		}
		count++
	}

	return results[:count], nil
}
func BenchmarkRawSqlAndFetchAllRowsWithFetchOrdersDirectlyWithPoolAndSetCnnPool(b *testing.B) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=1000"
	db, err := sql.Open("sqlserver", sqlServerDns)
	if err != nil {
		b.Fatal(err)
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute)
	db.SetConnMaxIdleTime(time.Minute)

	defer db.Close()

	// Warmup
	sqlExec := `SELECT [order].[order_id] AS [OrderId], [order].[created_at] AS [CreatedAt], [order].[created_by] AS [CreatedBy], [order].[note] AS [Note], [order].[updated_at] AS [UpdatedAt], [order].[updated_by] AS [UpdatedBy], [order].[version] AS [Version] FROM [order] ORDER BY [order].[order_id] ASC, [order].[version] ASC OFFSET 0 ROWS FETCH NEXT 10000 ROWS ONLY`
	for i := 0; i < 5; i++ {
		rows, err := db.Query(sqlExec)
		if err != nil {
			b.Fatal(err)
		}
		rows.Close()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// start := time.Now()
		sqlExec := `SELECT [order].[order_id] AS [OrderId], [order].[created_at] AS [CreatedAt], [order].[created_by] AS [CreatedBy], [order].[note] AS [Note], [order].[updated_at] AS [UpdatedAt], [order].[updated_by] AS [UpdatedBy], [order].[version] AS [Version] FROM [order] ORDER BY [order].[order_id] ASC, [order].[version] ASC OFFSET 0 ROWS FETCH NEXT 10000 ROWS ONLY`
		rows, err := db.Query(sqlExec)
		if err != nil {
			b.Fatal(err)
		}
		_, err = FetchOrdersDirectlyWithPool(rows)
		// fmt.Printf("Tong thoi gian tu luc query den fetch toan bo 1000 dong la: \t\t%d \t\tnns\n", time.Since(start).Nanoseconds())
		if err != nil {
			b.Fatal(err)
		}
	}
}
