package vapi

// type injectorUtilType struct {
// }

// func (injector *injectorUtilType) LoadByType(typ reflect.Type) (*reflect.Value, error) {

// 	var newMethod reflect.Method
// 	foungNewMethod := false
// 	for i := 0; i < typ.NumMethod(); i++ {
// 		if typ.Method(i).Name == "New" {
// 			newMethod = typ.Method(i)
// 			foungNewMethod = true

// 			break
// 		}
// 	}
// 	if !foungNewMethod {
// 		errMsg := fmt.Sprintf("New function was not found in %s. injector need New function", typ.String())
// 		return nil, errors.New(errMsg)

// 	}
// 	if newMethod.Type.NumOut() != 1 {
// 		errMsg := fmt.Sprintf("New function of %s must return 1 value (error or nil)", typ.String())
// 		return nil, errors.New(errMsg)
// 	}
// 	instanceType := typ
// 	if instanceType.Kind() == reflect.Ptr {
// 		instanceType = instanceType.Elem()
// 	}

// 	ret := reflect.New(instanceType)
// 	retVal := newMethod.Func.Call([]reflect.Value{ret})
// 	if retVal[0].Interface() != nil {
// 		return nil, retVal[0].Interface().(error)
// 	}
// 	return &ret, nil
// }

// var InjectorUtil = &injectorUtilType{}

// func LoadInject[T any]() (*T, error) {
// 	ret, err := InjectorUtil.LoadByType(reflect.TypeFor[*T]())
// 	if err != nil {
// 		return nil, err
// 	}
// 	retT := ret.Interface()
// 	fx := retT.(T)
// 	return &fx, nil

// }
