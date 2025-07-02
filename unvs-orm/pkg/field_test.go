package orm_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func mssql() orm.DialectCompiler {
	return &orm.MssqlDialect
}
func TestFieldMssql(t *testing.T) {

	f := orm.CreateNumberField[int64]("table.name")
	r, err := orm.Compiler.Ctx(mssql()).Resolve(f)
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name]", r.Syntax)
	assert.Equal(t, 0, len(r.Args))
}
func TestListOfFieldBinariesListOfMethods(t *testing.T) {
	fnList := []string{
		"Eq",
		"Ne",
		"Gt",
		"Lt",
		"Ge",
		"Le",
		"As",
		"Set",
		"Get",
		"Add",
		"Sub",
		"Mul",
		"Div",
		"Mod",
		"IsNull",
		"IsNotNull",
		"In",
		"NotIn",
		"Between",
		"NotBetween",
	}

	fn := orm.CreateNumberField[int64]("table.name")

	typ := reflect.TypeOf(&fn)

	for _, fnx := range fnList {
		if _, ok := typ.MethodByName(fnx); !ok {
			t.Error(fmt.Errorf("method %s was not found in NumberField[int64]", fnx))
		}

	}

}
func TestFieldAlias(t *testing.T) {

	TestFieldMssql(t)
	TestListOfFieldBinariesListOfMethods(t)

	f := orm.CreateNumberField[int64]("table.name")

	r, err := orm.Compiler.Ctx(mssql()).Resolve(f.As("alias"))
	assert.NoError(t, err)
	assert.Equal(t, "[table].[name] AS [alias]", r.Syntax)
	assert.Empty(t, r.Args)
}

type structTest struct {
	Expr     interface{}
	Expected string
}

func createListOfFieldBinaries(cmp *orm.CompilerUtils, testVal interface{}) []structTest {
	fnList := []string{
		"Eq->=",
		"Ne->!=",
		"Gt->>",
		"Lt-><",
		"Ge->>=",
		"Le-><=",
	}
	ret := []structTest{}

	fn := orm.CreateNumberField[int64]("table.name")
	fnVal := reflect.ValueOf(fn)
	typ := reflect.TypeOf(fn)

	for _, fex := range fnList {
		fx := strings.Split(fex, "->")[0]
		fxe := strings.Split(fex, "->")[1]

		if m, ok := typ.MethodByName(fx); ok {
			// m := fnVal.MethodByName(fx)
			// mc := m.Func.Call()

			outPut := m.Func.Call([]reflect.Value{fnVal, reflect.ValueOf(testVal)})
			Expected, err := cmp.Resolve(testVal)
			if err != nil {
				panic(err)
			}
			if Expected.Syntax == "?" {
				ret = append(ret, structTest{
					Expr:     outPut[0].Interface(),
					Expected: fmt.Sprintf("[table].[name] %s ?", fxe),
				})
			} else {
				ret = append(ret, structTest{
					Expr:     outPut[0].Interface(),
					Expected: fmt.Sprintf("[table].[name] %s %s", fxe, Expected.Syntax),
				})
			}
		} else {
			panic(fmt.Errorf("method %s was not found in NumberField[int64]", fx))
		}

	}
	return ret
}

func TestBinaryField(t *testing.T) {
	TestFieldAlias(t)
	cmp := orm.Compiler.Ctx(mssql())
	testData := createListOfFieldBinaries(cmp, 123)
	for _, td := range testData {
		r, err := orm.Compiler.Ctx(mssql()).Resolve(td.Expr)
		assert.NoError(t, err)
		assert.Equal(t, td.Expected, r.Syntax)
		assert.Equal(t, 1, len(r.Args))
	}

}
