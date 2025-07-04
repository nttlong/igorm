package pkgquerybuilder

import (
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func mssql() orm.DialectCompiler {
	return &orm.MssqlDialect
}
func TestOrderQuery(t *testing.T) {

	repo := orm.Repository[OrderRepository]()

	w := repo.Orders.OrderId.Eq(1).And(
		repo.Orders.Note.Eq("test"),
	)
	orderQuery := orm.From(
		repo.Orders, //<-- repo.Orders.Join(repo.OrderItems, repo.Orders.OrderId.Eq(repo.OrderItems.OrderId))
	).Where(
		w,
	).Select(repo.Orders.OrderId.Max().As("max_order_id"),
		repo.Orders.Note,
	).GroupBy(repo.Orders.Note).Having(
		repo.Orders.Note.Eq("test"),
	)
	sql, err := orderQuery.ToSql(mssql())
	expectedSql := "SELECT MAX([T0].[order_id]) AS [max_order_id], [T0].[note] FROM orders WHERE [T0].[order_id] = ? AND [T0].[note] = ? GROUP BY [T0].[note] HAVING [T0].[note] = ?"
	assert.Empty(t, err)
	assert.Equal(t, expectedSql, sql.Sql)
	assert.Equal(t, []interface{}{1, "test", "test"}, sql.Args)

}
func TestJoinExpre(b *testing.T) {
	repo := orm.Repository[OrderRepository]()
	on := repo.Orders.OrderId.Eq(repo.OrderItems.OrderId)
	join := repo.Orders.Join(repo.OrderItems, on)
	ctx := orm.JoinCompiler.Ctx(mssql())
	joinRes, err := ctx.Resolve(join)
	assert.NoError(b, err)
	expectedSql := "[orders] AS [T1] INNER JOIN [order_items] AS [T2] ON [T1].[order_id] = [T2].[order_id]"
	assert.Equal(b, expectedSql, joinRes.Syntax)

}
func TestJoinExpr2(b *testing.T) {
	repo := orm.Repository[OrderRepository]()
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
