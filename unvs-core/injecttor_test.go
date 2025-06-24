package unvscore

import (
	"dbx"
	"strconv"
	"testing"

	di "unvs.di"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Code        di.Scoped[TestStruct, int]
	Name        di.Scoped[TestStruct, string]
	Description di.Transient[TestStruct, string]
}
type FN func(arg ...interface{}) interface{}
type BaseService struct {
	db Singleton[BaseService, *dbx.DBX]
}
type TestStruct2 struct {
	db *dbx.DBXTenant
}

// o 1 dich vu nao do khi goi TestStruct2.GetUserName() thai can co database
// o 1 dich vu khac neu chi su dung TestStruct2.SaveFile() thi kg viec gi phai connect den db
func TestInject(t *testing.T) {

	fx := Resolve[TestStruct](struct {
		Code
	})
	fx.Code.Init = func(owner TestStruct) int {
		return 123
	}
	fx.Name.Init = func(owner TestStruct) string {
		return "Hello World " + strconv.Itoa(owner.Code.Get())
	}
	fx.Description.OnGet = func(owner TestStruct) string {
		return "Description " + strconv.Itoa(owner.Code.Get())
	}

	if fx.Code.Owner == fx.Name.Owner {
		assert.Equal(t, fx.Code.Owner, fx.Name.Owner)
	}

	v := fx.Name.Get()
	assert.Equal(t, "Hello World 123", v)
	assert.Equal(t, 123, fx.Code.Get())
	v = fx.Description.Get()
	assert.Equal(t, "Hello World 123", v)

}
