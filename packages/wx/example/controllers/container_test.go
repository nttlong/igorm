package controllers

import (
	"fmt"
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

type TestController struct {
}
type Service struct {
}
type ConatainerTest struct {
	wx.Service
}

func (c *ConatainerTest) New() error {
	return nil
}

type HanderService struct {
	wx.Handler
}

//	func (lst *HanderService) New() error {
//		return nil
//	}
func (c *TestController) Post(ctx *struct {
	HanderService `route:"method:post"`
}, container *ConatainerTest) {

}

func TestHalderWithInject(t *testing.T) {
	// (&ConatainerTest{}).Register(func(svc *ConatainerTest) error {
	// 	fmt.Println("call register")
	// 	return nil
	// })
	mt := wx.GetMethodByName[TestController]("Post")
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
	// (&ConatainerTest{}).Register(func(svc *ConatainerTest) error {
	// 	// fmt.Println("call register")
	// 	return nil
	// })
	mt := wx.GetMethodByName[TestController]("Post")
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
