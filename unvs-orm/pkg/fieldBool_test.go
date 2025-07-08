package orm_test

import (
	"fmt"
	"reflect"
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestMethodOfBoolField(t *testing.T) {
	fnList := []string{
		"And",
		"Or",
		"Not",
	}

	fn := orm.CreateNumberField[int64]("table2.field2")
	expr := fn.Eq(10)
	typ := reflect.TypeOf(expr)
	assert.Equal(t, "*orm.BoolField", typ.String())

	for _, mn := range fnList {
		if _, ok := typ.MethodByName(mn); !ok {
			t.Fatal(mn, fmt.Sprintf("method %s not found in %s", mn, typ.String()))

		}

	}

}
func TestAnd(t *testing.T) {

	cmp := orm.Compiler.Ctx(mssql())

	fn := orm.CreateNumberField[int64]("table2.field2")
	fn2 := orm.CreateNumberField[int64]("table1.field1")
	expr := fn.Eq(10)
	r, err := cmp.ResolveWithoutTableAlias(expr)
	assert.NoError(t, err)
	assert.Equal(t, "[table2].[field2] = ?", r.Syntax)
	assert.Equal(t, []interface{}{10}, r.Args)
	expr = expr.And(fn2.Eq(20))
	r, err = cmp.ResolveWithoutTableAlias(expr)
	assert.NoError(t, err)
	assert.Equal(t, "[table2].[field2] = ? AND [table1].[field1] = ?", r.Syntax)
	assert.Equal(t, []interface{}{10, 20}, r.Args)

	expr = expr.Or(fn2.Ge(30))
	r, err = cmp.ResolveWithoutTableAlias(expr)
	assert.NoError(t, err)
	assert.Equal(t, "[table2].[field2] = ? AND [table1].[field1] = ? OR [table1].[field1] >= ?", r.Syntax)
}
