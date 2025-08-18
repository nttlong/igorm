package tests

import (
	"fmt"
	"reflect"
	"testing"
	"wx"
	"wx/example/exampleControllers/controllers/media"

	"github.com/stretchr/testify/assert"
)

var Ping2Method *reflect.Method

func TestGetPing2Method(t *testing.T) {
	mt := wx.GetMethodByName[media.Media]("Ping2")
	assert.NotEmpty(t, mt)
	Ping2Method = mt

}
func TestGetReceiverFromPing2(t *testing.T) {
	TestGetPing2Method(t)
	ret, err := wx.Helper.GetReceiverTypeFromMethod(*Ping2Method)
	controllerName := wx.Helper.FindControllerName(*ret)
	assert.Equal(t, controllerName, "media")
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)
	assert.Equal(t, *ret, reflect.TypeOf(&media.Media{}))
	fmt.Println((*ret).Elem().PkgPath())
	t.Log((*ret).Elem().String())
}
func TestInspectPing2Method(t *testing.T) {
	TestGetPing2Method(t)
	infor, err := wx.Helper.GetHandlerInfo(*Ping2Method)

	assert.NoError(t, err)
	assert.NotEmpty(t, infor)
	assert.Equal(t, "media/test", infor.Uri)

}
