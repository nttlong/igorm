package mockcontroller

import (
	"net/http"
	"testing"
	"wx"
	"wx/libs"
	_ "wx/libs"
)

type ControllerInject struct {
	PassSvc wx.Depend[libs.PasswordService]
}

func (c *ControllerInject) New() {
	c.PassSvc.Init(func() (*libs.PasswordService, error) {
		return &libs.PasswordService{}, nil
	})

}

type CreateUserData struct {
	Name     string
	Password string
}

func (c *ControllerInject) CreateUser(ctx *wx.Handler, data CreateUserData) (interface{}, error) {
	passSvc, err := c.PassSvc.Ins()
	if err != nil {
		return nil, err
	}
	txt, err := passSvc.HashPassword(data.Password + "@" + data.Name)
	if err != nil {
		return nil, err
	}

	return struct {
		Text string
	}{
		Text: txt,
	}, nil
}
func TestCreateUser(t *testing.T) {
	mt := wx.GetMethodByName[ControllerInject]("CreateUser")
	mtInfo, err := wx.Helper.GetHandlerInfo(*mt)
	if err != nil {
		t.Error(err)
	}
	requestBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuilder.PostJson("/api/"+mtInfo.UriHandler, CreateUserData{
		Name:     "John",
		Password: "123456",
	})
	for i := 0; i < 5; i++ {
		requestBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {
			wx.Helper.ReqExec.Invoke(*mtInfo, r, w)

		})
	}

}
func BenchmarkCreateUser(b *testing.B) {
	mt := wx.GetMethodByName[ControllerInject]("CreateUser")
	mtInfo, _ := wx.Helper.GetHandlerInfo(*mt)

	requestBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuilder.PostJson("/api/"+mtInfo.UriHandler, CreateUserData{
		Name:     "John",
		Password: "123456",
	})
	b.Run("test", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			requestBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {
				wx.Helper.ReqExec.Invoke(*mtInfo, r, w)

			})
		}
	})

}
func (c *ControllerInject) GetUser(ctx *struct {
	wx.Handler `route:"@/get-user;method:get"`
}) (interface{}, error) {
	return nil, nil

}
func TestGetUser(t *testing.T) {
	mt := wx.GetMethodByName[ControllerInject]("GetUser")
	mtInfo, err := wx.Helper.GetHandlerInfo(*mt)
	if err != nil {
		t.Error(err)
	}
	requestBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuilder.Get("/api/" + mtInfo.UriHandler)
	for i := 0; i < 5; i++ {
		requestBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {
			wx.Helper.ReqExec.Invoke(*mtInfo, r, w)

		})
	}

}
