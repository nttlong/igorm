package handlers

import (
	"fmt"
	"reflect"
)

func (h *helperType) FindNewMehodOfHandler(info *HandlerInfo) error {
	// handlerType := info.Method.Type.In(info.IndexOfArg)
	typ := info.TypeOfArgs
	if typ.Kind() == reflect.Struct {
		prtType := reflect.PointerTo(typ)

		return fmt.Errorf("argument %d in method  %s.%s is not a pointer to struct; expected %s but got %s ", info.IndexOfArg, info.ReceiverType.String(), info.Method.Name, prtType.String(), typ.String())
	}
	for j := 0; j < typ.NumMethod(); j++ {
		if typ.Method(j).Name == "New" {
			mt := typ.Method(j)
			if mt.Type.NumOut() != 1 {
				return fmt.Errorf("%s.New must return one value (error)", typ.String())
			}
			errorType := reflect.TypeOf((*error)(nil)).Elem()

			if mt.Type.Out(0) != errorType {
				return fmt.Errorf("%s.New must return one value (error)", typ.String())
			}
			info.NewMethodOfHandler = &mt
			return nil

		}
	}
	handlerType := info.TypeOfArgsElem
	// if handlerType.Kind() == reflect.Ptr {
	// 	handlerType = handlerType.Elem()
	// }
	fieldIndex := []int{}
	for _, x := range info.FieldIndex {
		fieldIndex := append(fieldIndex, x)
		typ := handlerType.FieldByIndex(fieldIndex).Type
		if typ.Kind() == reflect.Struct {
			typ = reflect.PointerTo(typ)
		}
		for j := 0; j < typ.NumMethod(); j++ {
			if typ.Method(j).Name == "New" {
				mt := typ.Method(j)
				if mt.Type.NumOut() != 1 {
					return fmt.Errorf("%s.New must return one value (error)", typ.String())
				}
				errorType := reflect.TypeOf((*error)(nil)).Elem()

				if mt.Type.Out(0) != errorType {
					return fmt.Errorf("%s.New must return one value (error)", typ.String())
				}
				info.NewMethodOfHandler = &mt
				return nil

			}
		}
	}
	return nil
}
