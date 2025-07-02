package main

import (
	"database/sql"
	"log"
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

var db *sql.DB

func init() {
	var err error
	// Cập nhật DSN đúng với môi trường của bạn
	dsn := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=1000"
	db, err = sql.Open("sqlserver", dsn)
	if err != nil {
		log.Fatalf("Open DB error: %v", err)
	}

	// Optional: kiểm tra kết nối
	if err := db.Ping(); err != nil {
		log.Fatalf("Ping DB error: %v", err)
	}
}

func BenchmarkSelectOptimized(b *testing.B) {
	b.ReportAllocs()

	stmt, err := db.Prepare(`
SELECT 
	[order].[order_id], 
	[order].[created_at], 
	[order].[created_by], 
	[order].[note], 
	[order].[updated_at], 
	[order].[updated_by], 
	[order].[version] 
FROM [order] 
ORDER BY [order].[order_id], [order].[version] 
OFFSET 0 ROWS FETCH NEXT 10000 ROWS ONLY`)
	if err != nil {
		b.Fatalf("prepare failed: %v", err)
	}
	defer stmt.Close()

	orders := make([]OrderData, 0, 10000)

	for i := 0; i < b.N; i++ {
		orders = orders[:0] // reuse slice
		rows, err := stmt.Query()
		if err != nil {
			b.Fatalf("query failed: %v", err)
		}
		for rows.Next() {
			var o OrderData
			err := rows.Scan(
				&o.OrderId,
				&o.CreatedAt,
				&o.CreatedBy,
				&o.Note,
				&o.UpdatedAt,
				&o.UpdatedBy,
				&o.Version,
			)
			if err != nil {
				b.Fatalf("scan failed: %v", err)
			}
			orders = append(orders, o)
		}
		rows.Close()
	}
}
