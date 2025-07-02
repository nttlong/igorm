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
		"Set",
		"Get",
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
	}

	fn := orm.CreateDateTimeField("table.name")

	typ := reflect.TypeOf(fn)

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
	cmp := orm.Compiler.Ctx(mssql())
	fn := createDateTimeField("table.name")
	expr1 := fn.Eq("2021-01-01")
	r, err := cmp.Resolve(expr1)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] = ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	expr2 := fn.Ne("2021-01-01")
	r, err = cmp.Resolve(expr2)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] != ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	expr3 := fn.Lt("2021-01-01")
	r, err = cmp.Resolve(expr3)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] < ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	expr4 := fn.Le("2021-01-01")
	r, err = cmp.Resolve(expr4)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] <= ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	expr5 := fn.Gt("2021-01-01")
	r, err = cmp.Resolve(expr5)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] > ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	expr6 := fn.Ge("2021-01-01")
	r, err = cmp.Resolve(expr6)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] >= ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	expr7 := fn.IsNull()
	r, err = cmp.Resolve(expr7)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] IS NULL", r.Syntax)
	expr8 := fn.IsNotNull()
	r, err = cmp.Resolve(expr8)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] IS NOT NULL", r.Syntax)
	expr9 := fn.Between("2021-01-01", "2021-01-31")
	r, err = cmp.Resolve(expr9)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] BETWEEN ? AND ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	assert.Equal(t, "2021-01-31", r.Args[1])
	expr10 := fn.NotBetween("2021-01-01", "2021-01-31")
	r, err = cmp.Resolve(expr10)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "[table].[name] NOT BETWEEN ? AND ?", r.Syntax)
	assert.Equal(t, "2021-01-01", r.Args[0])
	assert.Equal(t, "2021-01-31", r.Args[1])
	expr11 := fn.Day()
	r, err = cmp.Resolve(expr11)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "DAY([table].[name])", r.Syntax)
	expr12 := fn.Month()
	r, err = cmp.Resolve(expr12)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "MONTH([table].[name])", r.Syntax)
	expr13 := fn.Year()
	r, err = cmp.Resolve(expr13)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "YEAR([table].[name])", r.Syntax)
	expr14 := fn.Hour()
	r, err = cmp.Resolve(expr14)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "HOUR([table].[name])", r.Syntax)
	expr15 := fn.Minute()
	r, err = cmp.Resolve(expr15)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "MINUTE([table].[name])", r.Syntax)
	expr16 := fn.Second()
	r, err = cmp.Resolve(expr16)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "SECOND([table].[name])", r.Syntax)
	expr17 := fn.Format("YYYY-MM-DD")
	r, err = cmp.Resolve(expr17)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "FORMAT([table].[name],?)", r.Syntax)
	assert.Equal(t, "YYYY-MM-DD", r.Args[0])

}
