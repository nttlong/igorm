package eorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelfJoin(t *testing.T) {
	for i := 0; i < 5; i++ {

		join := "Order AS Child INNER JOIN Order AS Parent ON Child.ParentCode = Parent.Code"
		builder := SqlBuilder.From(join).Select()
		sql, args := builder.ToSql(dialectFactory.Create("mssql"))
		assert.NoError(t, builder.Err)
		assert.Equal(t, "SELECT [Child].*, [Parent].* FROM [orders] AS [Child] INNER JOIN [orders] AS [Parent] ON [Child].[ParentCode] = [Parent].[Code]", sql)
		assert.Equal(t, 0, len(args))
	}
}
func TestSqlBuilderSelectField(t *testing.T) {
	for i := 0; i < 5; i++ {
		joinExpr := "Departments INNER JOIN User ON User.Code = Departments.Code INNER JOIN Check ON Check.Name = 'John'"
		builder := SqlBuilder.From(joinExpr).Select(
			`Departments.Code AS DepartmentCode,
			User.Name AS UserName,
			Check.Name as CheckName,
			concat(User.FirstName,' ',User.LastName) AS FullName,`,

			"", 12,
		)
		sql, args := builder.ToSql(dialectFactory.Create("mssql"))
		assert.NoError(t, builder.Err)
		assert.Equal(t, "SELECT [T1].[code] AS [DepartmentCode], [T2].[name] AS [UserName], [T3].[name] AS [CheckName], CONCAT([T2].[first_name], N' ', [T2].[last_name]) AS [FullName] FROM [departments] AS [T1] INNER JOIN [users] AS [T2] ON [T2].[code] = [T1].[code] INNER JOIN [checks] AS [T3] ON [T3].[name] = N'John'", sql)
		assert.Equal(t, 2, len(args))
	}
}
func BenchmarkBuilderSelectField(b *testing.B) {
	for i := 0; i < b.N; i++ {
		joinExpr := "Departments INNER JOIN User ON User.Code = Departments.Code INNER JOIN Check ON Check.Name = 'John'"
		builder := SqlBuilder.From(joinExpr).Select(
			`Departments.Code AS DepartmentCode,
			User.Name AS UserName,
			Check.Name as CheckName,
			concat(User.FirstName,' ',User.LastName) AS FullName,`,

			"", 12,
		)
		sql, args := builder.ToSql(dialectFactory.Create("mssql"))
		assert.NoError(b, builder.Err)
		assert.Equal(b, "SELECT [T1].[code] AS [DepartmentCode], [T2].[name] AS [UserName], [T3].[name] AS [CheckName], CONCAT([T2].[first_name], N' ', [T2].[last_name]) AS [FullName] FROM [departments] AS [T1] INNER JOIN [users] AS [T2] ON [T2].[code] = [T1].[code] INNER JOIN [checks] AS [T3] ON [T3].[name] = N'John'", sql)
		assert.Equal(b, 2, len(args))
	}
}
func TestSqlBuilder(t *testing.T) {
	for i := 0; i < 5; i++ {
		joinExpr := "Departments INNER JOIN User ON User.Code = Departments.Code INNER JOIN Check ON Check.Name = 'John'"
		builder := SqlBuilder.From(joinExpr).Select()
		sql, args := builder.ToSql(dialectFactory.Create("mssql"))
		assert.NoError(t, builder.Err)
		assert.Equal(t, "SELECT [T1].*, [T2].*, [T3].* FROM [departments] INNER JOIN [users] ON [T1].[code] = [T2].[code] INNER JOIN [checks] ON [T3].[name] = N'John'", sql)
		assert.Equal(t, 0, len(args))
	}

}
func BenchmarkSqlBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		joinExpr := "Departments INNER JOIN User ON User.Code = Departments.Code INNER JOIN Check ON Check.Name = 'John'"
		builder := SqlBuilder.From(joinExpr).Select()
		sql, args := builder.ToSql(dialectFactory.Create("mssql"))
		if builder.Err != nil {
			b.Error(builder.Err)
		}
		if sql != "SELECT [T1].*, [T2].*, [T3].* FROM [departments] INNER JOIN [users] ON [T1].[code] = [T2].[code] INNER JOIN [checks] ON [T3].[name] = N'John'" {
			b.Error("sql is not equal")
		}
		if len(args) != 0 {
			b.Error("args is not empty")
		}
	}
}
