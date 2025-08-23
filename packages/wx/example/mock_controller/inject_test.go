package mockcontroller

import (
	"net/http"
	"testing"
	"wx"
	"wx/libs"
	_ "wx/libs"

	"github.com/stretchr/testify/assert"
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

type DbService struct {
}

func (c *ControllerInject) GetUser(ctx *struct {
	wx.Handler `route:"@/get-user;method:get"`
}, db *wx.Depend[DbService]) (interface{}, error) {
	return nil, nil

}
func TestGetUser(t *testing.T) {
	mt := wx.GetMethodByName[ControllerInject]("GetUser")
	mtInfo, err := wx.Helper.GetHandlerInfo(*mt)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(mtInfo.IndexOfInjectors))
	requestBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuilder.Get("/api/" + mtInfo.UriHandler)
	for i := 0; i < 5; i++ {
		requestBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {
			wx.Helper.ReqExec.Invoke(*mtInfo, r, w)

		})
	}

}

type IService interface {
	GetName() string
}
type Resovler[T any, TImplement any] struct {
	Implementation TImplement
}

func (r *Resovler[T, TImplement]) Resolve() (T, error) {
	var t T
	return t, nil
}

type Service struct {
}

func (s *Service) New() error {
	// This method is used to initialize the Service instance
	return nil
}
func (s *Service) GetName() string {
	return "Service"
}

func (c *ControllerInject) GetUser2(ctx *struct {
	wx.Handler `route:"@/get-user;method:get"`
}, service1 *wx.Depend[Service]) (interface{}, error) {

	mySvc, err := service1.Ins() // get instance of Service
	if err != nil {
		return nil, err
	}

	return mySvc.GetName(), nil

}
func TestGetUser2(t *testing.T) {
	mt := wx.GetMethodByName[ControllerInject]("GetUser2")
	mtInfo, err := wx.Helper.GetHandlerInfo(*mt)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(mtInfo.IndexOfInjectors))
	requestBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuilder.Get("/api/" + mtInfo.UriHandler)
	for i := 0; i < 5; i++ {
		requestBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {
			wx.Helper.ReqExec.Invoke(*mtInfo, r, w)

		})
	}

}

type UserService struct {
	repo UserRepository // inject UserRepository dependency
}
type UserRepository struct {
}

func (u *UserService) GetUserById() string {
	return "UserService"
}
func (c *ControllerInject) GetUser3(ctx *struct {
	wx.Handler `route:"@/get-user/{UserId};method:get"`
	UserId     string
}, users *wx.Depend[UserService]) (interface{}, error) {

	uSvc, err := users.Ins() // get instance of UserService
	if err != nil {
		return nil, err
	}

	return uSvc.GetUserById(), nil

}
func TestGetGetUser3(t *testing.T) {
	mt := wx.GetMethodByName[ControllerInject]("GetUser3")
	mtInfo, err := wx.Helper.GetHandlerInfo(*mt)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 1, len(mtInfo.IndexOfInjectors))
}

func (u *UserService) New(context wx.HttpContext, Repo *wx.Depend[UserRepository]) error { // auto called by wx.Helper when the UserService is injected
	// This method is used to initialize the UserService instance
	userRepo, err := Repo.Ins() // get instance of UserRepository
	if err != nil {
		return err
	}
	u.repo = *userRepo // assign the UserRepository instance to the UserService
	return nil
}

func (c *ControllerInject) GetUserWithInjectService(ctx *struct {
	wx.Handler `route:"@/get-user/{UserId};method:post"`
	UserId     string
}, users *wx.HttpService[UserService]) (interface{}, error) {

	uSvc, err := users.Ins() // get instance of UserService
	if err != nil {
		return nil, err
	}

	return uSvc.GetUserById(), nil

}
func TestGetGetUserWithInjectService(t *testing.T) {
	mt := wx.GetMethodByName[ControllerInject]("GetUserWithInjectService")
	mtInfo, err := wx.Helper.GetHandlerInfo(*mt)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, -1, mtInfo.IndexOfRequestBody)
	assert.Equal(t, 1, len(mtInfo.IndexOfInjectorService))
	requestBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuilder.PostJson("/api/"+mtInfo.UriHandler, nil)
	for i := 0; i < 5; i++ {
		requestBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {
			wx.Helper.ReqExec.Invoke(*mtInfo, r, w)

		})
	}

}
func BenchmarkTestGetGetUserWithInjectService(b *testing.B) {
	mt := wx.GetMethodByName[ControllerInject]("GetUserWithInjectService")
	mtInfo, _ := wx.Helper.GetHandlerInfo(*mt)

	requestBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuilder.PostJson("/api/"+mtInfo.UriHandler, nil)
	for i := 0; i < b.N; i++ {
		requestBuilder.Handler(func(w http.ResponseWriter, r *http.Request) {
			wx.Helper.ReqExec.Invoke(*mtInfo, r, w)

		})
	}
}
