package services

import (
	"reflect"
)

type serviceUtilsType struct {
	pkgPath string
	// checkSingletonTypeName string
	// checkScopeTypeName string
}

func (svc *serviceUtilsType) checkSingletonTypeName() string {
	return GetCheckSingletonTypeName()
}
func (svc *serviceUtilsType) checkScopeTypeName() string {
	return GetCheckScopeTypeName()
}

var ServiceUtils *serviceUtilsType
var GetCheckSingletonTypeName func() string
var GetCheckScopeTypeName func() string

func init() {
	ServiceUtils = &serviceUtilsType{
		pkgPath: reflect.TypeOf(serviceUtilsType{}).PkgPath(),
		// checkSingletonTypeName: strings.Split(reflect.TypeOf(Singleton[any]{}).String(), "[")[0] + "[",
		// checkScopeTypeName:     strings.Split(reflect.TypeOf(Scoped[any]{}).String(), "[")[0] + "[",
	}
	// ServiceUtils.checkScopeTypeName = GetCheckScopeTypeName()
	// ServiceUtils.checkSingletonTypeName = GetCheckSingletonTypeName()

}
