package vapi

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
)

func (web *webHandlerRunnerType) ResolveReceiverValue(handler webHandler) (reflect.Value, error) {
	key := handler.apiInfo.ReceiverType.String() + "/webHandlerRunnerType/ResolveReceiverValue"
	ret, err := OnceCall(key, func() (*reflect.Value, error) {
		result := handler.initFunc.Call([]reflect.Value{})
		if result[1].IsValid() && !result[1].IsNil() {
			return nil, result[1].Interface().(error)
		}

		return &result[0], nil
	})
	return *ret, err
}

func (web *webHandlerRunnerType) ExecFormPost(handler webHandler, w http.ResponseWriter, r *http.Request) error {

	ReceiverValue, err := web.ResolveReceiverValue(handler)
	if err != nil {
		return err
	}
	var bodyData reflect.Value

	// Duyệt tất cả key/value trong form
	if handler.apiInfo.IndexOfRequestBody > -1 {
		bodyData = reflect.New(handler.apiInfo.TypeOfRequestBodyElem)

		for key, values := range r.Form {
			field, ok := handler.apiInfo.TypeOfRequestBodyElem.FieldByNameFunc(func(s string) bool {
				return strings.EqualFold(s, key)
			})

			if !ok {
				continue
			}
			fieldType := field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
			if fieldType.Kind() == reflect.Slice {
				bodyData.Elem().FieldByIndex(field.Index).Set(reflect.ValueOf(values))
			} else if fieldType.Kind() == reflect.Struct {
				if len(values) == 0 {
					continue
				}

				fieldValue := reflect.New(fieldType)

				err := json.Unmarshal([]byte(values[0]), fieldValue.Elem().Addr().Interface())
				if err != nil {
					return err
				}
				bodyData.Elem().FieldByIndex(field.Index).Set(fieldValue.Elem())

			} else {
				bodyData.Elem().FieldByIndex(field.Index).Set(reflect.ValueOf(values[0]))
			}
		}
		if len(handler.apiInfo.FormUploadFile) > 0 {
			for _, index := range handler.apiInfo.FormUploadFile {
				field := handler.apiInfo.TypeOfRequestBodyElem.Field(index)
				fieldType := field.Type
				if fieldType.Kind() == reflect.Ptr {
					fieldType = fieldType.Elem()
				}

				var fileValues []*multipart.FileHeader
				for fileName, fhs := range r.MultipartForm.File {
					if strings.EqualFold(fileName, field.Name) {
						fileValues = fhs
						break

					}

				}
				if len(fileValues) == 0 {
					continue
				}

				if fieldType.Kind() == reflect.Slice {

					if fieldType.Elem() == reflect.TypeOf(&multipart.FileHeader{}).Elem() {
						dataValues := make([]multipart.FileHeader, len(fileValues))
						for i, fh := range fileValues {
							dataValues[i] = *fh
						}
						bodyData.Elem().FieldByIndex(field.Index).Set(reflect.ValueOf(dataValues))
					} else if fieldType.Elem() == reflect.TypeOf(&multipart.FileHeader{}) {

						bodyData.Elem().FieldByIndex(field.Index).Set(reflect.ValueOf(fileValues))

					} else if fieldType.Elem() == reflect.TypeOf((*multipart.File)(nil)).Elem() {
						dataValues := make([]multipart.File, len(fileValues))
						for i, fh := range fileValues {
							if fh == nil {
								continue
							}

							f, errOpen := fh.Open()
							if errOpen != nil {
								return errOpen
							}
							if f == nil {
								continue
							}
							dataValues[i] = f
						}
						bodyData.Elem().FieldByIndex(field.Index).Set(reflect.ValueOf(dataValues))
					}

				} else if fieldType == reflect.TypeOf(multipart.FileHeader{}) {
					bodyData.Field(index).Set(reflect.ValueOf(*fileValues[0]))
				} else if fieldType == reflect.TypeOf(&multipart.FileHeader{}).Elem() {
					bodyData.Field(index).Set(reflect.ValueOf(fileValues[0]))
				} else if fieldType == reflect.TypeOf((*multipart.File)(nil)).Elem() {
					file, err := fileValues[0].Open()
					if err != nil {
						return err
					}
					bodyData.Elem().FieldByIndex(field.Index).Set(reflect.ValueOf(file))

				}

			}
		}

	}
	args := make([]reflect.Value, handler.apiInfo.Method.Type.NumIn())
	args[0] = ReceiverValue
	if handler.apiInfo.IndexOfRequestBody > -1 {
		args[handler.apiInfo.IndexOfRequestBody] = bodyData

	}
	if handler.apiInfo.ReceiverTypeElem.Kind() == reflect.Ptr {
		handler.apiInfo.ReceiverTypeElem = handler.apiInfo.ReceiverTypeElem.Elem()
	}

	context := reflect.New(handler.apiInfo.TypeOfArgsElem)

	context.Elem().FieldByName("Req").Set(reflect.ValueOf(r))
	context.Elem().FieldByName("Res").Set(reflect.ValueOf(w))

	args[handler.apiInfo.IndexOfArg] = context
	for i := 1; i < handler.apiInfo.Method.Type.NumIn(); i++ {
		if handler.apiInfo.Method.Type.In(i).Kind() != reflect.Ptr {
			if args[i].Kind() == reflect.Ptr {
				args[i] = args[i].Elem()

			}
		}

	}
	retArgs := handler.apiInfo.Method.Func.Call(args)
	if len(retArgs) > 0 {
		if err, ok := retArgs[len(retArgs)-1].Interface().(error); ok {
			return err
		}
		if len(retArgs) > 2 {
			retIntefaces := []interface{}{}
			for i := 0; i < len(retArgs)-1; i++ {
				retIntefaces = append(retIntefaces, retArgs[i].Interface())
			}

			retArgs = retArgs[0 : len(retArgs)-2]
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(retIntefaces)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(retArgs[0].Elem().Interface())
		}
		// Ví dụ: trả về dạng JSON

	}

	return nil
}
