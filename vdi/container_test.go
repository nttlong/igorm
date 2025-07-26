package vdi_test

import (
	"crypto/rand"
	"fmt"
	"math/big"
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

// C: Nested singleton (ServiceA → ServiceB)
func BenchmarkTestNestedSingleton(b *testing.B) {
	type ServiceB struct {
		Value string
	}
	type ServiceA struct {
		B *ServiceB
	}
	type ServiceC struct {
		B *ServiceB
	}
	type AppContainer struct {
		B vdi.Singleton[AppContainer, *ServiceB]
		A vdi.Singleton[AppContainer, *ServiceA]
		C vdi.Transient[AppContainer, *ServiceC]

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
	o4 := app.Get().C.Owner
	fmt.Println(o4)

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
func TestTransition(t *testing.T) {
	type ServiceB struct {
		Value string
	}
	type ServiceA struct {
		B *ServiceB
	}
	type ServiceC struct {
		B *ServiceB
	}
	type AppContainer struct {
		B vdi.Transient[AppContainer, *ServiceB]
		A vdi.Transient[AppContainer, *ServiceA]

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
	assert.NoError(t, err)
	c := app.Get()
	a := c.A.Get()
	b := c.B.Get()
	assert.Equal(t, "Hello", b.Value)
	assert.Equal(t, "Hello", a.B.Value)
	c.Name = "test"
	assert.Equal(t, "test", c.Name)
}
func BenchmarkTestTransientA(b *testing.B) {
	type Logger struct {
		Owner interface{}
		Value int
	}

	type AppContainer struct {
		Log vdi.Transient[AppContainer, *Logger]
	}

	app, _ := vdi.RegisterContainer(func(c *AppContainer) error {
		c.Log.Owner = c
		c.Log.Init = func(owner *AppContainer) *Logger {
			return &Logger{Value: 123}
		}
		return nil
	})

	c := app.Get()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log := c.Log.Get()
		_ = log.Value
	}
}

func BenchmarkTestTransientB(b *testing.B) {
	type Logger struct {
		Value int
	}

	type AppContainer struct {
		Log vdi.Transient[AppContainer, *Logger]
	}

	app, _ := vdi.RegisterContainer(func(c *AppContainer) error {
		c.Log.Init = func(owner *AppContainer) *Logger {
			return &Logger{Value: 456}
		}
		return nil
	})

	c := app.Get()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log := c.Log.Get()
		_ = log.Value
	}
}
func BenchmarkTestTransientC(b *testing.B) {
	type Logger struct {
		ID int
	}

	type Service struct {
		Log *Logger
	}

	type AppContainer struct {
		Log     vdi.Transient[AppContainer, *Logger]
		Service vdi.Transient[AppContainer, *Service]
	}

	app, _ := vdi.RegisterContainer(func(c *AppContainer) error {
		c.Log.Init = func(owner *AppContainer) *Logger {
			return &Logger{ID: 1}
		}
		c.Service.Init = func(owner *AppContainer) *Service {
			return &Service{
				Log: owner.Log.Get(),
			}
		}
		return nil
	})

	c := app.Get()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc := c.Service.Get()
		_ = svc.Log.ID
	}
}
func BenchmarkTestTransientD(b *testing.B) {
	type Config struct {
		Name string
	}

	type Logger struct {
		Config *Config
	}

	type Service struct {
		Log *Logger
	}

	type AppContainer struct {
		Config  vdi.Transient[AppContainer, *Config]
		Log     vdi.Transient[AppContainer, *Logger]
		Service vdi.Transient[AppContainer, *Service]
	}

	app, _ := vdi.RegisterContainer(func(c *AppContainer) error {
		c.Config.Init = func(owner *AppContainer) *Config {
			return &Config{Name: "default"}
		}
		c.Log.Init = func(owner *AppContainer) *Logger {
			return &Logger{Config: owner.Config.Get()}
		}
		c.Service.Init = func(owner *AppContainer) *Service {
			return &Service{Log: owner.Log.Get()}
		}
		return nil
	})

	c := app.Get()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc := c.Service.Get()
		_ = svc.Log.Config.Name
	}
}
func BenchmarkTestTransientE(b *testing.B) {
	type Config struct {
		Value string
	}
	type Logger struct {
		Config *Config
	}
	type Repository struct {
		Logger *Logger
	}
	type Service struct {
		Repo *Repository
	}
	type Handler struct {
		Service *Service
	}

	type AppContainer struct {
		Config     vdi.Transient[AppContainer, *Config]
		Logger     vdi.Transient[AppContainer, *Logger]
		Repository vdi.Transient[AppContainer, *Repository]
		Service    vdi.Transient[AppContainer, *Service]
		Handler    vdi.Transient[AppContainer, *Handler]
	}

	app, _ := vdi.RegisterContainer(func(c *AppContainer) error {
		c.Config.Init = func(owner *AppContainer) *Config {
			return &Config{Value: "deep"}
		}
		c.Logger.Init = func(owner *AppContainer) *Logger {
			return &Logger{Config: owner.Config.Get()}
		}
		c.Repository.Init = func(owner *AppContainer) *Repository {
			return &Repository{Logger: owner.Logger.Get()}
		}
		c.Service.Init = func(owner *AppContainer) *Service {
			return &Service{Repo: owner.Repository.Get()}
		}
		c.Handler.Init = func(owner *AppContainer) *Handler {
			return &Handler{Service: owner.Service.Get()}
		}
		return nil
	})

	c := app.Get()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := c.Handler.Get()
		_ = h.Service.Repo.Logger.Config.Value
	}
}
func BenchmarkMixedInjection(b *testing.B) {
	type Config struct {
		Name string
	}

	type Logger struct {
		Config *Config
	}

	type DbContext struct {
		ID big.Int
	}

	type Repository struct {
		Db  *DbContext
		Log *Logger
	}

	type Service struct {
		Repo *Repository
	}

	type AppContainer struct {
		Config  vdi.Singleton[AppContainer, *Config]
		Logger  vdi.Singleton[AppContainer, *Logger]
		Context vdi.Scoped[AppContainer, *DbContext]
		Repo    vdi.Transient[AppContainer, *Repository]
		Svc     vdi.Transient[AppContainer, *Service]
	}

	root, _ := vdi.RegisterContainer(func(c *AppContainer) error {
		c.Config.Init = func(owner *AppContainer) *Config {
			return &Config{Name: "app"}
		}
		c.Logger.Init = func(owner *AppContainer) *Logger {
			return &Logger{Config: owner.Config.Get()}
		}
		c.Context.Init = func(owner *AppContainer) *DbContext {
			val, _ := rand.Prime(rand.Reader, 128)
			return &DbContext{ID: *val}
		}
		c.Repo.Init = func(owner *AppContainer) *Repository {
			return &Repository{
				Db:  owner.Context.Get(),
				Log: owner.Logger.Get(),
			}
		}
		c.Svc.Init = func(owner *AppContainer) *Service {
			return &Service{
				Repo: owner.Repo.Get(),
			}
		}
		return nil
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Mỗi vòng lặp tạo 1 container mới => mô phỏng scope mới
		c := root.Get()

		svc := c.Svc.Get()
		_ = svc.Repo.Db.ID
		_ = svc.Repo.Log.Config.Name
	}
}
