package wx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type newMthodInfoOfService struct {
	IndexOfHttpContext     int
	IndexOfInjectords      []int
	IndexInjectHttpService []int // index of inject http service, if not found, it will be empty
}
type serviceUtilType struct {
	prefixHttpServiceCheck      string
	packagePathHttpServiceCheck string
}

func (s *serviceUtilType) isInjectHttpService(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return false
	}

	if typ.PkgPath() == s.packagePathHttpServiceCheck && strings.HasPrefix(typ.Name(), s.prefixHttpServiceCheck) {
		return true
	}
	return false
}
func (s *serviceUtilType) isHttpServiceinternal(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return false
	}

	if typ == reflect.TypeOf(HttpContext{}) || typ == reflect.TypeOf(&HttpContext{}) {
		return true
	}
	return false
}

type initIsHttpService struct {
	val  bool
	once sync.Once
}

var cachIsHttpService sync.Map

func (s *serviceUtilType) isHttpService(typ reflect.Type) bool {
	actual, _ := cachIsHttpService.LoadOrStore(typ, &initIsHttpService{})
	isHttpServiceItem := actual.(*initIsHttpService)
	isHttpServiceItem.once.Do(func() {
		isHttpServiceItem.val = s.isHttpServiceinternal(typ)
	})
	return isHttpServiceItem.val
}

type initExtractInfo struct {
	val  *newMthodInfoOfService
	once sync.Once
	err  error
}

var cacheExtractInfo sync.Map

func (s *serviceUtilType) ExtractInfo(method reflect.Method) (*newMthodInfoOfService, error) {
	actual, _ := cacheExtractInfo.LoadOrStore(method, &initExtractInfo{})
	initExtractInfoItem := actual.(*initExtractInfo)
	initExtractInfoItem.once.Do(func() {
		initExtractInfoItem.val, initExtractInfoItem.err = s.extractInfo(method)
	})
	return initExtractInfoItem.val, initExtractInfoItem.err
}
func (s *serviceUtilType) extractInfo(method reflect.Method) (*newMthodInfoOfService, error) {
	ret := &newMthodInfoOfService{
		IndexOfHttpContext: -1,
		IndexOfInjectords:  []int{},
	}
	unknownArgs := []int{}
	unknownArgsType := []reflect.Type{}
	for i := 1; i < method.Type.NumIn(); i++ {
		typ := method.Type.In(i)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct {
			continue
		}

		if s.isHttpService(typ) {
			if ret.IndexOfHttpContext == -1 {
				ret.IndexOfHttpContext = i
			} else {
				return nil, fmt.Errorf("%s.%s has more than one HttpService argument", method.Type.In(0).String(), method.Name)
			}
			continue
		}
		if Helper.IsGenericDepen(typ) {
			ret.IndexOfInjectords = append(ret.IndexOfInjectords, i)
			continue
		}
		if s.isInjectHttpService(typ) {
			ret.IndexInjectHttpService = append(ret.IndexInjectHttpService, i)
			continue
		}
		unknownArgs = append(unknownArgs, i)
		unknownArgsType = append(unknownArgsType, typ)

	}
	if ret.IndexOfHttpContext == -1 {
		return nil, fmt.Errorf("%s.%s has no HttpService argument", method.Type.In(0).String(), method.Name)
	}
	if len(unknownArgs) > 0 {
		reciverType := method.Type.In(0)
		if reciverType.Kind() == reflect.Ptr {
			reciverType = reciverType.Elem()
		}
		uneceptabeTypes := []string{}
		for _, arg := range unknownArgs {
			uneceptabeTypes = append(uneceptabeTypes, method.Type.In(arg).String()+fmt.Sprintf(" at args %d", arg))
		}

		uneceptabeType := strings.Join(uneceptabeTypes, ", ")
		msg := fmt.Sprintf("%s.%s has uneceptabe types: %v", reciverType.Name(), method.Name, uneceptabeType)
		return nil, errors.New(msg)
	}
	return ret, nil
}

var ServiceUtil = &serviceUtilType{
	prefixHttpServiceCheck:      strings.Split(reflect.TypeOf(HttpService[any]{}).Name(), "[")[0] + "[",
	packagePathHttpServiceCheck: reflect.TypeOf(HttpService[any]{}).PkgPath(),
}
