package mockcontroller

import (
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

type UrlPathTest struct {
}

func (u *UrlPathTest) New() error {

	// Initialize if needed
	return nil
}
func (u *UrlPathTest) Download(ctx *struct {
	wx.Handler `route:"uri:@/{*FilePath};method:get"`
	FilePath   string
}) (any, error) {
	// Here you would implement the logic to handle the download
	// For now, just return nil to indicate success
	return &struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    "200",
		Message: "success",
	}, nil
}
func TestDownload(t *testing.T) {
	// Create a new instance of UrlPathTest
	method := wx.GetMethod[UrlPathTest]("Download")
	assert.NotEmpty(t, method, "Method should not be empty")
	info, err := wx.Helper.GetHandlerInfo(*method)
	assert.NotEmpty(t, info, "Info should not be empty")
	assert.NoError(t, err, "Should not return an error")
	builder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	//builder.Get("/api/" + info.UriHandler[0:len(info.UriHandler)-1]) // Assuming the last part is the file path
	// builder.Handler(func(w http.ResponseWriter, r *http.Request) {
	// 	wx.Helper.ReqExec.Invoke(*info, r, w)
	// })
	builder.Get("/api/" + info.UriHandler[0:len(info.UriHandler)-1]) // Example file path
	builder.ServerHandler(func() (any, error) {
		data, err := wx.Helper.ReqExec.Invoke(*info, builder.Req, &builder.Res)
		return data, err
	})
	assert.Equal(t, 404, builder.Res.Code, "Response code should be 200")
	builder.Get("/api/" + info.UriHandler) // Example file path
	builder.ServerHandler(func() (any, error) {
		return wx.Helper.ReqExec.Invoke(*info, builder.Req, &builder.Res)
	})
	assert.Equal(t, 200, builder.Res.Code, "Response code should be 200")

}
