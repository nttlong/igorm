package vapi

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSingletonInstance(t *testing.T) {
	typ := reflect.TypeOf(&S3Utils{})
	val := serviceUtils.CreateSingletonInstance(typ)
	obj := val.Elem().Interface().(S3Utils)
	assert.NotNil(t, obj)
	vale := serviceUtils.CreateSingletonInstance(typ)
	obj1 := vale.Elem().Interface().(S3Utils)
	assert.NotNil(t, obj1)
	assert.Equal(t, obj, obj1)

}
func BenchmarkTestCreateSingletonInstance(t *testing.B) {
	for i := 0; i < t.N; i++ {
		typ := reflect.TypeOf(&S3Utils{})
		serviceUtils.CreateSingletonInstance(typ)
		// obj := val.Elem().Interface().(S3Utils)
		// assert.NotEmpty(t, obj)
	}

}
