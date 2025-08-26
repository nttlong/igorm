package mockcontroller

import (
	"testing"
	"wx"

	"wx/handlers"
	_ "wx/handlers"

	"github.com/stretchr/testify/assert"
)

type Service1 struct{}
type Service2 struct{}

type Controller1 struct {
	Svc1 *wx.Depend[Service1]
	svc2 *wx.Global[Service2]
	fx   Service2
}

func (c *Controller1) New(svc2 *wx.Global[Service2]) error {
	c.svc2 = svc2
	return nil

}
func (c *Controller1) Post(ctx *wx.Handler) (interface{}, error) {
	t, err := c.svc2.Ins()
	if err != nil {
		return nil, err
	}
	c.fx = t

	return &struct {
		Name string
		Age  int
	}{
		Name: "John",
		Age:  30,
	}, nil
}

func TestCallHandlerPost(t *testing.T) {
	mt := wx.GetMethodByName[Controller1]("Post")

	mtInfo, err := handlers.Helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuild.PostJson("/api/"+mtInfo.UriHandler, nil)
	req, res := requestBuild.Build()

	ret, err := wx.Helper.ReqExec.DoJsonPost(*mtInfo, req, res)
	assert.NoError(t, err)
	t.Log(ret)

	t.Log(mtInfo)
}
func BenchmarkTestCallHandlerPost(b *testing.B) {
	mt := wx.GetMethodByName[Controller1]("Post")
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
func (c *Controller1) PostWithBody(ctx *wx.Handler, body *struct {
	Name string
	Age  int
}) (interface{}, error) {
	return body, nil
}
func TestCallHandlerPostWithBody(t *testing.T) {
	mt := wx.GetMethodByName[Controller1]("PostWithBody")

	mtInfo, err := handlers.Helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuild.PostJson("/api/"+mtInfo.UriHandler, nil)
	req, res := requestBuild.Build()

	ret, err := wx.Helper.ReqExec.DoJsonPost(*mtInfo, req, res)
	assert.NoError(t, err)
	assert.Empty(t, ret)
	requestBuild.PostJson("/api/"+mtInfo.UriHandler, &struct {
		Name string
		Age  int
	}{
		Name: "John",
		Age:  30,
	})
	req, res = requestBuild.Build()
	ret2, err := wx.Helper.ReqExec.DoJsonPost(*mtInfo, req, res)

	assert.NotEmpty(t, ret2)
}
func BenchmarkTestCallHandlerPostWithBody(b *testing.B) {
	mt := wx.GetMethodByName[Controller1]("PostWithBody")

	mtInfo, _ := handlers.Helper.GetHandlerInfo(*mt)

	b.Run("ParallePostEmptyBody", func(b *testing.B) {
		requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()
		requestBuild.PostJson("/api/"+mtInfo.UriHandler, nil)
		req, res := requestBuild.Build()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				wx.Helper.ReqExec.DoJsonPost(*mtInfo, req, res)
			}
		})

	})
	b.Run("ParallePostNotEmptyBody", func(b *testing.B) {
		requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()
		requestBuild.PostJson("/api/"+mtInfo.UriHandler, &struct {
			Name string
			Age  int
		}{
			Name: "John",
			Age:  30,
		})
		req, res := requestBuild.Build()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				wx.Helper.ReqExec.DoJsonPost(*mtInfo, req, res)
			}
		})

	})
	b.Run("PostEmptyBody", func(b *testing.B) {
		requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()
		requestBuild.PostJson("/api/"+mtInfo.UriHandler, nil)
		req, res := requestBuild.Build()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wx.Helper.ReqExec.DoJsonPost(*mtInfo, req, res)
		}

	})
	b.Run("PostNotEmptyBody", func(b *testing.B) {
		requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()
		requestBuild.PostJson("/api/"+mtInfo.UriHandler, &struct {
			Name string
			Age  int
		}{
			Name: "John",
			Age:  30,
		})
		req, res := requestBuild.Build()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wx.Helper.ReqExec.DoJsonPost(*mtInfo, req, res)
		}

	})
}
func (c *Controller1) PostWithRequireBody(ctx *wx.Handler, body struct {
	Name string
	Age  int
}) (interface{}, error) {
	return body, nil
}
func TestCallHandlerPostWithRequireBody(t *testing.T) {
	mt := wx.GetMethodByName[Controller1]("PostWithRequireBody")

	mtInfo, err := handlers.Helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuild.PostJson("/api/"+mtInfo.UriHandler, nil)
	req, res := requestBuild.Build()

	ret, err := wx.Helper.ReqExec.DoJsonPost(*mtInfo, req, res)
	assert.Error(t, err)
	assert.Empty(t, ret)
	requestBuild.PostJson("/api/"+mtInfo.UriHandler, &struct {
		Name string
		Age  int
	}{
		Name: "John",
		Age:  30,
	})
	req, res = requestBuild.Build()
	ret2, err := wx.Helper.ReqExec.DoJsonPost(*mtInfo, req, res)

	assert.NotEmpty(t, ret2)
}
func BenchmarkPostWithRequireBody(b *testing.B) {
	mt := wx.GetMethodByName[Controller1]("PostWithRequireBody")

	mtInfo, _ := handlers.Helper.GetHandlerInfo(*mt)

	requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()
	requestBuild.PostJson("/api/"+mtInfo.UriHandler, nil)

	b.Run("PtrToStruct", func(b *testing.B) {

		req, res := requestBuild.Build()
		requestBuild.PostJson("/api/"+mtInfo.UriHandler, &struct {
			Name string
			Age  int
		}{
			Name: "John",
			Age:  30,
		})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {

			wx.Helper.ReqExec.DoJsonPost(*mtInfo, req, res)
		}
	})
	b.Run("PtrToNil", func(b *testing.B) {
		req, res := requestBuild.Build()
		requestBuild.PostJson("/api/"+mtInfo.UriHandler, nil)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			wx.Helper.ReqExec.DoJsonPost(*mtInfo, req, res)
		}

	})

}
