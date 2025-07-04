package orm_test

import (
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestNumberField_NumberField(t *testing.T) {
	cmp := orm.Compiler.Ctx(mssql())
	TestListOfFieldBinariesListOfMethods(t)

	fn := orm.CreateNumberField[int]("table.name")
	testData := createListOfFieldBinaries(cmp, fn)

	for _, td := range testData {
		r, err := orm.Compiler.Ctx(mssql()).Resolve(nil, td.Expr)
		assert.NoError(t, err)
		assert.Equal(t, td.Expected, r.Syntax)
		assert.Equal(t, 0, len(r.Args))
	}

}
