/*
this code is used to call a method of a struct dynamically use as private method
*/
package dynacall

import (
	"encoding/json"
	"reflect"
)

/*
for args, the first argument is the owner ex: func(u*user) methodName ... owner is u

and the rest of the arguments are the parameters of the method

for example: func(u*user) update(name string, age int) error
args = []interface{}{u, "john", 25}
*/
func invoke(method reflect.Method, args []interface{}, injector interface{}) (interface{}, error) {
	funcType := method.Type

	ownerType := funcType.In(0) //first argument is the owner ex: func(u*user) methodName ... owner is u
	ownerInstance := createReceiverInstance(ownerType, injector)
	invokeArgs := make([]reflect.Value, len(args)+1)
	invokeArgs[0] = ownerInstance
	argTypes := []reflect.Type{}
	for i := 1; i < funcType.NumIn(); i++ {
		// paramVal := args[i-1]
		argType := funcType.In(i)
		argTypes = append(argTypes, argType)
	}
	_, structInstance, f, err := CreateDynamicStruct(argTypes)
	if err != nil {
		return nil, err
	}

	wrapperArgs := map[string]interface{}{}
	for i := 0; i < len(args); i++ {
		wrapperArgs[f[i].Name] = args[i]
	}
	jsonBytes, err := json.Marshal(wrapperArgs)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, structInstance.Addr().Interface())
	if err != nil {
		return nil, err
	}
	for i := 0; i < structInstance.NumField(); i++ {
		invokeArgs[i+1] = structInstance.Field(i)

	}

	ret := method.Func.Call(invokeArgs)
	return ret, nil
}
