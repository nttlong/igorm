package main

import (
	"testing"

	"entgo.io/ent/dialect/sql"
	entSql "entgo.io/ent/dialect/sql"
	"github.com/stretchr/testify/assert"
	// alias để tránh trùng tên
)

func TestSelfJoin(t *testing.T) {

}
func BenchmarkTestEntSelfJoin(b *testing.B) {

	builder := entSql.Dialect("mssql") // dùng MSSQL để tạo cú pháp có [brackets]
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q := builder.
			Select("*").
			From(entSql.Table("Order").As("Child")).
			Join(entSql.Table("Order").As("Parent")).
			OnP(entSql.ExprP("Child.ParentCode = Parent.Code"))

		sqlText, args := q.Query()
		assert.NoError(b, nil) // vì không có err trả về
		assert.Equal(b, "SELECT * FROM `Order` AS `Child` JOIN `Order` AS `Parent` ON Child.ParentCode = Parent.Code", sqlText)
		assert.Equal(b, 0, len(args))
	}
}
func BenchmarkTestEntSelectJoin(b *testing.B) {
	builder := entSql.Dialect("postgres")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q := builder.
			Select(
				"T1.code AS DepartmentCode",
				"T2.name AS UserName",
				"T3.name AS CheckName",
				"CONCAT(T2.first_name, ' ', T2.last_name) AS FullName",
			).
			From(entSql.Table("departments").As("T1")).
			Join(entSql.Table("users").As("T2")).
			On("T2.code", "T1.code").
			Join(entSql.Table("checks").As("T3")).
			On("T3.name", "?") // có thể binding 'John' nếu cần
		sqlText, _ := q.Query()
		assert.Equal(b, "SELECT T1.code AS DepartmentCode, T2.name AS UserName, T3.name AS CheckName, CONCAT(T2.first_name, ' ', T2.last_name) AS FullName FROM \"departments\" AS \"T1\" JOIN \"users\" AS \"T2\" ON \"T2.code\" = \"T1.code\" JOIN \"checks\" AS \"T3\" ON \"T3.name\" = \"?\"", sqlText)
	}
}

func TestEntSelectJoin(t *testing.T) {
	builder := sql.Dialect("postgres")
	q := builder.
		Select(
			"T1.code AS DepartmentCode",
			"T2.name AS UserName",
			"T3.name AS CheckName",
			"CONCAT(T2.first_name, ' ', T2.last_name) AS FullName",
		).
		From(entSql.Table("departments").As("T1")).
		Join(entSql.Table("users").As("T2")).
		On("T2.code", "T1.code").
		Join(entSql.Table("checks").As("T3")).
		On("T3.name", "?") // có thể binding 'John' nếu cần
	sqlText, _ := q.Query()
	assert.Equal(t, "SELECT T1.code AS DepartmentCode, T2.name AS UserName, T3.name AS CheckName, CONCAT(T2.first_name, ' ', T2.last_name) AS FullName FROM \"departments\" AS \"T1\" JOIN \"users\" AS \"T2\" ON \"T2.code\" = \"T1.code\" JOIN \"checks\" AS \"T3\" ON \"T3.name\" = \"?\"", sqlText)

}
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
