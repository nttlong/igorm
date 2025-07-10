package pkgquerybuilder

import (
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestAliAsTableSelfJoin(t *testing.T) {
	repo := orm.Repository[OrderRepository]()
	customer2 := repo.Customers.Alias("customer2")

	sql := customer2.Name.RightJoin(repo.Customers.Name).Select(customer2.Name)
	dialect := mssql()
	compilerResult := sql.Compile(dialect)
	assert.NoError(t, compilerResult.Err())
	sqlExpected := "SELECT [T1].[name] AS [name] FROM [customers] AS [T2] RIGHT JOIN [customers] AS [T1] ON [T1].[name] = [T2].[name]"
	assert.Equal(t, sqlExpected, compilerResult.String())
}
func TestComplexJoin1(t *testing.T) {
	dialect := mssql() //<-- create mssql dialect
	//ctx := orm.JoinCompiler.Ctx(mssql()) //<-- create compiler context for mssql dialect
	repo := orm.Repository[OrderRepository]()

	joinOneTo2Tables := repo.Orders.OrderId.LeftJoin( //<-- join Order and OrderItem tables and select Order.Note, Order.CreatedAt, Order.UpdatedAt, Order.CreatedBy, OrderItem.Product
		repo.OrderItems.OrderId,
		repo.Customers.CustomerId,
		repo.Invoices.CreatedAt.Year(),
	)
	sql := joinOneTo2Tables.Select( //<-- from the joined tables select Order.Note, Order.CreatedAt, Order.UpdatedAt, Order.CreatedBy, OrderItem.Product
		repo.Orders.Note,
		repo.Orders.CreatedAt,
		repo.Orders.UpdatedAt,
		repo.Orders.CreatedBy,
		repo.OrderItems.Product,
	)
	compilerResult := sql.Compile(dialect)
	assert.NoError(t, compilerResult.Err())
	sqlExpected := "SELECT [T1].[note] AS [note], [T1].[created_at] AS [created_at], [T1].[updated_at] AS [updated_at], [T1].[created_by] AS [created_by], [T2].[product] AS [product] FROM [orders] AS [T1] LEFT JOIN [order_items] AS [T2] ON [T1].[order_id] = [T2].[order_id] LEFT JOIN [customers] AS [T3] ON [T1].[order_id] = [T3].[customer_id] LEFT JOIN [invoices] AS [T4] ON [T1].[order_id] = YEAR([T4].[created_at])"
	assert.Equal(t, sqlExpected, compilerResult.String())
	sqlText := compilerResult.String()

	assert.Equal(t, []interface{}(nil), compilerResult.Args)
	assert.Equal(t, sqlText, compilerResult.String())

}
func TestComplexJoin2(t *testing.T) {
	dialect := mssql() //<-- create mssql dialect
	//ctx := orm.JoinCompiler.Ctx(mssql()) //<-- create compiler context for mssql dialect
	repo := orm.Repository[OrderRepository]()
	sql := repo.Customers.Note.Like(repo.Customers.Email).Select(repo.Customers.Note)
	compilerResult := sql.Compile(dialect)
	assert.NoError(t, compilerResult.Err())
	sqlExpected := "SELECT [T1].[note] AS [note] FROM [customers] AS [T1] WHERE [T1].[note] LIKE [T1].[email]"
	assert.Equal(t, sqlExpected, compilerResult.String())

}
