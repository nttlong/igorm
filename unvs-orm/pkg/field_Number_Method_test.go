package orm_test

import (
	"fmt"
	"reflect"
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestMethodOfNumberField(t *testing.T) {
	fnList := []string{
		"Eq",
		"Ne",
		"Gt",
		"Lt",
		"Ge",
		"Le",
		"IsNull",
		"IsNotNull",
		"Between",
		"NotBetween",
		"In",
		"NotIn",
		"Add",
		"Sub",
		"Mul",
		"Div",
		"Mod",
		"Pow",
		"As",
		"Text",
	}

	fn := orm.CreateNumberField[int]("table.name")
	// expr := fn.Eq(10)
	typ := reflect.TypeOf(&fn)
	assert.Equal(t, "*orm.NumberField[int]", typ.String())

	for _, mn := range fnList {
		if _, ok := typ.MethodByName(mn); !ok {
			t.Fatal(mn, fmt.Sprintf("method %s not found in %s", mn, typ.String()))

		}

	}

}
func TestNumberField(t *testing.T) {
	TestMethodOfNumberField(t)

	fn := orm.CreateNumberField[int]("table.name")
	cmp := orm.Compiler.Ctx(mssql())
	expr := fn.Eq(10)
	r, err := cmp.ResolveWithoutTableAlias(expr)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] = ?", r.Syntax)
	assert.Equal(t, []interface{}{10}, r.Args)
	expr1 := fn.Gt(10)
	r, err = cmp.ResolveWithoutTableAlias(expr1)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] > ?", r.Syntax)
	assert.Equal(t, []interface{}{10}, r.Args)
	expr2 := fn.Lt(10)
	r, err = cmp.ResolveWithoutTableAlias(expr2)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] < ?", r.Syntax)
	assert.Equal(t, []interface{}{10}, r.Args)
	expr3 := fn.Ge(10)
	r, err = cmp.ResolveWithoutTableAlias(expr3)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] >= ?", r.Syntax)
	assert.Equal(t, []interface{}{10}, r.Args)
	expr4 := fn.Le(10)
	r, err = cmp.ResolveWithoutTableAlias(expr4)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] <= ?", r.Syntax)
	assert.Equal(t, []interface{}{10}, r.Args)
	expr5 := fn.IsNull()
	r, err = cmp.ResolveWithoutTableAlias(expr5)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] IS NULL", r.Syntax)
	assert.Equal(t, []interface{}(nil), r.Args)
	expr6 := fn.IsNotNull()
	r, err = cmp.ResolveWithoutTableAlias(expr6)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] IS NOT NULL", r.Syntax)
	assert.Equal(t, []interface{}(nil), r.Args)
	expr7 := fn.Between(10, 20)
	r, err = cmp.ResolveWithoutTableAlias(expr7)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] BETWEEN ? AND ?", r.Syntax)
	assert.Equal(t, []interface{}{10, 20}, r.Args)
	expr8 := fn.NotBetween(10, 20)
	r, err = cmp.ResolveWithoutTableAlias(expr8)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] NOT BETWEEN ? AND ?", r.Syntax)
	assert.Equal(t, []interface{}{10, 20}, r.Args)
	expr9 := fn.In(10, 20, 30)
	r, err = cmp.ResolveWithoutTableAlias(expr9)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] IN (?,?,?)", r.Syntax)
	assert.Equal(t, []interface{}{10, 20, 30}, r.Args)
	expr10 := fn.NotIn(10, 20, 30)
	r, err = cmp.ResolveWithoutTableAlias(expr10)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] NOT IN (?,?,?)", r.Syntax)
	assert.Equal(t, []interface{}{10, 20, 30}, r.Args)
	expr11 := fn.Add(10)
	r, err = cmp.ResolveWithoutTableAlias(expr11)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] + ?", r.Syntax)
	assert.Equal(t, []interface{}{10}, r.Args)
	r, err = cmp.ResolveWithoutTableAlias(fn.Text())
	assert.NoError(t, err)
	assert.Equal(t, "CONVERT(NVARCHAR(50), [table].[name])", r.Syntax)
}
