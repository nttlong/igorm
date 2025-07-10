package eorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEormJoin(t *testing.T) {
	ej := &exprJoin{

		context: &exprCompileContext{
			tables: []string{},
			alias:  map[string]string{},
			schema: map[string]bool{
				"User": true,
			},
			dialect:     dialectFactory.Create("mssql"),
			IsBuildJoin: true,
		},
	}
	err := ej.build("Departments INNER JOIN User ON User.Code = Departments.Code INNER JOIN Check ON Check.Name = 'John'")
	assert.NoError(t, err)
	assert.Equal(t, "[departments] AS [T1] INNER JOIN [User] AS [T2] ON [T2].[Code] = [T1].[code] INNER JOIN [checks] AS [T3] ON [T3].[name] = N'John'", ej.content)

}
func BenchmarkEormJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ej := &exprJoin{

			context: &exprCompileContext{
				tables: []string{},
				alias:  map[string]string{},
				schema: map[string]bool{
					"User": true,
				},
				dialect:     dialectFactory.Create("mssql"),
				IsBuildJoin: true,
			},
		}
		err := ej.build("Departments INNER JOIN User ON User.Code = Departments.Code INNER JOIN Check ON Check.Name = 'John'")
		if err != nil {
			b.Fail()
		}
	}
}
