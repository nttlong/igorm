package unvsdi

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type unvsdiUtils struct {
	currentPackage    string
	once              sync.Once
	injectorTypeNames []string
}

func (u *unvsdiUtils) getCurrentPackage() string {
	u.once.Do(func() {
		u.currentPackage = reflect.TypeOf(unvsdiUtils{}).PkgPath()
	})
	return u.currentPackage
}
func (u *unvsdiUtils) isInjector(field reflect.StructField) bool {
	fmt.Println(field.Type.PkgPath())
	if field.Type.PkgPath() != u.getCurrentPackage() {
		return false
	}
	fieldTypeName := field.Type.String()
	isInjector := false
	for _, injectorTypeName := range u.injectorTypeNames {
		if strings.HasPrefix(fieldTypeName, injectorTypeName+"[") {
			isInjector = true
			break
		}
	}

	return isInjector
}

var utils = &unvsdiUtils{
	injectorTypeNames: []string{ // list of injector type names
		strings.Split(reflect.TypeOf(Singleton[any, any]{}).String(), "[")[0],
		strings.Split(reflect.TypeOf(Scoped[any, any]{}).String(), "[")[0],
		strings.Split(reflect.TypeOf(Transient[any, any]{}).String(), "[")[0],
	},
}
