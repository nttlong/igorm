package vapi

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnceCall(t *testing.T) {
	mt := GetMethodByName[User]("Get2")
	for i := 0; i < 10; i++ {
		ret, err := OnceCall(mt.Type.In(2).PkgPath()+"/helperType/IsInjector", func() (*bool, error) {
			typ := mt.Type.In(2)
			fmt.Println(typ.String())
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			if typ.PkgPath() == reflect.TypeOf(Inject[string]{}).PkgPath() {
				if strings.Contains(typ.Name(), "Inject[") {
					ret := true
					return &ret, nil
				}

			}

			ret := false
			return &ret, nil
		})
		assert.Nil(t, err)
		assert.Equal(t, true, *ret)
	}
}
