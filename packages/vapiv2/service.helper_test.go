package vapi

import (
	"net/http"
	httptest "net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FileUtils struct {
}
type S3Utils struct {
}
type BaseService struct {
	S3 Scoped[S3Utils]
}

type Service1 struct {
	BaseService
	Files Singleton[FileUtils]
	//S3    *Scoped[S3Utils]
}

func (s *Service1) New() error {
	s.Files.Init(func() (*FileUtils, error) {
		return &FileUtils{}, nil
	})
	s.S3.Init(func(ctx *ServiceContext) (*S3Utils, error) {
		return &S3Utils{}, nil
	})
	// s.S3.Init(func(ctx *ServiceContext) (*S3Utils, error) {
	// 	return &S3Utils{}, nil
	// })
	return nil
}
func createMockServiceContext() ServiceContext {
	ctx := ServiceContext{
		Req: httptest.NewRequest(
			http.MethodPost,
			"/api/data?foo=bar",
			strings.NewReader(`{"name": "John"}`),
		),
		Res: httptest.NewRecorder(),
	}
	return ctx
}
func TestCreateServiceContext(t *testing.T) {
	ctx := &ServiceContext{
		Req: httptest.NewRequest(
			http.MethodPost,
			"/api/data?foo=bar",
			strings.NewReader(`{"name": "John"}`),
		),
		Res: httptest.NewRecorder(),
	}
	t.Log(ctx.Req)
}
func TestInjectorNew(t *testing.T) {
	for i := 0; i < 5; i++ {
		typ := reflect.TypeFor[Service1]()
		ctx := createMockServiceContext()
		svcVal, err := serviceUtils.NewService(typ, ctx.Req, ctx.Res)
		assert.NoError(t, err)
		assert.NotNil(t, svcVal)
	}
}
func BenchmarkInjectorNew(b *testing.B) {
	typ := reflect.TypeFor[*Service1]()
	for i := 0; i < b.N; i++ {
		svcVal, err := serviceUtils.NewServiceOptimize(typ)
		assert.NoError(b, err)
		assert.NotNil(b, svcVal)
	}

}
func BenchmarkNewService(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := createMockServiceContext()
		svcVal, err := serviceUtils.NewService(reflect.TypeFor[*Service1](), ctx.Req, ctx.Res)
		assert.NoError(b, err)
		assert.NotNil(b, svcVal)
	}

}
func BenchmarkNewServiceOptimize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		svcVal, err := serviceUtils.NewServiceOptimize(reflect.TypeFor[*Service1]())
		assert.NoError(b, err)
		assert.NotNil(b, svcVal)
	}

}
