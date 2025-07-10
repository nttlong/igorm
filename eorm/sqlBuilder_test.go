package eorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
