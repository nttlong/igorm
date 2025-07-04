package orm_test

import (
	"reflect"
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestAllMethodsAreImplementedOfDateTimeField(t *testing.T) {
	fnList := []string{
		"Eq",
		"Ne",
		"Lt",
		"Le",
		"Gt",
		"Ge",
		"As",

		"In",
		"NotIn",
		"IsNull",
		"IsNotNull",
		"Between",
		"NotBetween",
		"Day",
		"Month",
		"Year",
		"Hour",
		"Minute",
		"Second",
		"Format",
		"Min",
		"Max",
		"Count",
	}

	fn := orm.CreateDateTimeField("table.name")

	typ := reflect.TypeOf(&fn)

	for _, fn := range fnList {
		if _, ok := typ.MethodByName(fn); !ok {
			t.Errorf("method %s not found in type %s", fn, typ)
		}
	}

}
func createDateTimeField(fullName string) *orm.DateTimeField {
	ret := orm.CreateDateTimeField(fullName)
	return &ret
}
func TestDateTiemField(t *testing.T) {
	TestAllMethodsAreImplementedOfDateTimeField(t)
	cmp := orm.Compiler.Ctx(mssql())
	fn := createDateTimeField("table.name")
	expr1 := fn.Eq("2021-01-01")
	r, err := cmp.Resolve(nil, expr1)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] = ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	expr2 := fn.Ne("2021-01-01")
	r, err = cmp.Resolve(nil, expr2)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] != ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	expr3 := fn.Lt("2021-01-01")
	r, err = cmp.Resolve(nil, expr3)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] < ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	expr4 := fn.Le("2021-01-01")
	r, err = cmp.Resolve(nil, expr4)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] <= ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	expr5 := fn.Gt("2021-01-01")
	r, err = cmp.Resolve(nil, expr5)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] > ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	expr6 := fn.Ge("2021-01-01")
	r, err = cmp.Resolve(nil, expr6)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] >= ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	expr7 := fn.IsNull()
	r, err = cmp.Resolve(nil, expr7)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] IS NULL", r.Syntax)
	expr8 := fn.IsNotNull()
	r, err = cmp.Resolve(nil, expr8)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] IS NOT NULL", r.Syntax)
	expr9 := fn.Between("2021-01-01", "2021-01-31")
	r, err = cmp.Resolve(nil, expr9)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] BETWEEN ? AND ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	assert.Equal(t, "2021-01-31", r.Args[1])
	expr10 := fn.NotBetween("2021-01-01", "2021-01-31")
	r, err = cmp.Resolve(nil, expr10)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] NOT BETWEEN ? AND ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	assert.Equal(t, "2021-01-31", r.Args[1])
	expr11 := fn.Day()
	r, err = cmp.Resolve(nil, expr11)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "DAY([table].[name])", r.Syntax)
	expr12 := fn.Month()
	r, err = cmp.Resolve(nil, expr12)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "MONTH([table].[name])", r.Syntax)
	expr13 := fn.Year()
	r, err = cmp.Resolve(nil, expr13)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "YEAR([table].[name])", r.Syntax)
	expr14 := fn.Hour()
	r, err = cmp.Resolve(nil, expr14)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "HOUR([table].[name])", r.Syntax)
	expr15 := fn.Minute()
	r, err = cmp.Resolve(nil, expr15)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "MINUTE([table].[name])", r.Syntax)
	expr16 := fn.Second()
	r, err = cmp.Resolve(nil, expr16)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "SECOND([table].[name])", r.Syntax)
	expr17 := fn.Format("YYYY-MM-DD")
	r, err = cmp.Resolve(nil, expr17)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "FORMAT([table].[name],?)", r.Syntax)
	assert.Equal(t, "YYYY-MM-DD", r.Args[0])
	expr18 := fn.Min()
	r, err = cmp.Resolve(nil, expr18)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "MIN([table].[name])", r.Syntax)
	expr19 := fn.Max()
	r, err = cmp.Resolve(nil, expr19)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "MAX([table].[name])", r.Syntax)
	expr20 := fn.Count()
	r, err = cmp.Resolve(nil, expr20)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "COUNT([table].[name])", r.Syntax)

}
