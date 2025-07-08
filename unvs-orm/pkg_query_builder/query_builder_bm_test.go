package pkgquerybuilder

import (
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestCompileExprToSqlSyntax(b *testing.T) {
	repo := orm.Repository[OrderRepository]()
	ctx := orm.Compiler.Ctx(mssql()) //<-- chuan bi bo compiler, la duy nhat cho moi database driver, truong hop nay dang dung voi mssql
	w := repo.Orders.OrderId.Eq(1).And(
		repo.Orders.Note.Eq("test"),
	)

	stm, err := ctx.ResolveWithoutTableAlias(w)
	assert.NoError(b, err)
	assert.Equal(b, "[orders].[order_id] = ? AND [orders].[note] = ?", stm.Syntax)
	assert.Equal(b, []interface{}{1, "test"}, stm.Args)
	stm2, err := ctx.ResolveWithTableAlias(w)
	assert.NoError(b, err)
	assert.Equal(b, "[T1].[order_id] = ? AND [T1].[note] = ?", stm2.Syntax)
	assert.Equal(b, []interface{}{1, "test"}, stm.Args)
	assert.Equal(b, []string{"orders"}, stm2.Tables)
	assert.Equal(b, map[string]string{"orders": "T1"}, stm2.GetTableAliasMap())

}
func TestSimpleQuery(t *testing.T) {
	repo := orm.Repository[OrderRepository]()
	w := repo.Orders.OrderId.Eq(1).And(
		repo.Orders.Note.Eq("test"),
	)
	orderQuery := repo.Orders.Filter(
		w,
	).Select(repo.Orders.OrderId.Max().As("max_order_id"),
		repo.Orders.Note,
	).GroupBy(repo.Orders.Note).Having(
		repo.Orders.Note.Eq("test"),
	)
	sql, err := orderQuery.ToSql(mssql())
	assert.NoError(t, err)
	expectedSql := "SELECT MAX([T1].[order_id]) AS [max_order_id], [T1].[note] FROM [orders] AS [T1] WHERE [T1].[order_id] = ? AND [T1].[note] = ? GROUP BY [T1].[note] HAVING [T1].[note] = ?"
	assert.Equal(t, expectedSql, sql.Sql)
}
func BenchmarkCompileExprWithoutTableAliasToSqlSyntax(b *testing.B) {
	repo := orm.Repository[OrderRepository]()
	ctx := orm.Compiler.Ctx(mssql()) //<-- chuan bi bo compiler, la duy nhat cho moi database driver, truong hop nay dang dung voi mssql
	for i := 0; i < b.N; i++ {
		w := repo.Orders.OrderId.Eq(1).And(
			repo.Orders.Note.Eq("test"),
		)

		stm, err := ctx.ResolveWithoutTableAlias(w)
		assert.NoError(b, err)
		assert.Equal(b, "[orders].[order_id] = ? AND [orders].[note] = ?", stm.Syntax)
		assert.Equal(b, []interface{}{1, "test"}, stm.Args)
	}

}
func BenchmarkCompileExprWithTableAliasToSqlSyntax(b *testing.B) {
	repo := orm.Repository[OrderRepository]()
	ctx := orm.Compiler.Ctx(mssql()) //<-- chuan bi bo compiler, la duy nhat cho moi database driver, truong hop nay dang dung voi mssql
	for i := 0; i < b.N; i++ {
		w := repo.Orders.OrderId.Eq(1).And(
			repo.Orders.Note.Eq("test"),
		)

		stm, err := ctx.ResolveWithTableAlias(w)
		assert.NoError(b, err)
		assert.Equal(b, "[T1].[order_id] = ? AND [T1].[note] = ?", stm.Syntax)
		assert.Equal(b, []interface{}{1, "test"}, stm.Args)
		assert.Equal(b, []string{"orders"}, stm.Tables)
		assert.Equal(b, map[string]string{"orders": "T1"}, stm.GetTableAliasMap())
	}

}
func BenchmarkQueryBuilder(b *testing.B) {
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < b.N; i++ {

		w := repo.Orders.OrderId.Eq(1).And(
			repo.Orders.Note.Eq("test"),
		)
		orderQuery := repo.Orders.Filter(
			w,
		).Select(repo.Orders.OrderId.Max().As("max_order_id"),
			repo.Orders.Note,
		).GroupBy(repo.Orders.Note).Having(
			repo.Orders.Note.Eq("test"),
		)
		b.Log(orderQuery)
		sql, err := orderQuery.ToSql(mssql())
		assert.NoError(b, err)
		expectedSql := "SELECT MAX([T1].[order_id]) AS [max_order_id], [T1].[note] FROM [orders] AS[T1] WHERE [T1].[order_id] = ? AND [T1].[note] = ? GROUP BY [T1].[note] HAVING [T1].[note] = ?"
		assert.Equal(b, expectedSql, sql.Sql)
	}

}
func BenchmarkTestJoinExpr(b *testing.B) {
	repo := orm.Repository[OrderRepository]()
	b.Log(repo)
	for i := 0; i < b.N; i++ {

		on := repo.Orders.OrderId.Eq(repo.OrderItems.OrderId)
		join := repo.Orders.Join(repo.OrderItems, on)
		ctx := orm.JoinCompiler.Ctx(mssql())
		joinRes, err := ctx.Resolve(join)
		assert.NoError(b, err)
		expectedSql := "[orders] AS [T1] INNER JOIN [order_items] AS [T2] ON [T1].[order_id] = [T2].[order_id]"
		assert.Equal(b, expectedSql, joinRes.Syntax)
	}

}
func BenchmarkTestJoinExpr2(b *testing.B) {
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < b.N; i++ {

		on := repo.Orders.OrderId.Add(1).Eq(repo.OrderItems.OrderId)
		join := repo.Orders.Join(repo.OrderItems, on)
		ctx := orm.JoinCompiler.Ctx(mssql())
		joinRes, err := ctx.Resolve(join)
		assert.NoError(b, err)
		expectedSql := "[orders] AS [T1] INNER JOIN [order_items] AS [T2] ON [T1].[order_id] + ? = [T2].[order_id]"
		assert.Equal(b, expectedSql, joinRes.Syntax)
		assert.Equal(b, []interface{}{1}, joinRes.Args)
		on2 := repo.Orders.Note.Len().Eq(repo.OrderItems.OrderId)
		expectedSql2 := "[orders] AS [T1] INNER JOIN [order_items] AS [T2] ON LEN([T1].[note]) = [T2].[order_id]"
		join2 := repo.Orders.Join(repo.OrderItems, on2)
		joinRes2, err := ctx.Resolve(join2)
		assert.NoError(b, err)
		assert.Equal(b, expectedSql2, joinRes2.Syntax)
		on3 := repo.Orders.Note.Len().Eq(repo.OrderItems.OrderId.Add(1))
		expectedSql3 := "[orders] AS [T1] INNER JOIN [order_items] AS [T2] ON LEN([T1].[note]) = [T2].[order_id] + ?"
		join3 := repo.Orders.Join(repo.OrderItems, on3)
		joinRes3, err := ctx.Resolve(join3)
		assert.NoError(b, err)
		assert.Equal(b, expectedSql3, joinRes3.Syntax)
		assert.Equal(b, []interface{}{1}, joinRes3.Args)
	}
}
func BenchmarkTestJoin3Tables(b *testing.B) {
	ctx := orm.JoinCompiler.Ctx(mssql())
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < b.N; i++ {

		join2 := repo.Invoices.Join(
			repo.InvoiceDetails, repo.Invoices.InvoiceId.Eq(repo.InvoiceDetails.InvoiceId),
		).Join(
			repo.Items, repo.InvoiceDetails.ItemId.Eq(repo.Items.ItemId),
		)
		joinRes2, err := ctx.Resolve(join2)
		assert.NoError(b, err)
		//"[orders] AS [T3] INNER JOIN [invoice_details] AS [T2] ON [T1].[invoice_id] = [T2].[invoice_id] INNER JOIN [] AS [] ON [T3].[order_id] = [T1].[order_id]"
		//"[orders] AS [T3] INNER JOIN [invoice_details] AS [T2] ON [T1].[invoice_id] = [T2].[invoice_id] INNER JOIN [] AS [] ON [T3].[order_id] = [T1].[order_id]"
		expectedSql2 := "[invoices] AS [T1] INNER JOIN [invoice_details] AS [T2] ON [T1].[invoice_id] = [T2].[invoice_id] INNER JOIN [items] AS [T3] ON [T2].[item_id] = [T3].[item_id]"
		assert.Equal(b, expectedSql2, joinRes2.Syntax)
	}
}
func BenchmarkTestJoin3Tables2(b *testing.B) {
	ctx := orm.JoinCompiler.Ctx(mssql())
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < b.N; i++ {
		join2 := repo.Invoices.Join(
			repo.InvoiceDetails, repo.Invoices.InvoiceId.Eq(repo.InvoiceDetails.InvoiceId),
		).Join(
			repo.Customers, repo.Invoices.CustomerId.Eq(repo.Customers.CustomerId),
		).Join(
			repo.PaymentMethods, repo.Invoices.PaymentMethodId.Eq(repo.PaymentMethods.PaymentMethodId),
		).Join(
			repo.Items, repo.InvoiceDetails.ItemId.Eq(repo.Items.ItemId),
		)
		joinRes2, err := ctx.Resolve(join2)
		assert.NoError(b, err)

		expectedSql2 := "[invoices] AS [T1] INNER JOIN [invoice_details] AS [T2] ON [T1].[invoice_id] = [T2].[invoice_id] INNER JOIN [customers] AS [T3] ON [T1].[customer_id] = [T3].[customer_id] INNER JOIN [payment_methods] AS [T4] ON [T1].[payment_method_id] = [T4].[payment_method_id] INNER JOIN [items] AS [T5] ON [T2].[item_id] = [T5].[item_id]"
		assert.Equal(b, expectedSql2, joinRes2.Syntax)
	}
}
func BenchmarkTestLeftJoinExpr(b *testing.B) {
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < b.N; i++ {
		on := repo.Orders.OrderId.Eq(repo.OrderItems.OrderId)
		join := repo.Orders.LeftJoin(repo.OrderItems, on).LeftJoin(
			repo.Invoices, repo.Invoices.OrderId.Eq(repo.OrderItems.OrderId),
		)
		ctx := orm.JoinCompiler.Ctx(mssql())
		joinRes, err := ctx.Resolve(join)
		assert.NoError(b, err)
		expectedSql := "[orders] AS [T1] LEFT JOIN [order_items] AS [T2] ON [T1].[order_id] = [T2].[order_id] LEFT JOIN [invoices] AS [T3] ON [T3].[order_id] = [T2].[order_id]"
		assert.Equal(b, expectedSql, joinRes.Syntax)
	}

}
func BenchmarkTestJoinByUsingDirectlyQueryable(b *testing.B) {
	repo := orm.Repository[OrderRepository]()
	ctx := orm.JoinCompiler.Ctx(mssql())
	for i := 0; i < b.N; i++ {

		innerJoin := repo.Invoices.OrderId.Eq(repo.OrderItems.OrderId).Join(
			repo.Invoices.Version.Eq(1).And( //<-- will be compile as join condition even this is AND not join
				repo.Invoices.CustomerId.Eq(repo.Customers.CustomerId), //<-- be cause new table appear in
			),
		)
		expectedInnerJoinClause := "[invoices] AS [T1] INNER JOIN [order_items] AS [T2] ON [T1].[order_id] = [T2].[order_id] AND [T1].[version] = ? INNER JOIN [customers] AS [T3] ON [T1].[customer_id] = [T3].[customer_id]"

		innerJoinClauseRes, err := ctx.ResolveBoolFieldAsJoin(nil, nil, innerJoin)
		assert.NoError(b, err)

		assert.Equal(b, expectedInnerJoinClause, innerJoinClauseRes.Syntax)
		assert.Equal(b, []interface{}{1}, innerJoinClauseRes.Args)
		leftJoin := repo.Invoices.OrderId.Eq(repo.OrderItems.OrderId).LeftJoin(
			repo.Invoices.Version.Eq(100).And( //<-- will be compile as join condition even this is AND not join
				repo.Invoices.CustomerId.Eq(repo.Customers.CustomerId), //<-- be cause new table appear in
			),
		)
		leftJoinExpectedSql := "[invoices] AS [T1] LEFT JOIN [order_items] AS [T2] ON [T1].[order_id] = [T2].[order_id] AND [T1].[version] = ? LEFT JOIN [customers] AS [T3] ON [T1].[customer_id] = [T3].[customer_id]"

		leftJoinRes, err := ctx.ResolveBoolFieldAsJoin(nil, nil, leftJoin)
		assert.NoError(b, err)

		assert.Equal(b, leftJoinExpectedSql, leftJoinRes.Syntax)
		assert.Equal(b, []interface{}{100}, leftJoinRes.Args)

		rightJoin := repo.Invoices.OrderId.Eq(repo.OrderItems.OrderId).RightJoin(
			repo.Invoices.Version.Eq(1).And( //<-- will be compile as join condition even this is AND not join
				repo.Invoices.CustomerId.Eq(repo.Customers.CustomerId), //<-- be cause new table appear in
			),
		)
		expectedRightJoinClause := "[invoices] AS [T1] RIGHT JOIN [order_items] AS [T2] ON [T1].[order_id] = [T2].[order_id] AND [T1].[version] = ? RIGHT JOIN [customers] AS [T3] ON [T1].[customer_id] = [T3].[customer_id]"

		rightJoinClauseRes, err := ctx.ResolveBoolFieldAsJoin(nil, nil, rightJoin)
		assert.NoError(b, err)

		assert.Equal(b, expectedRightJoinClause, rightJoinClauseRes.Syntax)
		assert.Equal(b, []interface{}{1}, rightJoinClauseRes.Args)
	}

}
