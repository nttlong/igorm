package vgrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestService struct {
	Start int
}

func (t *TestService) AddMore(data int) {
	t.Start += data
}

func TestAddSvc(t *testing.T) {
	AddSingletonService(func() (*TestService, error) {
		return &TestService{
			Start: 0,
		}, nil
	})
	ret, err := Call("vgrpc.TestService.AddMore", []interface{}{})
	assert.NoError(t, err)

	assert.Equal(t, len(ret), 0)
}
func BenchmarkTestAddSvc(b *testing.B) {
	AddSingletonService(func() (*TestService, error) {
		return &TestService{
			Start: 0,
		}, nil
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		ret, err := Call("vgrpc.TestService.Hello", []interface{}{})
		assert.NoError(b, err)

		assert.Equal(b, len(ret), 0)
	}
}
func BenchmarkDirectCall(b *testing.B) {

	b.ResetTimer()
	b.Run("DirectCall", func(b *testing.B) {
		svc := &TestService{}
		for i := 0; i < b.N; i++ {
			svc.AddMore(i)
		}

	})
	b.Run("ReflectCall", func(b *testing.B) {
		AddSingletonService(func() (*TestService, error) {
			return &TestService{}, nil
		})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ret, err := Call("vgrpc.TestService.AddMore", []interface{}{i})
			assert.NoError(b, err)

			assert.Equal(b, len(ret), 0)
		}

	})
	b.Run("ReflectCall2", func(b *testing.B) {
		AddSingletonService(func() (*TestService, error) {
			return &TestService{}, nil
		})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ret, err := Call2("vgrpc.TestService.AddMore", []interface{}{i})
			assert.NoError(b, err)

			assert.Equal(b, len(ret), 0)
		}

	})

}
