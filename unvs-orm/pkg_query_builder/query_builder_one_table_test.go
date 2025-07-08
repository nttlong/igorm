package pkgquerybuilder

import (
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestSelectOneTable(b *testing.T) {
	dialect := mssql() //<-- create mssql dialect
	//ctx := orm.JoinCompiler.Ctx(mssql()) //<-- create compiler context for mssql dialect
	repo := orm.Repository[OrderRepository]()
	sql := repo.Orders.Select(
		// repo.Orders.Note,
		// repo.Orders.CreatedAt,
		// repo.Orders.UpdatedAt,
		// repo.Orders.CreatedBy,
		// repo.OrderItems.Product,
		repo.Expr("SUM(Orders.orderId) as total_quantity"),
	)
	compilerResult, err := sql.ToSql(dialect)
	assert.NoError(b, err)
	sqlExpected := "SELECT [orders].[note], [orders].[created_at], [orders].[updated_at], [orders].[created_by], [order_items].[product] FROM [orders] WHERE [orders].[order_id] = ?"
	assert.Equal(b, sqlExpected, compilerResult.Sql)
	assert.Equal(b, []interface{}{1}, compilerResult.Args)

}

func BenchmarkTestSelectOneTable(b *testing.B) {
	dialect := mssql() //<-- create mssql dialect
	//ctx := orm.JoinCompiler.Ctx(mssql()) //<-- create compiler context for mssql dialect
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < b.N; i++ {
		sql := repo.Orders.Select(
			repo.Orders.Note,
			repo.Orders.CreatedAt,
			repo.Orders.UpdatedAt,
			repo.Orders.CreatedBy,
			repo.OrderItems.Product,
		)
		compilerResult, err := sql.ToSql(dialect)
		assert.NoError(b, err)
		sqlExpected := "SELECT [orders].[note], [orders].[created_at], [orders].[updated_at], [orders].[created_by], [order_items].[product] FROM [orders]"
		assert.Equal(b, sqlExpected, compilerResult.Sql)
		assert.Equal(b, []interface{}{}, compilerResult.Args)
	}

}
