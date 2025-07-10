package pkgquerybuilder

import (
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

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
		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err())
		sqlExpected := "SELECT [orders].[note], [orders].[created_at], [orders].[updated_at], [orders].[created_by], [order_items].[product] FROM [orders]"
		assert.Equal(b, sqlExpected, compilerResult.String())
		assert.Equal(b, []interface{}{}, compilerResult.Args)
	}

}
