package vdi_test

import (
	"fmt"
	"testing"
	"vdi"

	unvsdi "vdi"

	"github.com/stretchr/testify/assert"
)

type Logger struct {
	ID string
}

type Service struct {
	Logger *Logger `inject:""`
}

func TestScoped_ResolveDifferentInstances(t *testing.T) {
	root := unvsdi.NewRootContainer()
	root.RegisterScoped(func() *Logger {
		return &Logger{ID: "scoped-logger"}
	})

	scope1 := root.CreateScope()
	scope2 := root.CreateScope()

	l1, _ := scope1.ResolveByType(unvsdi.TypeOf[*Logger]())
	l2, _ := scope1.ResolveByType(unvsdi.TypeOf[*Logger]())
	l3, _ := scope2.ResolveByType(unvsdi.TypeOf[*Logger]())

	assert.Same(t, l1, l2)    // ✅ l1 và l2 trong cùng scope, phải giống nhau
	assert.NotSame(t, l1, l3) // ✅ l3 thuộc scope khác, phải khác
}

func TestSingleton_SameInstance(t *testing.T) {
	root := unvsdi.NewRootContainer()
	root.RegisterSingleton(func() *Logger {
		return &Logger{ID: "singleton"}
	})

	scope1 := root.CreateScope()
	scope2 := root.CreateScope()

	l1, _ := scope1.ResolveByType(unvsdi.TypeOf[*Logger]())
	l2, _ := scope1.ResolveByType(unvsdi.TypeOf[*Logger]())
	l3, _ := scope2.ResolveByType(unvsdi.TypeOf[*Logger]())

	assert.Same(t, l1, l2)
	assert.Same(t, l1, l3)
}

func TestTransient_AlwaysNew(t *testing.T) {
	root := unvsdi.NewRootContainer()
	root.RegisterTransient(func() *Logger {
		return &Logger{ID: "transient"}
	})

	scope := root.CreateScope() //<-- cho nay phai la CreateTransient chứ

	l1, _ := scope.ResolveByType(unvsdi.TypeOf[*Logger]())
	l2, _ := scope.ResolveByType(unvsdi.TypeOf[*Logger]())

	assert.NotSame(t, l1, l2)
}

func TestInjectFields(t *testing.T) {
	root := unvsdi.NewRootContainer()
	root.RegisterSingleton(func() *Logger {
		return &Logger{ID: "log-1"}
	})

	scope := root.CreateScope()
	svc := &Service{}
	err := unvsdi.InjectFields(scope, svc)
	assert.NoError(t, err)
	assert.Equal(t, "log-1", svc.Logger.ID)
}

func TestResolveUnknownType(t *testing.T) {
	root := unvsdi.NewRootContainer()
	scope := root.CreateScope()

	_, err := scope.ResolveByType(unvsdi.TypeOf[*Logger]())
	assert.Error(t, err)
}
func TestDI_ResolveSingleton(t *testing.T) {
	container := vdi.NewRootContainer()

	container.RegisterSingleton(func() *Logger {
		return &Logger{ID: "singleton"}
	})

	val, err := container.ResolveByType(vdi.TypeOf[*Logger]())
	assert.NoError(t, err)

	logger1 := val.(*Logger)
	assert.Equal(t, "singleton", logger1.ID)
}

func TestDI_ResolveTransient(t *testing.T) {
	container := unvsdi.NewRootContainer()
	container.RegisterTransient(func() *Logger {
		return &Logger{ID: "transient"}
	})

	val1, err := container.ResolveByType(unvsdi.TypeOf[*Logger]())
	assert.NoError(t, err)

	val2, err := container.ResolveByType(unvsdi.TypeOf[*Logger]())
	assert.NoError(t, err)
	if val1 != val2 {
		fmt.Println("val1:", val1, "val2:", val2) //<-- no da chay qua cho nay
	}

	assert.NotSame(t, val1, val2)
	//<-- nhung cho nay lai assert.NotEqual lai bao loi Error:      	Should not be: &vdi_test.Logger{ID:"transient"}
}

func TestDI_ScopedOverride(t *testing.T) {
	root := unvsdi.NewRootContainer()
	root.RegisterScoped(func() *Logger {
		return &Logger{ID: "root-scope"}
	})

	scope1 := root.CreateScope()
	scope2 := root.CreateScope()

	l1, err := scope1.ResolveByType(unvsdi.TypeOf[*Logger]()) //<-- khi goi cai nay
	assert.NoError(t, err)
	l2, err := scope1.ResolveByType(unvsdi.TypeOf[*Logger]())
	assert.NoError(t, err)
	l3, err := scope2.ResolveByType(unvsdi.TypeOf[*Logger]())
	assert.NoError(t, err)
	assert.NotSame(t, l1, l2)
	assert.NotSame(t, l1, l3)
}

func TestDI_InjectFields(t *testing.T) {
	container := unvsdi.NewRootContainer()
	container.RegisterSingleton(func() *Logger {
		return &Logger{ID: "injected"}
	})

	svc := &Service{}
	err := unvsdi.InjectFields(container, svc)
	assert.NoError(t, err)
	assert.Equal(t, "injected", svc.Logger.ID)
}
