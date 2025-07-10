package pkgquerybuilder

import (
	"testing"
	"time"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestSelect(b *testing.T) {
	dialect := mssql() //<-- create mssql dialect
	//ctx := orm.JoinCompiler.Ctx(mssql()) //<-- create compiler context for mssql dialect
	repo := orm.Repository[OrderRepository]()
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
	assert.NoError(b, compilerResult.Err())
	sqlExpected := "SELECT [T1].[note] AS [note], [T1].[created_at] AS [created_at], [T1].[updated_at] AS [updated_at], [T1].[created_by] AS [created_by], [T2].[product] AS [product] FROM [orders] AS [T1] INNER JOIN [order_items] AS [T2] ON [T1].[order_id] = [T2].[order_id]"
	assert.Equal(b, sqlExpected, compilerResult.String())
	sqlText := compilerResult.String()

	assert.Equal(b, []interface{}(nil), compilerResult.Args)
	assert.Equal(b, sqlText, compilerResult.String())

}
func TestSelectLeftJoin(t *testing.T) {
	dialect := mssql() //<-- create mssql dialect
	//ctx := orm.JoinCompiler.Ctx(mssql()) //<-- create compiler context for mssql dialect
	repo := orm.Repository[OrderRepository]()
	join3Table := repo.Orders.OrderId.LeftJoin( //<-- join Order and OrderItem tables and select Order.Note, Order.CreatedAt, Order.UpdatedAt, Order.CreatedBy, OrderItem.Product
		repo.OrderItems.OrderId,
	).LeftJoin(repo.Customers.CustomerId)
	t.Log(join3Table)
	join2Table := repo.Orders.OrderId.LeftJoin( //<-- join Order and OrderItem tables and select Order.Note, Order.CreatedAt, Order.UpdatedAt, Order.CreatedBy, OrderItem.Product
		repo.OrderItems.OrderId)
	t.Log(join2Table)
	join3Table2 := repo.Orders.OrderId.LeftJoin( //<-- join Order and OrderItem tables and select Order.Note, Order.CreatedAt, Order.UpdatedAt, Order.CreatedBy, OrderItem.Product
		repo.OrderItems.OrderId,
	).RightJoin(repo.Customers.CustomerId)
	t.Log(join3Table2)
	joinOneTo2Tables := repo.Orders.OrderId.LeftJoin( //<-- join Order and OrderItem tables and select Order.Note, Order.CreatedAt, Order.UpdatedAt, Order.CreatedBy, OrderItem.Product
		repo.OrderItems.OrderId,
		repo.Customers.CustomerId,
	)
	sql := joinOneTo2Tables.Select(
		repo.Orders.Note,
		repo.Orders.CreatedAt,
		repo.Orders.UpdatedAt,
		repo.Orders.CreatedBy,
		repo.OrderItems.Product,
	)
	compilerResult := sql.Compile(dialect)
	assert.NoError(t, compilerResult.Err())
	sqlExpected := "SELECT [T1].[note] AS [note], [T1].[created_at] AS [created_at], [T1].[updated_at] AS [updated_at], [T1].[created_by] AS [created_by], [T2].[product] AS [product] FROM [orders] AS [T1] LEFT JOIN [order_items] AS [T2] ON [T1].[order_id] = [T2].[order_id]"
	assert.Equal(t, sqlExpected, compilerResult.String())
	sqlText := compilerResult.String()

	assert.Equal(t, []interface{}(nil), compilerResult.Args)
	assert.Equal(t, sqlText, compilerResult.String())

}
func TestSelectWhere(b *testing.T) {
	dialect := mssql() //<-- create mssql dialect
	//ctx := orm.JoinCompiler.Ctx(mssql()) //<-- create compiler context for mssql dialect
	repo := orm.Repository[OrderRepository]()
	when := time.Now()
	content := "test"
	sql := repo.Orders.OrderId.Eq( //<-- join Order and OrderItem tables and select Order.Note, Order.CreatedAt, Order.UpdatedAt, Order.CreatedBy, OrderItem.Product
		repo.OrderItems.OrderId,
	).Select(
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
	assert.NoError(b, compilerResult.Err())

	sqlExpected := "SELECT [T1].[note] AS [Note], [T1].[created_at] AS [CreatedAt], [T1].[updated_at] AS [UpdatedAt], [T1].[created_by] AS [CreatedBy], [T2].[product] AS [Product] FROM [orders] AS [T1]  JOIN [order_items] AS [T2] ON [T1].[order_id] = [T2].[order_id] WHERE [T1].[note] = ? AND [T1].[updated_at] = ?"
	assert.Equal(b, sqlExpected, compilerResult.String())
	sqlText := compilerResult.String()

	assert.Equal(b, []interface{}{content, when}, compilerResult.Args)
	assert.Equal(b, sqlText, compilerResult.String())

}
func TestSelectWhereWithJoinByStatement(b *testing.T) {
	dialect := mssql() //<-- create mssql dialect
	//ctx := orm.JoinCompiler.Ctx(mssql()) //<-- create compiler context for mssql dialect
	repo := orm.Repository[OrderRepository]()
	when := time.Now()
	content := "test"
	sql := repo.Join("customer.customerId=invoice.customerId AND customer.name=?", "John").Select(
		repo.Orders.Note,
		repo.Orders.CreatedAt,
		repo.Orders.UpdatedAt,
		repo.Orders.CreatedBy,
		repo.OrderItems.Product,
		repo.Expr("len(customer.name)+order.orderid"),
	).Where(
		repo.Orders.Note.Eq(content).And(
			repo.Orders.UpdatedAt.Eq(when),
		),
	)
	compilerResult := sql.Compile(dialect)
	assert.NoError(b, compilerResult.Err())

	sqlExpected := "SELECT [orders].[note] AS [Note], [orders].[created_at] AS [CreatedAt], [orders].[updated_at] AS [UpdatedAt], [orders].[created_by] AS [CreatedBy], [order_items].[product] AS [Product], len([T1].[name]) + [T2].[orderid] FROM [customers] AS [T1] INNER JOIN [invoices] AS [T2] ON [T1].[customer_id] = [T2].[customer_id] AND [T1].[name] = ? WHERE [orders].[note] = ? AND [orders].[updated_at] = ?"
	assert.Equal(b, sqlExpected, compilerResult.String())
	sqlText := compilerResult.String()

	assert.Equal(b, []interface{}{"John", "test", when}, compilerResult.Args)
	assert.Equal(b, sqlText, compilerResult.String())

}
func BenchmarkSelectWhereWithJoinByStatement(b *testing.B) {
	dialect := mssql() //<-- create mssql dialect
	//ctx := orm.JoinCompiler.Ctx(mssql()) //<-- create compiler context for mssql dialect
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < b.N; i++ {
		when := time.Now()
		content := "test"
		sql := repo.Join("customer.customerId=invoice.customerId AND customer.name=?", "John").Select(
			repo.Orders.Note,
			repo.Orders.CreatedAt,
			repo.Orders.UpdatedAt,
			repo.Orders.CreatedBy,
			repo.OrderItems.Product,
			repo.Expr("len(customer.name)+order.orderid"),
		).Where(
			repo.Orders.Note.Eq(content).And(
				repo.Orders.UpdatedAt.Eq(when),
			),
		)
		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err())

		sqlExpected := "SELECT [orders].[note] AS [Note], [orders].[created_at] AS [CreatedAt], [orders].[updated_at] AS [UpdatedAt], [orders].[created_by] AS [CreatedBy], [order_items].[product] AS [Product], len([T1].[name]) + [T2].[orderid] FROM [customers] AS [T1] INNER JOIN [invoices] AS [T2] ON [T1].[customer_id] = [T2].[customer_id] AND [T1].[name] = ? WHERE [orders].[note] = ? AND [orders].[updated_at] = ?"
		assert.Equal(b, sqlExpected, compilerResult.String())
		sqlText := compilerResult.String()

		assert.Equal(b, []interface{}{"John", "test", when}, compilerResult.Args)
		assert.Equal(b, sqlText, compilerResult.String())
	}

}
