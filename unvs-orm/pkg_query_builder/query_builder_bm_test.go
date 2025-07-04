package pkgquerybuilder

import (
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func BenchmarkQueryBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		repo := orm.Repository[OrderRepository]()

		w := repo.Orders.OrderId.Eq(1).And(
			repo.Orders.Note.Eq("test"),
		)
		orderQuery := orm.From(repo.Orders).Where(
			w,
		).Select(repo.Orders.OrderId.Max().As("max_order_id"),
			repo.Orders.Note,
		).GroupBy(repo.Orders.Note).Having(
			repo.Orders.Note.Eq("test"),
		)
		sql, err := orderQuery.ToSql(mssql())
		expectedSql := "SELECT MAX([T0].[order_id]) AS [max_order_id], [T0].[note] FROM orders WHERE [T0].[order_id] = ? AND [T0].[note] = ? GROUP BY [T0].[note] HAVING [T0].[note] = ?"
		assert.Empty(b, err)
		assert.Equal(b, expectedSql, sql.Sql)
		assert.Equal(b, []interface{}{1, "test", "test"}, sql.Args)
	}

}
