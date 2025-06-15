package dynacall

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

/*
for args, the first argument is the owner ex: func(u*user) methodName ... owner is u

and the rest of the arguments are the parameters of the method

for example: func(u*user) update(name string, age int) error
args = []interface{}{u, "john", 25}
*/
func invoke(method reflect.Method, arg interface{}, injector interface{}) (interface{}, error) {
	funcType := method.Type

	ownerType := funcType.In(0) //first argument is the owner ex: func(u*user) methodName ... owner is u
	ownerInstance := createReceiverInstance(ownerType, injector)
	args := []interface{}{}
	if arg != nil {
		val := reflect.ValueOf(arg)
		if val.Type().Kind() == reflect.Slice {
			for i := 0; i < val.Len(); i++ {
				args = append(args, val.Index(i).Interface())
			}
			// args = arg.([]interface{})
		} else {
			args = append(args, arg)
		}
	}

	invokeArgs := make([]reflect.Value, len(args)+1)
	invokeArgs[0] = ownerInstance
	// argTypes := []reflect.Type{}
	for i := 1; i < funcType.NumIn(); i++ {
		// // paramVal := args[i-1]
		argType := funcType.In(i)
		valueType := reflect.TypeOf(args[i-1])
		if valueType.Kind() == reflect.Ptr {
			valueType = valueType.Elem()
		}
		if argType == valueType {
			invokeArgs[i] = reflect.ValueOf(args[i-1])
			continue
		}
		if argType.Kind() == reflect.Ptr {
			if argType.Elem() == valueType {
				invokeArgs[i] = reflect.ValueOf(args[i-1])
				continue
			}
		}
		kn := valueType.Kind().String()
		fmt.Println("kn: ", kn)
		if argType.Kind() == reflect.Struct && valueType.Kind() == reflect.Struct {
			invokeArgs[i] = reflect.ValueOf(args[i-1])
			continue
		}
		test := fmt.Sprintf("argType: %v, valueType: %v\n", argType, valueType)
		fmt.Println(test)

		if valueType.ConvertibleTo(argType) {

			val := reflect.ValueOf(args[i-1])
			if invokeArgs[i].Kind() == reflect.Ptr {
				invokeArgs[i] = val.Addr()
			} else {
				invokeArgs[i] = val
			}

		} else {
			if argType.Kind() == reflect.Ptr {
				argType = argType.Elem()
			}
			if argType == reflect.TypeOf(time.Now()) {
				if strTime, ok := args[i-1].(string); ok {
					timeVale, err := time.Parse(time.RFC3339, strTime)
					if err != nil {
						return nil, CallError{
							Code: CallErrorCodeInvalidArgs,
							Err:  fmt.Errorf("invalid time format, expected RFC3339, got %s", strTime),
						}
					}
					invokeArgs[i] = reflect.ValueOf(timeVale)
					continue

				}
			}
			if argType == reflect.TypeOf(uuid.UUID{}) {
				if strUUID, ok := args[i-1].(string); ok {
					uuidVale, err := uuid.Parse(strUUID)
					if err != nil {
						return nil, CallError{
							Code: CallErrorCodeInvalidArgs,
							Err:  fmt.Errorf("invalid uuid format, expected uuid, got %s", strUUID),
						}
					}
					invokeArgs[i] = reflect.ValueOf(uuidVale)
					continue
				}

			}

			return nil, CallError{Err: CallError{
				Code: CallErrorCodeInvalidArgs,
				Err:  fmt.Errorf("invalid args, expected %s, got %s", argType, args[i-1]),
			}}
		}
		// argTypes = append(argTypes, argType)

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
