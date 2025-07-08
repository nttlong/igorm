package orm_test

import (
	"fmt"
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestNumberField_NumberField(t *testing.T) {
	cmp := orm.Compiler.Ctx(mssql())
	TestListOfFieldBinariesListOfMethods(t)

	fn := orm.CreateNumberField[int]("table.name")
	// vx := fn.Eq(fn)
	// compiler := orm.Compiler.Ctx(mssql())
	// r1, err1 := compiler.ResolveWithoutTableAlias(vx)
	// assert.NoError(t, err1)
	// assert.Equal(t, "[table].[name] = [table].[name]", r1.Syntax)
	testData := createListOfFieldBinaries(cmp, fn)

	for _, td := range testData {
		r, err := orm.Compiler.Ctx(mssql()).ResolveWithoutTableAlias(td.Expr)
		if td.Expected != r.Syntax {
			panic(fmt.Errorf("error at: %s, got: %s", td.Fn, td.Op))
		}
		assert.NoError(t, err)
		assert.Equal(t, td.Expected, r.Syntax)
		assert.Equal(t, 0, len(r.Args))
	}

}
