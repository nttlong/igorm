package mockcontroller

import (
	"net/http"
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

type ControllerUriParams struct {
}

func (ctl *ControllerUriParams) Post(ctx *struct {
	wx.Handler `route:"{param1}/{param2}"`
	Param1     string
	Param2     string
}) (any, error) {

	return nil, nil
}
func TestControllerUriParamsPost(t *testing.T) {
	mt := wx.GetMethodByName[ControllerUriParams]("Post")
	assert.NotEmpty(t, mt)
	info, err := wx.Helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	assert.Equal(t, "controller-uri-params/", info.UriHandler)
	assert.Equal(t, "POST", info.HttpMethod)
	assert.Equal(t, "controller-uri-params/{param1}/{param2}", info.Uri)
	reqBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	reqBuilder.PostJson("api/"+info.UriHandler+"abc/cde", nil)
	reqBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {
		wx.Helper.ReqExec.CreateHandlerContext(*info, r, w)
		wx.Helper.ReqExec.DoJsonPost(*info, r, w)
	})

}
func BenchmarkControllerUriParamsPost(b *testing.B) {
	mt := wx.GetMethodByName[ControllerUriParams]("Post")

	info, _ := wx.Helper.GetHandlerInfo(*mt)

	reqBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	reqBuilder.PostJson("api/"+info.UriHandler+"abc/cde", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reqBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {

			wx.Helper.ReqExec.DoJsonPost(*info, r, w)
		})
	}
}
func (ctl *ControllerUriParams) PostQuery(ctx *struct {
	wx.Handler `route:"?param1={param1}&param2={param2}"`
	Param1     []string
	Param2     []string
}) (any, error) {

	return nil, nil
}
func TestControllerPostQuery(t *testing.T) {
	mt := wx.GetMethodByName[ControllerUriParams]("PostQuery")
	assert.NotEmpty(t, mt)
	info, err := wx.Helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	assert.Equal(t, "controller-uri-params", info.UriHandler)
	assert.Equal(t, "POST", info.HttpMethod)
	assert.Equal(t, "controller-uri-params?param1={param1}&param2={param2}", info.Uri)
	reqBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	reqBuilder.PostJson("api/"+info.UriHandler+"?param1=abc&param1=cde", nil)
	reqBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {

		wx.Helper.ReqExec.Invoke(*info, r, w)
	})

}
