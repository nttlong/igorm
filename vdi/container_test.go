package vdi_test

import (
	"reflect"
	"testing"
	"unsafe"
	"vdi"

	"github.com/stretchr/testify/assert"
)

type ConfigService struct {
	path string
}

func NewConfigService(path string) *ConfigService {
	return &ConfigService{path: path}
}

type AppContainer struct {
	Config vdi.Singleton[AppContainer, *ConfigService]
}

func TestAppContainer(t *testing.T) {
	for i := 0; i < 10; i++ {
		app, err := vdi.RegisterContainer(func(svc *AppContainer) error {
			svc.Config.Init = func(owner *AppContainer) *ConfigService {
				return NewConfigService("config.json")
			}
			return nil

		})
		a := app.Get()
		f1 := a.Config.Get()
		f2 := a.Config.Get()
		assert.Same(t, f1, f2)

		assert.NoError(t, err)
		assert.NotEmpty(t, app)
	}

}
func BenchmarkTestAppContainer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		app, err := vdi.RegisterContainer(func(svc *AppContainer) error {
			svc.Config.Init = func(owner *AppContainer) *ConfigService {
				return NewConfigService("config.json")
			}
			return nil

		})
		a := app.Get()
		f1 := a.Config.Get()
		f2 := a.Config.Get()
		assert.Same(b, f1, f2)

		assert.NoError(b, err)
		assert.NotEmpty(b, app)
	}

}

// A: Reuse singleton trong cùng container
func BenchmarkTestAppContainer_ReuseSingleton(b *testing.B) {
	type ConfigService struct {
		Path string
	}
	type AppContainer struct {
		Config vdi.Singleton[AppContainer, *ConfigService]
	}
	app, err := vdi.RegisterContainer(func(svc *AppContainer) error {
		svc.Config.Init = func(owner *AppContainer) *ConfigService {
			return &ConfigService{Path: "config.json"}
		}
		return nil
	})
	assert.NoError(b, err)
	container := app.Get()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cfg := container.Config.Get()
		_ = cfg
	}
}

// B: Transient — khởi tạo mới mỗi lần gọi
func BenchmarkTestTransient(b *testing.B) {
	type TransientService struct {
		ID int
	}
	factory := func() *TransientService {
		return &TransientService{ID: 123}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = factory()
	}
}

// C: Nested singleton (ServiceA → ServiceB)
func BenchmarkTestNestedSingleton(b *testing.B) {
	type ServiceB struct {
		Value string
	}
	type ServiceA struct {
		B *ServiceB
	}
	type AppContainer struct {
		B    vdi.Singleton[AppContainer, *ServiceB]
		A    vdi.Singleton[AppContainer, *ServiceA]
		Name string
	}
	app, err := vdi.RegisterContainer(func(svc *AppContainer) error {
		svc.B.Init = func(owner *AppContainer) *ServiceB {
			return &ServiceB{Value: "Hello"}
		}
		svc.A.Init = func(owner *AppContainer) *ServiceA {
			return &ServiceA{B: owner.B.Get()}
		}
		return nil
	})
	o1 := app.Get().A.Owner
	o2 := app.Get().B.Owner
	o3 := app.Get()
	o3.Name = "test"

	va := reflect.ValueOf(o1)
	vb := reflect.ValueOf(o2)

	// Lấy pointer bên trong nếu là pointer
	if va.Kind() == reflect.Ptr {
		va = va.Elem()
	}
	if vb.Kind() == reflect.Ptr {
		vb = vb.Elem()
	}
	ptrA := unsafe.Pointer(va.UnsafeAddr())
	ptrB := unsafe.Pointer(vb.UnsafeAddr())
	print(ptrA == ptrB)

	assert.NoError(b, err)

	c := app.Get()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		a := c.A.Get()
		_ = a.B.Value
	}
}
