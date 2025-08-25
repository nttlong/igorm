package wx

import (
	"reflect"
	"strings"
	"wx/handlers"
)

type Form[T any] struct {
	Data *T
}

func isGenericForm(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return false
	}
	if typ.PkgPath() == reflect.TypeOf(Form[any]{}).PkgPath() && strings.HasPrefix(strings.Split(typ.Name(), "[")[0]+"[", "Form[") {
		return true
	}
	return false
}
func init() {
	handlers.IsGenericForm = isGenericForm
}
