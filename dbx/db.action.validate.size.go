package dbx

import (
	"reflect"
	"strconv"
)

func validateSize(entity interface{}) *DBXError {
	isMultiple := false

	if entity == nil {
		return nil
	}
	data := reflect.ValueOf(entity)
	typ := reflect.TypeOf(entity)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		data = data.Elem()

	}
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
		isMultiple = true
	}
	if typ.Kind() == reflect.Array {
		typ = typ.Elem()
		isMultiple = true
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()

	}
	et, err := CreateEntityType(typ)
	if err != nil {
		return nil
	}
	if isMultiple {
		return validateSizeMultiple(et, entity)
	} else {
		return validateSizeSingle(et, data)
	}

}
func validateSizeSingle(et *EntityType, data reflect.Value) *DBXError {
	for _, f := range et.EntityFields {
		if f.MaxLen > 0 && f.NonPtrFieldType.Kind() == reflect.String {
			v := data.FieldByName(f.Name)
			if v.Kind() == reflect.Ptr {
				if v.Elem() == reflect.Zero(v.Type()).Elem() {

					continue
				}
				v = v.Elem()

			}
			valData := v.Interface().(string)

			if len(valData) > f.MaxLen {
				return &DBXError{
					Code:           DBXErrorCodeInvalidSize,
					Message:        "Field " + f.Name + " exceeds the maximum length of " + strconv.Itoa(f.MaxLen),
					TableName:      et.TableName,
					ConstraintName: "",
					Fields:         []string{f.Name},
					Values:         []string{valData},
					MaxSize:        f.MaxLen,
				}
			}
		}
	}
	return nil

}

func validateSizeMultiple(et *EntityType, entity interface{}) *DBXError {
	panic("not implemented")
}
