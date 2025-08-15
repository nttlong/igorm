package vapi

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSingletonFields(t *testing.T) {
	ret := serviceUtils.GetSingletonFields(reflect.TypeFor[*Service1]())
	t.Log(ret)
	field := reflect.TypeOf(Service1{}).FieldByIndex(ret[0])
	assert.Equal(t, "Files", field.Name)

}
func BenchmarkGetSingletonFields(t *testing.B) {
	for i := 0; i < t.N; i++ {
		ret := serviceUtils.GetSingletonFields(reflect.TypeFor[*Service1]())

		reflect.TypeOf(Service1{}).FieldByIndex(ret[0])
		// assert.Equal(t, "Files", field.Name)
	}
}
