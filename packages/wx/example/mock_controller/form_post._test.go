package mockcontroller

import (
	"fmt"
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

type FormPostTest struct {
}

func (f *FormPostTest) Post(ctx *wx.Handler, form wx.Form[struct {
	Username string `json:"username"`
	Password string `json:"password"`
}]) {
	fmt.Println(form.Data.Username)
	fmt.Println(form.Data.Password)
}
func TestFormPost(t *testing.T) {
	method := wx.GetMethod[FormPostTest]("Post")
	assert.NotNil(t, method, "Method should not be nil")
	info, err := wx.Helper.GetHandlerInfo(*method)
	assert.NoError(t, err, "Error should be nil")
	assert.NotNil(t, info)
	assert.Equal(t, 2, info.IndexOfRequestBody)
	assert.Equal(t, true, info.IsFormUpload)

}
