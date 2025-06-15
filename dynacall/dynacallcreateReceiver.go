package dynacall

import (
	"reflect"
)

func createReceiverInstance(receiverType reflect.Type, injector interface{}) reflect.Value {
	receiverStructType := receiverType.Elem()
	receiverInstance := reflect.New(receiverStructType)
	if injector != nil {
		applyInjector(receiverType, receiverInstance, injector)
	}
	return receiverInstance
}
func applyInjector(receiverType reflect.Type, receiverInstance reflect.Value, injector interface{}) {
	injectorType := reflect.TypeOf(injector)
	injectorValue := reflect.ValueOf(injector)
	if injectorType.Kind() == reflect.Ptr {
		injectorType = injectorType.Elem()
		injectorValue = injectorValue.Elem()
	}
	for i := 0; i < injectorType.NumField(); i++ {
		field := injectorType.Field(i)
		fieldName := field.Name

		fieldValue := injectorValue.FieldByName(fieldName)

		if fieldValue.IsValid() {

			receiverInstanceEle := receiverInstance
			if receiverInstance.Kind() == reflect.Ptr {
				receiverInstanceEle = receiverInstance.Elem()
			}

			receiverField := receiverInstanceEle.FieldByName(fieldName)
			if receiverField.IsValid() {
				receiverField.Set(fieldValue)
			}
		}
	}
	if receiverType.Kind() == reflect.Ptr {
		receiverType = receiverType.Elem()
		receiverInstance = receiverInstance.Elem()
	}
	for i := 0; i < receiverType.NumField(); i++ {
		field := receiverType.Field(i)

		if field.Anonymous && field.Type != reflect.TypeOf(Caller{}) {
			if field.Type.Kind() == reflect.Ptr {
				field.Type = field.Type.Elem()
			}
			if field.Type.Kind() == reflect.Ptr {
				field.Type = field.Type.Elem()
			}
			fieldValue := receiverInstance.Field(i)
			if fieldValue.IsValid() {

				applyInjector(field.Type, fieldValue, injector)
			}
		}
	}

}
