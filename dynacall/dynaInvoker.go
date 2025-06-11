/*
this code is used to call a method of a struct dynamically use as private method
*/
package dynacall

import (
	"encoding/json"
	"fmt"
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
	if len(f) != len(args) {
		return nil, CallError{
			Code: CallErrorCodeInvalidArgs,
			Err:  fmt.Errorf("invalid number of arguments, expected %d, got %d", len(f), len(args)),
		}
	}
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
	retData := make([]interface{}, method.Type.NumOut())
	for i := 0; i < method.Type.NumOut(); i++ {
		//get all the return values and fecth into retData
		retData[i] = ret[i].Interface()
	}
	//outPutType := method.funcType.Out(0)
	retErr := retData[len(retData)-1]
	if retErr != nil {
		return nil, retErr.(error)
	}
	if len(retData) == 2 {
		return retData[0], nil
	}
	return ret[0 : len(ret)-1], nil
}
