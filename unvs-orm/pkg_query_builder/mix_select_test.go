package pkgquerybuilder

import (
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestSelectMixMode(b *testing.T) {
	dialect := mssql() //<-- create mssql dialect
	//ctx := orm.JoinCompiler.Ctx(mssql()) //<-- create compiler context for mssql dialect
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < 5; i++ {

		sql := repo.Orders.Select(
			//repo.Orders.Note,
			// repo.Orders.CreatedAt,
			// repo.Orders.UpdatedAt,
			// repo.Orders.CreatedBy,
			// repo.OrderItems.Product,
			repo.Expr("SUM(Orders.orderId) as total_quantity"),
		)
		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err())
		sqlExpected := "SELECT SUM ([orders].[order_id]) AS [total_quantity] FROM [orders]"
		assert.Equal(b, sqlExpected, compilerResult.String())
		assert.Equal(b, []interface{}(nil), compilerResult.Args)
	}

}
func BenchmarkTestSelectMixMode(b *testing.B) {
	//go test -bench=^BenchmarkTestSelectMixMode$ -benchmem -run=^$ ./pkg_query_builder -cpuprofile=cpu11.prof
	dialect := mssql() //<-- create mssql dialect
	//ctx := orm.JoinCompiler.Ctx(mssql()) //<-- create compiler context for mssql dialect
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < b.N; i++ {
		// for j := 0; j < 5; j++ {
		sql := repo.LeftJoin("Orders.OrderId = OrderItems.OrderId").Select(
			repo.Expr("SUM(Orders.orderId) as total_quantity"),
		)

		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err())
		sqlExpected := "SELECT SUM ([T1].[order_id]) AS [total_quantity] FROM [orders] AS [T1] LEFT JOIN [order_items] AS [T2] ON [T1].[order_id] = [T2].[order_id]"
		assert.Equal(b, sqlExpected, compilerResult.String())
		assert.Equal(b, []interface{}(nil), compilerResult.Args)
		// }
	}
}
