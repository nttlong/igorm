package mockcontroller

import (
	"net/http"
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

type UserInfo struct {
	Username string `json:"username"`
	UserId   string `json:"userId"`
}

func (u *UserInfo) New(ctx *wx.AuthContext) error {
	return nil

}
func TestUserContext(t *testing.T) {
	// Create a new UserContext

}

type AuthTest struct {
	User wx.UserClaims
}

func (a *AuthTest) New() {
	// fmt.Println("New method called")

}

var handlerEror error

func (a *AuthTest) Post(ctx *wx.Handler, data *UserInfo, user *wx.Auth[UserInfo]) (interface{}, error) {
	// var db DbService
	_, err := user.Get()
	if err != nil {
		handlerEror = err
		return nil, err
	}
	// fmt.Println("User service:", userSvc)

	return nil, nil
}
func TestAuthTestPost(t *testing.T) {
	mt := wx.GetMethodByName[AuthTest]("Post")
	for i := 1; i < mt.Type.NumIn(); i++ {
		typ := mt.Type.In(i)
		_, err := wx.Helper.DepenAuthFind(typ)
		assert.NoError(t, err)

	}
	mtInfo, err := wx.Helper.GetHandlerInfo(*mt)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 3, mtInfo.IndexOfAuthClaimsArg)
	assert.Equal(t, []int{}, mtInfo.IndexOfAuthClaims)
	assert.Equal(t, 2, mtInfo.IndexOfRequestBody)
	assert.Equal(t, 1, mtInfo.IndexOfArg)
	assert.Equal(t, []int{}, mtInfo.IndexOfInjectors)

	requestBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuilder.PostJson("/api/"+mtInfo.UriHandler, nil)

	requestBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {
		wx.Helper.ReqExec.Invoke(*mtInfo, r, w)
	})

}
func BenchmarkAuthTestPost(b *testing.B) {
	mt := wx.GetMethodByName[AuthTest]("Post")

	mtInfo, _ := wx.Helper.GetHandlerInfo(*mt)

	requestBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuilder.PostJson("/api/"+mtInfo.UriHandler, nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		requestBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {
			wx.Helper.ReqExec.Invoke(*mtInfo, r, w)
		})
	}

}
