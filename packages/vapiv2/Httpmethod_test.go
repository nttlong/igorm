package vapi

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpMethod(t *testing.T) {
	ret := httpUtilsTypeInstance.InspectHttpMethodFromType(reflect.TypeFor[*struct {
		HttpGet
	}]())
	assert.Equal(t, "get", ret)

}
func TestGetTags(t *testing.T) {
	ret := httpUtilsTypeInstance.InspectHttpMethodFromType(reflect.TypeFor[*struct {
		HttpPost `route:"/test"`
	}]())
	assert.Equal(t, "post", ret)

}
