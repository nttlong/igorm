package tests

import (
	"fmt"
	"reflect"
	"testing"
	"wx"
	"wx/example/exampleControllers/controllers/media"

	"github.com/stretchr/testify/assert"
)

var PingMethod *reflect.Method

func TestGetPingMethod(t *testing.T) {
	mt := wx.GetMethodByName[media.Media]("Ping")
	assert.NotEmpty(t, mt)
	PingMethod = mt

}
func TestGetReceiver(t *testing.T) {
	TestGetPingMethod(t)
	ret, err := wx.Helper.GetReceiverTypeFromMethod(*PingMethod)
	controllerName := wx.Helper.FindControllerName(*ret)
	assert.Equal(t, controllerName, "media")
	assert.NoError(t, err)
	assert.NotEmpty(t, ret)
	assert.Equal(t, *ret, reflect.TypeOf(&media.Media{}))
	fmt.Println((*ret).Elem().PkgPath())
	t.Log((*ret).Elem().String())
}
func TestInspectPingMethod(t *testing.T) {
	TestGetPingMethod(t)
	infor, err := wx.Helper.GetHandlerInfo(*PingMethod)

	assert.NoError(t, err)
	assert.NotEmpty(t, infor)
	assert.Equal(t, infor.Uri, "media/ping")

}
