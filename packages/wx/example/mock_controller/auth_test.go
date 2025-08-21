package mockcontroller

import (
	"net/http"
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

type AuthTest struct {
	User wx.UserClaims
}

func (a *AuthTest) New() {

}
func (a *AuthTest) Post(ctx *wx.Handler, user wx.UserClaims) (interface{}, error) {
	return nil, nil
}
func TestAuthTestPost(t *testing.T) {
	mt := wx.GetMethodByName[AuthTest]("Post")
	mtInfo, err := wx.Helper.GetHandlerInfo(*mt)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 0, mtInfo.IndexOfAuthClaimsArg)
	assert.Equal(t, []int{0, 0}, mtInfo.IndexOfAuthClaims)
	requestBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuilder.PostJson("/api/"+mtInfo.UriHandler, nil)
	requestBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {
		wx.Helper.ReqExec.Invoke(*mtInfo, r, w)
	})

}
