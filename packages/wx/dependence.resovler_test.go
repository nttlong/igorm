package wx

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct {
}
type B struct {
}

func (b *B) New() error {
	return nil
}
func (a *A) New(b *B) error {
	return nil
}
func TestResolveType(t *testing.T) {
	r, err := DepenResolvers.ResolveType(reflect.TypeFor[A]())
	assert.Nil(t, err)
	assert.NotNil(t, r)
}
func BenchmarkResolveType(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r, err := DepenResolvers.ResolveType(reflect.TypeFor[A]())
		assert.Nil(b, err)
		assert.NotNil(b, r)
	}

}
func BenchmarkResolveTypeOnce(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r, err := DepenResolvers.ResolveTypeOnce(reflect.TypeFor[A]())
		assert.Nil(b, err)
		assert.NotNil(b, r)
	}

}
