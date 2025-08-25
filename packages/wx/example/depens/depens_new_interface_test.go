package depens

import (
	"reflect"
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

type IService1 interface {
}
type Service1 struct {
}

func (svc *Service1) New() (IService1, error) {
	return svc, nil
}

type IService interface {
	GetData() any
}
type TestService struct {
	Svc1 IService1
}

func (svc *TestService) New(svc1 *wx.Depend[Service1]) (IService, error) {
	val, err := svc1.Ins()
	if err != nil {
		return nil, err
	}
	svc.Svc1 = val

	return svc, nil
}
func (svc *TestService) GetData() any {
	return svc

}
func TestNewDependTestService(t *testing.T) {
	val, err := wx.Helper.DependNew(reflect.TypeOf(&TestService{}))
	assert.NoError(t, err)
	assert.NotNil(t, val)

}
func BenchmarkNewDependTestService(b *testing.B) {
	for i := 0; i < b.N; i++ {
		val, err := wx.Helper.DependNew(reflect.TypeOf(&TestService{}))
		assert.NoError(b, err)
		assert.NotNil(b, val)
	}

}
