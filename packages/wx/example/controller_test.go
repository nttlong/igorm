package example

import (
	"testing"
	"wx"
	"wx/handlers"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	mt := wx.GetMethodByName[Media]("ListOfFiles")

	mtInfo, err := handlers.Helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuild.PostJson("/api/"+mtInfo.UriHandler, nil)
	req, res := requestBuild.Build()

	wx.Helper.ReqExec.DoJsonPost(*mtInfo, req, res)
	// wx.LoadController(func() (*Media, error) {
	// 	return &Media{}, nil
	// })
	// assert.NoError(t, err)
	t.Log(mtInfo)
}
func BenchmarkTestGet(b *testing.B) {
	mt := wx.GetMethodByName[Media]("ListOfFiles")

	mtInfo, _ := handlers.Helper.GetHandlerInfo(*mt)

	requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuild.PostJson("/api/"+mtInfo.UriHandler, nil)
	req, res := requestBuild.Build()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wx.Helper.ReqExec.DoJsonPost(*mtInfo, req, res)
		}

	})

}
