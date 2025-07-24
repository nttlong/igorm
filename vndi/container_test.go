package vdi_test

import (
	"testing"

	vdi "vndi"

	"github.com/stretchr/testify/assert"
)

type Logger struct {
	ID string
}

func BenchmarkTestScoped_ReuseWithinScope(b *testing.B) {
	root := vdi.NewRootContainer()
	root.RegisterScoped(func() *Logger {
		return &Logger{ID: "scoped"}
	})
	scope := root.CreateScope()
	t := vdi.TypeOf[*Logger]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = scope.ResolveByType(t)
	}
}

func BenchmarkTestTransient_AlwaysNew(b *testing.B) {
	root := vdi.NewRootContainer()
	root.RegisterTransient(func() *Logger {
		return &Logger{ID: "transient"}
	})
	scope := root.CreateScope()
	t := vdi.TypeOf[*Logger]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = scope.ResolveByType(t)
	}
}

func TestSingleton_SameEverywhere(t *testing.T) {
	root := vdi.NewRootContainer()
	root.RegisterSingleton(func() *Logger {
		return &Logger{ID: "singleton"}
	})

	scope1 := root.CreateScope()
	scope2 := root.CreateScope()

	l1, _ := scope1.ResolveByType(vdi.TypeOf[*Logger]())
	l2, _ := scope2.ResolveByType(vdi.TypeOf[*Logger]())
	l3, _ := root.ResolveByType(vdi.TypeOf[*Logger]())

	assert.Same(t, l1, l2)
	assert.Same(t, l1, l3)
}
func BenchmarkTestSingleton_SameEverywhere(t *testing.B) {
	for i := 0; i < t.N; i++ {
		root := vdi.NewRootContainer()
		root.RegisterSingleton(func() *Logger {
			return &Logger{ID: "singleton"}
		})

		scope1 := root.CreateScope()
		scope2 := root.CreateScope()

		l1, _ := scope1.ResolveByType(vdi.TypeOf[*Logger]())
		l2, _ := scope2.ResolveByType(vdi.TypeOf[*Logger]())
		l3, _ := root.ResolveByType(vdi.TypeOf[*Logger]())

		assert.Same(t, l1, l2)
		assert.Same(t, l1, l3)
	}
}
func TestScoped_UniquePerScope(t *testing.T) {
	root := vdi.NewRootContainer()
	root.RegisterScoped(func() *Logger {
		return &Logger{ID: "scoped"}
	})

	scope1 := root.CreateScope()
	scope2 := root.CreateScope()

	l1, _ := scope1.ResolveByType(vdi.TypeOf[*Logger]())
	l2, _ := scope1.ResolveByType(vdi.TypeOf[*Logger]()) // same as l1
	l3, _ := scope2.ResolveByType(vdi.TypeOf[*Logger]()) // different

	assert.Same(t, l1, l2)
	assert.NotSame(t, l1, l3)
}

func TestTransient_AlwaysNew(t *testing.T) {
	root := vdi.NewRootContainer()
	root.RegisterTransient(func() *Logger {
		return &Logger{ID: "transient"}
	})

	scope := root.CreateScope()

	l1, _ := scope.ResolveByType(vdi.TypeOf[*Logger]())
	l2, _ := scope.ResolveByType(vdi.TypeOf[*Logger]())

	assert.NotSame(t, l1, l2)
}
