package dynacall

import (
	"encoding/json"
	"reflect"
)

type Invoker struct {
	// Data interface{}
	// fields     map[string]reflect.StructField
	Args       interface{}
	method     reflect.Method
	callerPath string
}

func (r *Invoker) New(callerPath string) error {
	inputType, method, err := GetInputTypeOfCallerPath(callerPath)
	if err != nil {
		return err
	}
	val := reflect.New(inputType[0]).Elem()
	// val := fields[i]
	ptr := reflect.New(inputType[0]) // Tạo một con trỏ mới đến kiểu của val (ví dụ: *main.MyData)
	ptr.Elem().Set(val)

	r.method = *method

	r.callerPath = callerPath
	r.Args = ptr.Interface()
	return nil
}
func NewInvoker(callerPath string) (*Invoker, error) {
	invoker := &Invoker{}
	err := invoker.New(callerPath)
	if err != nil {
		return nil, err
	}
	return invoker, nil

}

func (r *Invoker) Injector(injector interface{}) func() (interface{}, error) {
	return func() (interface{}, error) {
		return Call(r.callerPath, r.Args, injector)
	}
}
func (r *Invoker) LoadJSON(jsonStr string) error {
	err := json.Unmarshal([]byte(jsonStr), r.Args)
	if err != nil {
		return err
	}
	return nil
}
