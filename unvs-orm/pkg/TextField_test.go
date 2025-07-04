package orm_test

import (
	"reflect"
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func TestListOfFieldTextFields(t *testing.T) {
	fnList := []string{
		"Eq",
		"Ne",
		"As",

		"Like",
		"NotLike",
		"In",
		"NotIn",
		"IsNull",
		"IsNotNull",
		"Between",
		"NotBetween",
		"Lower",
		"Upper",
		"Trim",
		"LTrim",
		"RTrim",
		"Len",
		"Concat",
	}

	fn := orm.CreateTextField("table.name")

	typ := reflect.TypeOf(&fn)

	for _, fn := range fnList {
		if _, ok := typ.MethodByName(fn); !ok {
			t.Errorf("method %s not found in type %s", fn, typ)
		}
	}

}
func createTextField(fullName string) *orm.TextField {

	fn := orm.CreateTextField(fullName)
	return &fn
}

func TestTextField(t *testing.T) {
	cmp := orm.Compiler.Ctx(mssql())
	fn := createTextField("table.name")
	expr := fn.IsNotNull()
	r, err := cmp.Resolve(nil, expr)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] IS NOT NULL", r.Syntax)
	expr1 := fn.Eq("abc")
	r, err = cmp.Resolve(nil, expr1)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] = ?", r.Syntax)
	assert.Equal(t, "abc", r.Args[0])
	expr2 := fn.In([]string{"a", "b", "c"})
	r, err = cmp.Resolve(nil, expr2)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] IN (?,?,?)", r.Syntax)
	fnFirstName := createTextField("table.first_name")
	fnLastName := createTextField("table.last_name")
	expr3 := fnFirstName.Concat(" ", fnLastName)
	r, err = cmp.Resolve(nil, expr3)
	assert.NoError(t, err)
	assert.Equal(t, "CONCAT([table].[first_name],?,[table].[last_name])", r.Syntax)
	assert.Equal(t, " ", r.Args[0])

}
