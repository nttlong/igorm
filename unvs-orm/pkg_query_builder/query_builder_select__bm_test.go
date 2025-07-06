package pkgquerybuilder

import (
	"testing"
	"time"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func BenchmarkTestSelect(b *testing.B) {
	dialect := mssql() //<-- create mssql dialect
	//ctx := orm.JoinCompiler.Ctx(mssql()) //<-- create compiler context for mssql dialect
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < b.N; i++ {
		sql := repo.Orders.OrderId.Eq( //<-- join Order and OrderItem tables and select Order.Note, Order.CreatedAt, Order.UpdatedAt, Order.CreatedBy, OrderItem.Product
			repo.OrderItems.OrderId,
		).Select(
			repo.Orders.Note,
			repo.Orders.CreatedAt,
			repo.Orders.UpdatedAt,
			repo.Orders.CreatedBy,
			repo.OrderItems.Product,
		)
		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err)
		sqlExpected := "SELECT [T1].[note] AS [Note], [T1].[created_at] AS [CreatedAt], [T1].[updated_at] AS [UpdatedAt], [T1].[created_by] AS [CreatedBy], [T2].[product] AS [Product] FROM [orders] AS [T1]  JOIN [order_items] AS [T2] ON [T1].[order_id] = [T2].[order_id]"
		assert.Equal(b, sqlExpected, compilerResult.SqlText)
		sqlText := compilerResult.SqlText

		assert.Equal(b, []interface{}(nil), compilerResult.Args)
		assert.Equal(b, sqlText, compilerResult.SqlText)
	}

}
func BenchmarkTestSelectWhere(b *testing.B) {
	dialect := mssql() //<-- create mssql dialect
	//ctx := orm.JoinCompiler.Ctx(mssql()) //<-- create compiler context for mssql dialect
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < b.N; i++ {
		when := time.Now()
		content := "test"
		sql := repo.Orders.CreatedAt.Eq( //<-- join Order and OrderItem tables and select Order.Note, Order.CreatedAt, Order.UpdatedAt, Order.CreatedBy, OrderItem.Product
			repo.OrderItems.CreatedAt,
		).Select(
			repo.Orders.OrderId.Text(),
			repo.Orders.Note,
			repo.Orders.CreatedAt,
			repo.Orders.UpdatedAt,
			repo.Orders.CreatedBy,
			repo.OrderItems.Product,
		).Where(
			repo.Orders.Note.Eq(content).And(
				repo.Orders.UpdatedAt.Eq(when),
			),
		)
		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err)

		sqlExpected := "SELECT CONVERT(NVARCHAR(50), [T1].[order_id]) AS [OrderId], [T1].[note] AS [Note], [T1].[created_at] AS [CreatedAt], [T1].[updated_at] AS [UpdatedAt], [T1].[created_by] AS [CreatedBy], [T2].[product] AS [Product] FROM [orders] AS [T1]  JOIN [order_items] AS [T2] ON [T1].[created_at] = [T2].[created_at] WHERE [T1].[note] = ? AND [T1].[updated_at] = ?"
		assert.Equal(b, sqlExpected, compilerResult.SqlText)
		sqlText := compilerResult.SqlText

		assert.Equal(b, []interface{}{content, when}, compilerResult.Args)
		assert.Equal(b, sqlText, compilerResult.SqlText)
	}

}
