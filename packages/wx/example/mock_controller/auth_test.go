package mockcontroller

import (
	"fmt"
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
	fmt.Println("New method called")

}
func (a *AuthTest) Post(ctx *wx.Handler, data *UserInfo, user *wx.Auth[UserInfo]) (interface{}, error) {
	userSvc, err := user.Get()
	if err != nil {
		return nil, err
	}
	fmt.Println("User service:", userSvc)

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
	// req, res := requestBuilder.Build()
	// newMt := wx.GetMethodByName[UserInfo]("New")
	// valOfCtx := reflect.ValueOf(&wx.AuthContext{
	// 	Req: req,
	// 	Res: res,
	// })
	// v, err := wx.Helper.DepenAuthCreateInstance(*newMt, valOfCtx)
	// assert.NoError(t, err)
	// assert.NotNil(t, v)
	// ele := v.Interface().(*UserInfo)

	// assert.NotNil(t, ele)
	// authVal, err := wx.Helper.DepenAuthCreate(reflect.TypeOf(&wx.Auth[UserInfo]{}))
	// assert.NoError(t, err)
	// assert.NotNil(t, authVal)
	// auth := authVal.Interface().(*wx.Auth[UserInfo])
	// val, err := auth.Get()
	// assert.NoError(t, err)
	// assert.NotNil(t, val)

	requestBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {
		wx.Helper.ReqExec.Invoke(*mtInfo, r, w)
	})

}
