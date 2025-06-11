package dynacall

import "reflect"

func createReceiverInstance(receiverType reflect.Type, injector interface{}) reflect.Value {
	receiverStructType := receiverType.Elem()
	receiverInstance := reflect.New(receiverStructType)
	if injector != nil {
		applyInjector(receiverInstance, injector)
	}
	return receiverInstance
}
func applyInjector(receiverInstance reflect.Value, injector interface{}) {
	injectorType := reflect.TypeOf(injector)
	injectorValue := reflect.ValueOf(injector)
	for i := 0; i < injectorType.NumField(); i++ {
		field := injectorType.Field(i)
		fieldName := field.Name

		fieldValue := injectorValue.FieldByName(fieldName)
		if fieldValue.IsValid() {
			receiverField := receiverInstance.Elem().FieldByName(fieldName)
			if receiverField.IsValid() {
				receiverField.Set(fieldValue)
			}
		}
	}
}
