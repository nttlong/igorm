package fapi

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ObjTest struct {
}

func (obj *ObjTest) Upload(c *Context) {
	panic("implement me")

}

//	func (obj *ObjTest) ApiMethod1(c *ProfileUpdate, getUser Inject[Singleton[AppConfig]]) struct {
//		Code string
//		Msg  string
//	} {
//
//		panic("Test")
//	}

func TestApiUtils(t *testing.T) {

	typ := reflect.TypeOf(&ObjTest{})
	mt := api.InspectMethod(typ.Method(0))
	assert.NotEmpty(t, mt)
	assert.Equal(t, "ApiMethod1", mt.method.Name)
	assert.Equal(t, false, mt.IsAbsUri())
	uri := mt.GetRoute()
	fmt.Println(uri)
	assert.Equal(t, true, mt.HasUploadFile())
	assert.Equal(t, true, mt.HasUploadFile())

}
