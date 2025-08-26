package controllers

import (
	"fmt"
	"reflect"
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

type TestController struct {
}
type Service struct {
}
type ConatainerTest struct {
	wx.Inject[ConatainerTest]
}

func (c *ConatainerTest) Post(ctx *wx.Handler, container *ConatainerTest) {

}

func TestTestController(t *testing.T) {
	ok := wx.Helper.Inject.IsInjectType(reflect.TypeOf(ConatainerTest{}))
	assert.True(t, ok)
	ok = wx.Helper.Inject.IsInjectType(reflect.TypeOf(Service{}))
	assert.False(t, ok)
	mt := wx.GetMethodByName[ConatainerTest]("Post")
	assert.NotEmpty(t, *mt)
	mtInfo, err := wx.Helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	assert.NotNil(t, mtInfo)
	assert.Equal(t, -1, mtInfo.IndexOfRequestBody, "POST")
	assert.Equal(t, []int{2}, mtInfo.IndexOfArgIsInject, "POST")

}
func TestHalderWithInject(t *testing.T) {
	(&ConatainerTest{}).Register(func(svc *ConatainerTest) error {
		fmt.Println("call register")
		return nil
	})
	mt := wx.GetMethodByName[ConatainerTest]("Post")
	assert.NotEmpty(t, *mt)
	handler, err := wx.Helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	build := wx.Helper.ReqExec.CreateMockRequestBuilder()
	build.PostJson(handler.UriHandler, nil)
	build.ServerHandler(func() (any, error) {
		wx.Helper.ReqExec.Invoke(*handler, build.Req, &build.Res)
		fmt.Println("OK")
		return nil, nil
	})

}
func BenchmarkHalderWithInject(b *testing.B) {
	(&ConatainerTest{}).Register(func(svc *ConatainerTest) error {
		// fmt.Println("call register")
		return nil
	})
	mt := wx.GetMethodByName[ConatainerTest]("Post")
	assert.NotEmpty(b, *mt)
	handler, err := wx.Helper.GetHandlerInfo(*mt)
	assert.NoError(b, err)
	build := wx.Helper.ReqExec.CreateMockRequestBuilder()
	build.PostJson(handler.UriHandler, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		build.ServerHandler(func() (any, error) {
			wx.Helper.ReqExec.Invoke(*handler, build.Req, &build.Res)
			// fmt.Println("OK")
			return nil, nil
		})
	}
}
