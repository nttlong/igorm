package wx

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type depenResolvers struct {
}

func (de *depenResolvers) ResolveDependsType(typ reflect.Type) (*reflect.Value, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	instanceVal := reflect.New(typ)
	de.resoleDepenceFiels(instanceVal, map[reflect.Type]bool{})

	return &instanceVal, nil
}
func (de *depenResolvers) ResolveType(typ reflect.Type) (*reflect.Value, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	//ret := reflect.New(typ)
	if isStructDepenType(typ) {
		ret, err := de.ResolveDependsType(typ)
		return ret, err
	}

	nm, err := de.FindNewMethod(typ)
	if err != nil {

		return nil, err
	}
	if nm == nil {
		return nil, nil
	}

	retVale, err := de.RunNewMethod(*nm)
	if err != nil {
		return nil, err
	}

	return retVale, nil

}

type intResolveType struct {
	ret  *reflect.Value
	err  error
	once sync.Once
}

var cacheResolveType sync.Map

func (de *depenResolvers) ResolveTypeOnce(typ reflect.Type) (*reflect.Value, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	key := typ.String()
	actual, _ := cacheResolveType.LoadOrStore(key, &intResolveType{})
	item := actual.(*intResolveType)
	item.once.Do(func() {
		item.ret, item.err = de.ResolveType(typ)

	})
	return item.ret, item.err
}

type initFindNewMethdodOnly struct {
	ret  *reflect.Method
	err  error
	once sync.Once
}

var cacheFindNewMethdodOnly sync.Map

func (de *depenResolvers) findNewMethdodOnly(typ reflect.Type) *reflect.Method {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	key := typ.String()
	actual, _ := cacheFindNewMethdodOnly.LoadOrStore(key, &initFindNewMethdodOnly{})
	item := actual.(*initFindNewMethdodOnly)
	item.once.Do(func() {
		item.ret = de.findNewMethdodOnlyNoCache(typ)
	})
	return item.ret
}
func (de *depenResolvers) findNewMethdodOnlyNoCache(typ reflect.Type) *reflect.Method {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	typPtr := reflect.PointerTo(typ)
	for i := 0; i < typPtr.NumMethod(); i++ {
		m := typPtr.Method(i)
		if m.Name == "New" {
			return &m
		}
	}
	return nil

}
func (de *depenResolvers) findNewMethodWithVisited(typ reflect.Type, visited map[reflect.Type]bool) (*reflect.Method, error) {

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if _, ok := visited[typ]; ok {
		msg := "recursive reference:/n"
		for k := range visited {
			msg += k.String() + "-->"
		}
		msg += typ.String() + ",in new method of " + typ.String() + ".New"
		return nil, errors.New(msg)
	}
	visited[typ] = true
	// typPtr := reflect.PointerTo(typ)
	m := de.findNewMethdodOnly(typ)
	if m == nil {
		if isTypeDepen(typ, map[reflect.Type]bool{}) {
			return nil, nil
		}
		return nil, errors.New("" + typ.String() + " does not have New method")
	}
	for j := 1; j < m.Type.NumIn(); j++ {
		typIn := m.Type.In(j)
		if typIn.Kind() == reflect.Ptr {
			typIn = typIn.Elem()
		}
		if typIn.Kind() != reflect.Struct {
			return nil, fmt.Errorf("args %d of %s.%s is not struct", j, typ.String(), m.Name)
		}
		if isTypeDepen(typIn, visited) {
			continue
		}
		nextNewhMethod, err := de.findNewMethodWithVisited(typIn, visited)
		if err != nil {
			return nil, err
		}
		if nextNewhMethod == nil {
			return nil, fmt.Errorf("args %d of %s.%s is does not have New method, plesae review %s", j, typ.String(), m.Name, typIn.String())
		}
		if !de.hasErrorReturn(*nextNewhMethod) {
			msgEror := fmt.Sprintf("%s.%s must return error", nextNewhMethod.Type.In(0).String(), m.Name)
			return nil, errors.New(msgEror)
		}

	}
	return m, nil

}
func (de *depenResolvers) hasErrorReturn(m reflect.Method) bool {
	// Lấy kiểu của phương thức
	methodType := m.Type

	// Lấy số lượng giá trị trả về
	numOut := methodType.NumOut()

	// Duyệt qua tất cả các giá trị trả về
	for i := 0; i < numOut; i++ {
		// Lấy kiểu của giá trị trả về thứ i
		returnType := methodType.Out(i)

		// Kiểm tra xem kiểu trả về đó có thỏa mãn interface 'error' không
		// TypeOf(new(error)).Elem() sẽ trả về reflect.Type của interface 'error'
		if returnType.Implements(reflect.TypeOf(new(error)).Elem()) {
			return true
		}
	}
	return false
}

func (de *depenResolvers) FindNewMethod(typ reflect.Type) (*reflect.Method, error) {
	return de.findNewMethodWithVisited(typ, map[reflect.Type]bool{})
}
func (de *depenResolvers) resoleDepenceFiels(ownerVal reflect.Value, visited map[reflect.Type]bool) (*reflect.Value, error) {
	typ := ownerVal.Type()

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, nil
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !isTypeDepen(field.Type, visited) {
			continue
		}
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		instanceVal := reflect.New(fieldType)
		if appField, ok := fieldType.FieldByName("Owner"); ok {
			fieldSet := instanceVal.Elem().FieldByName(appField.Name)
			fieldSet.Set(reflect.Zero(fieldSet.Type()))
			if fieldSet.IsValid() {
				fieldSet.Set(ownerVal)
			} else {
				fieldSet.Elem().Set(ownerVal)

			}

		}

		fieldSet := ownerVal.Elem().FieldByName(field.Name)

		if fieldSet.IsValid() {
			if fieldType.Kind() == reflect.Struct {
				fieldSet.Set(instanceVal.Elem())
				continue
			} else {
				fieldSet.Set(instanceVal)
			}

		} else {
			fieldSet.Set(reflect.ValueOf((instanceVal).Interface()))
		}
	}

	return &ownerVal, nil

}
func (de *depenResolvers) applyOwnerValue(ownerVal reflect.Value, typ reflect.Type, value reflect.Value) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if appField, ok := typ.FieldByName("Owner"); ok {

		// if value.Kind() == reflect.Ptr {
		// 	value = value.Elem()
		// }
		fieldSet := value.Elem().FieldByName(appField.Name)
		// if fieldSet.Kind() == reflect.Ptr {
		// 	fieldSet = fieldSet.Elem()
		// }
		// fieldSet.Set(reflect.Zero(fieldSet.Type()))
		if fieldSet.IsValid() {
			fieldSet.Set(ownerVal)
		}

	}

}

func (de *depenResolvers) RunNewMethodWithReceiver(retVale reflect.Value, nm reflect.Method) (*reflect.Value, error) {
	// if retVale.Kind() == reflect.Ptr {
	// 	retVale = retVale.Elem()
	// }
	de.resoleDepenceFiels(retVale, map[reflect.Type]bool{})
	args := make([]reflect.Value, nm.Type.NumIn())
	args[0] = retVale

	for i := 1; i < nm.Type.NumIn(); i++ {
		typ := nm.Type.In(i)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct {

			continue
		}
		if isStructDepenType(typ) {
			ret, err := de.ResolveType(typ)
			if err != nil {
				return nil, err
			}
			de.applyOwnerValue(retVale, typ, *ret)
			if nm.Type.In(i).Kind() == reflect.Ptr {
				args[i] = *ret
				continue
			} else {
				args[i] = (*ret).Elem()
			}

			continue
		}
		if isTypeDepen(typ, map[reflect.Type]bool{}) {
			ret, err := de.ResolveType(typ)
			if err != nil {
				return nil, err
			}
			if nm.Type.In(i).Kind() == reflect.Ptr {
				args[i] = *ret
				continue
			} else {
				args[i] = (*ret).Elem()
			}

			continue
		}
		nmOfArgs, err := de.FindNewMethod(typ)
		if err != nil {
			return nil, err
		}
		if nmOfArgs == nil {
			return nil, fmt.Errorf("method new was not found in %s", typ.String())
		}
		argval, err := de.RunNewMethod(*nmOfArgs)
		if err != nil {
			return nil, err
		}
		args[i] = *argval

	}

	ret := nm.Func.Call(args)
	if !ret[len(ret)-1].IsNil() {
		return nil, ret[len(ret)-1].Interface().(error)
	}
	return &retVale, nil

}
func (de *depenResolvers) RunNewMethod(nm reflect.Method) (*reflect.Value, error) {
	retVale := reflect.New(nm.Type.In(0).Elem())
	de.resoleDepenceFiels(retVale, map[reflect.Type]bool{})
	args := make([]reflect.Value, nm.Type.NumIn())
	args[0] = retVale

	for i := 1; i < nm.Type.NumIn(); i++ {
		typ := nm.Type.In(i)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct {

			continue
		}
		if isStructDepenType(typ) {
			ret, err := de.ResolveType(typ)
			if err != nil {
				return nil, err
			}
			if nm.Type.In(i).Kind() == reflect.Ptr {
				args[i] = *ret
				continue
			} else {
				args[i] = (*ret).Elem()
			}

			continue
		}
		if isTypeDepen(typ, map[reflect.Type]bool{}) {
			ret, err := de.ResolveType(typ)
			if err != nil {
				return nil, err
			}
			if nm.Type.In(i).Kind() == reflect.Ptr {
				args[i] = *ret
				continue
			} else {
				args[i] = (*ret).Elem()
			}

			continue
		}
		nmOfArgs, err := de.FindNewMethod(typ)
		if err != nil {
			return nil, err
		}
		if nmOfArgs == nil {
			return nil, fmt.Errorf("method new was not found in %s", typ.String())
		}
		argval, err := de.RunNewMethod(*nmOfArgs)
		if err != nil {
			return nil, err
		}
		args[i] = *argval

	}

	ret := nm.Func.Call(args)
	if !ret[len(ret)-1].IsNil() {
		return nil, ret[len(ret)-1].Interface().(error)
	}
	return &retVale, nil

}

var DepenResolvers = &depenResolvers{}

func New[T any]() (*T, error) {
	ret, err := DepenResolvers.ResolveType(reflect.TypeFor[T]())
	if err != nil {
		return nil, err
	}
	if ret == nil {
		return nil, nil
	}
	return ret.Interface().(*T), nil

}
