package main

import (
	ent "check/ent"
	_ "check/ent/schema"

	"log"
	"os"
	"runtime/pprof"

	// "your_project/ent"

	entSql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq" // hoặc driver khác nếu dùng PostgreSQL/MySQL/SQLite
)

func main() {
	// Init Ent client (dù không exec query, vẫn cần)
	client, err := ent.Open("postgres", "host=localhost port=5432 user=postgres dbname=testdb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Tạo CPU profile
	f, _ := os.Create("cpu22-ent.prof")
	//go tool pprof -http=localhost:8080 ./cpu19-ent.prof
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	// ctx := context.Background()

	// Loop để benchmark
	for i := 0; i < 1_000_000; i++ {
		builder := entSql.Dialect("postgres")

		// Tạo query: SELECT [users].[email] + [users].[user_id] FROM [users]
		q := builder.
			Select("*"). // hoặc chọn từng field nếu muốn
			From(entSql.Table("departments").As("T1")).
			Join(entSql.Table("users").As("T2")).
			On("T2.code", "T1.code").
			Join(entSql.Table("checks").As("T3")).
			On("T3.name", "?") // binding parameter

		q.Query()

	}

	pprof.StopCPUProfile()
}
