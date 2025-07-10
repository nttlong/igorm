package main

import (
	"testing"

	"entgo.io/ent/dialect/sql"
	entSql "entgo.io/ent/dialect/sql"
	"github.com/stretchr/testify/assert"
	// alias để tránh trùng tên
)

func TestXxx(t *testing.T) {
	builder := sql.Dialect("postgres")

	// Tạo query: SELECT [users].[email] + [users].[user_id] FROM [users]
	q := builder.
		Select("[users].[email] + [users].[user_id]"). // string thô
		From(sql.Table("users"))

	q.Query()
}
func BenchmarkTest(b *testing.B) {
	builder := entSql.Dialect("mssql")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		q := builder.
			Select(
				"[users].[email] + ?",     // Biểu thức với tham số
				"[users].[user_name] + ?", // Biểu thức với tham số
			).
			From(entSql.Table("users"))

		sqlStr, _ := q.Query()

		assert.Equal(b,
			"SELECT `[users].[email] + ?`, `[users].[user_name] + ?` FROM `users`",
			sqlStr,
		)

	}
}

func BenchmarkEntSelectJoin(b *testing.B) {
	builder := entSql.Dialect("mssql")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		s := builder.
			Select("SUM(T1.OrderId) AS total_quantity"). // ✅ string
			From(entSql.Table("orders").As("T1")).
			LeftJoin(entSql.Table("order_items").As("T2")).
			On(
				"T1.OrderId ", "T2.OrderId",
			)

		sql, arg := s.Query() // <- Build SQL tại đây!
		assert.Equal(b, "SELECT SUM(OrderId) AS total_quantity FROM `orders`", sql)
		assert.Equal(b, []interface{}(nil), arg)
	}
}
