package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	//github.com/google/uuid v1.6.0
	_ "github.com/microsoft/go-mssqldb"

	ef "unvs.ef"
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

func BenchmarkSelectWithWhereOrderByLimit(b *testing.B) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=10000&encrypt=disable"
	db, err := sql.Open("sqlserver", sqlServerDns)
	db.SetMaxOpenConns(75)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(time.Minute * 2)
	db.SetConnMaxIdleTime(time.Minute)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := ef.Repo[OrderRepository](db, true) // create repos

	// sql, arsg := cmd.ToSQL(repo.Dialect)

	argTime := &struct {
		compilerTime  int64
		execSqlTime   int64
		fetchDataTime int64
	}{
		compilerTime:  0,
		execSqlTime:   0,
		fetchDataTime: 0,
	}
	cmdSelect := repo.From(repo.Orders).Select(
		repo.Orders.OrderId,
		repo.Orders.CreatedAt,
		repo.Orders.CreatedBy,
		repo.Orders.Note,
		repo.Orders.UpdatedAt,
		repo.Orders.UpdatedBy,
		repo.Orders.Version,
	)
	cmdSelect = cmdSelect.Where(repo.Orders.OrderId.Gt(1000))
	cmd := cmdSelect.OrderBy(
		repo.Orders.OrderId.Desc(),
		repo.Orders.Version.Desc(),
	)
	cmd = cmd.Limit(10000)
	sql, args := cmd.ToSQL(repo.Dialect)
	for i := 0; i < 10; i++ {
		db.QueryContext(b.Context(), sql, args...)
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatalf("Prepare error: %v", err)
	}
	defer stmt.Close()
	b.ResetTimer()

	fmt.Println("\n---------------------------------------------------")
	for i := 0; i < b.N; i++ {

		//start := time.Now()

		// n := time.Since(start).Nanoseconds()
		// fmt.Printf("Exec sql: %d\n", n)
		// argTime.compilerTime += n
		// b.Log(fmt.Sprintf("Build sql time: \t\t%d \t\tnns\n", n))
		start := time.Now()
		// sql = "SELECT [order].[order_id] AS [OrderId], [order].[created_at] AS [CreatedAt], [order].[created_by] AS [CreatedBy], [order].[note] AS [Note], [order].[updated_at] AS [UpdatedAt], [order].[updated_by] AS [UpdatedBy], [order].[version] AS [Version] FROM [order] WHERE ([order].[order_id] > @p1) ORDER BY [order].[order_id] ASC, [order].[version] ASC OFFSET 0 ROWS FETCH NEXT 10000 ROWS ONLY"
		rows, err := stmt.QueryContext(b.Context(), args...)
		if err != nil {
			b.Fatal(err)
		}

		n := time.Since(start).Nanoseconds()
		fmt.Printf("Exec sql: %d ns\n", n)
		argTime.execSqlTime += n
		// b.Log(fmt.Sprintf("Call db.Query(sql, args...): \t\t%d \t\tnns\n", n))
		start = time.Now()
		_, err = ef.FetchAllRows(rows, reflect.TypeOf(OrderData{}))
		// _, err = ef.FetchAllRows(rows, reflect.TypeOf(OrderData{}))
		n = time.Since(start).Nanoseconds()
		fmt.Printf("ef.FetchAllRows: %d ns\n", n)
		fmt.Println("---------------------------------------------------")
		argTime.fetchDataTime += n

		if err != nil {
			b.Fatal(err)
		}
	}
	argTime.compilerTime = argTime.compilerTime / int64(b.N)
	argTime.execSqlTime = argTime.execSqlTime / int64(b.N)
	argTime.fetchDataTime = argTime.fetchDataTime / int64(b.N)
	b.Log(fmt.Sprintf("compilerTime: \t\t%d \t\tnns\n", argTime.compilerTime))
	b.Log(fmt.Sprintf("execSqlTime: \t\t%d \t\tnns\n", argTime.execSqlTime))
	b.Log(fmt.Sprintf("fetchDataTime: \t\t%d \t\tnns\n", argTime.fetchDataTime))
}
func BenchmarkSelectWithWhereOrderByLimitInlineSql(b *testing.B) {
	sqlServerDns := "sqlserver://sa:123456@localhost?database=aaa&fetchSize=10000&encrypt=disable"
	db, err := sql.Open("sqlserver", sqlServerDns)
	// db.SetMaxOpenConns(75)
	// db.SetMaxIdleConns(3)
	// db.SetConnMaxLifetime(time.Minute * 2)
	// db.SetConnMaxIdleTime(time.Minute)

	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := ef.Repo[OrderRepository](db, true) // create repos

	// sql, arsg := cmd.ToSQL(repo.Dialect)

	argTime := &struct {
		compilerTime  int64
		execSqlTime   int64
		fetchDataTime int64
	}{
		compilerTime:  0,
		execSqlTime:   0,
		fetchDataTime: 0,
	}

	b.ResetTimer()

	fmt.Println("\n---------------------------------------------------")
	for i := 0; i < b.N; i++ {

		cmdSelect := repo.From(repo.Orders).Select(
			repo.Orders.OrderId,
			repo.Orders.CreatedAt,
			repo.Orders.CreatedBy,
			repo.Orders.Note,
			repo.Orders.UpdatedAt,
			repo.Orders.UpdatedBy,
			repo.Orders.Version,
		)
		cmdSelect = cmdSelect.Where(repo.Orders.OrderId.Gt(1000).And(*repo.Orders.Version.Lt(1000)))
		cmd := cmdSelect.OrderBy(
			repo.Orders.OrderId.Desc(),
			repo.Orders.Version.Desc(),
		)
		cmd = cmd.Limit(10000)
		sql, args := cmd.ToSQL(repo.Dialect)
		start := time.Now()
		// sql = "SELECT [order].[order_id] AS [OrderId], [order].[created_at] AS [CreatedAt], [order].[created_by] AS [CreatedBy], [order].[note] AS [Note], [order].[updated_at] AS [UpdatedAt], [order].[updated_by] AS [UpdatedBy], [order].[version] AS [Version] FROM [order] WHERE ([order].[order_id] > @p1) ORDER BY [order].[order_id] ASC, [order].[version] ASC OFFSET 0 ROWS FETCH NEXT 10000 ROWS ONLY"
		rows, err := db.QueryContext(b.Context(), sql, args...)
		if err != nil {
			b.Fatal(err)
		}

		n := time.Since(start).Nanoseconds()
		fmt.Printf("Exec sql: %d ns\n", n)
		argTime.execSqlTime += n
		// b.Log(fmt.Sprintf("Call db.Query(sql, args...): \t\t%d \t\tnns\n", n))
		start = time.Now()
		_, err = ef.FetchAllRows(rows, reflect.TypeOf(OrderData{}))
		// _, err = ef.FetchAllRows(rows, reflect.TypeOf(OrderData{}))
		n = time.Since(start).Nanoseconds()
		fmt.Printf("ef.FetchAllRows: %d ns\n", n)
		fmt.Println("---------------------------------------------------")
		argTime.fetchDataTime += n

		if err != nil {
			b.Fatal(err)
		}
	}
	argTime.compilerTime = argTime.compilerTime / int64(b.N)
	argTime.execSqlTime = argTime.execSqlTime / int64(b.N)
	argTime.fetchDataTime = argTime.fetchDataTime / int64(b.N)
	b.Log(fmt.Sprintf("compilerTime: \t\t%d \t\tnns\n", argTime.compilerTime))
	b.Log(fmt.Sprintf("execSqlTime: \t\t%d \t\tnns\n", argTime.execSqlTime))
	b.Log(fmt.Sprintf("fetchDataTime: \t\t%d \t\tnns\n", argTime.fetchDataTime))
}
