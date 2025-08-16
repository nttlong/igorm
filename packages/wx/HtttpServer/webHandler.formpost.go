package htttpserver

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	vapiErr "wx/errors"
)

func (web *webHandlerRunnerType) ExecFormPost(handler WebHandler, w http.ResponseWriter, r *http.Request) error {

	ReceiverValue, err := web.ResolveReceiverValue(handler, r)
	if err != nil {
		return err
	}
	var bodyData reflect.Value

	// Duyệt tất cả key/value trong form
	if handler.ApiInfo.IndexOfRequestBody > -1 {
		bodyData = reflect.New(handler.ApiInfo.TypeOfRequestBodyElem)
		if len(handler.ApiInfo.FormUploadFile) > 0 {
			for _, index := range handler.ApiInfo.FormUploadFile {
				field := handler.ApiInfo.TypeOfRequestBodyElem.Field(index)
				fieldType := field.Type
				if fieldType.Kind() == reflect.Ptr {
					fieldType = fieldType.Elem()
				}
				// r.FormFile("Files") //<-- r *http.Request
				fileValues, ok := r.MultipartForm.File[field.Name]
				if !ok {
					if field.Type == reflect.TypeOf(multipart.FileHeader{}) {
						msgError := fmt.Sprintf("%s was not found,%s is required", field.Name, field.Name)
						return vapiErr.NewParamMissMatchError(msgError)

					}

				}
				delete(r.MultipartForm.File, field.Name)

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
					field := handler.ApiInfo.TypeOfRequestBodyElem.Field(index)
					fieldSet := bodyData.Elem().FieldByIndex(field.Index)

					if len(fileValues) == 0 {
						if fieldSet.Kind() == reflect.Ptr {
							fieldSet.Set(reflect.Zero(fieldSet.Type()))
						} else {
							fieldSet.Set(reflect.Zero(fieldSet.Type()))
						}
						continue

					}
					valueSet := reflect.ValueOf(*fileValues[0])
					if fieldSet.IsValid() && fieldSet.CanConvert(valueSet.Type()) {
						fieldSet.Set(valueSet)
					}

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

		for key, values := range r.Form {
			fmt.Println(key)
			field, ok := handler.ApiInfo.TypeOfRequestBodyElem.FieldByNameFunc(func(s string) bool {
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
				fieldSetValue := bodyData.Elem().FieldByIndex(field.Index)
				valueSet := reflect.ValueOf(values)
				if fieldSetValue.CanConvert(valueSet.Type()) {
					fieldSetValue.Set(valueSet)
					// continue
				}

				// bodyData.Elem().FieldByIndex(field.Index).Set(reflect.ValueOf(values))
			} else if fieldType.Kind() == reflect.Struct {
				if len(values) == 0 {
					continue
				}

				fieldValue := reflect.New(fieldType)

				err := json.Unmarshal([]byte(values[0]), fieldValue.Elem().Addr().Interface())
				if err != nil {
					return err
				}
				fieldSetValue := bodyData.Elem().FieldByIndex(field.Index)
				valueSet := fieldValue.Elem()
				if fieldSetValue.CanConvert(valueSet.Type()) {
					fieldSetValue.Set(valueSet)
					// continue
				}

			} else {
				fieldSetValue := bodyData.Elem().FieldByIndex(field.Index)
				valueSet := reflect.ValueOf(values[0])

				if fieldSetValue.CanConvert(valueSet.Type()) {
					fieldSetValue.Set(valueSet)
					// continue
				}

				// bodyData.Elem().FieldByIndex(field.Index).Set(reflect.ValueOf(values[0]))
			}
		}

	}
	args := make([]reflect.Value, handler.ApiInfo.Method.Type.NumIn())
	args[0] = ReceiverValue
	if handler.ApiInfo.IndexOfRequestBody > -1 {
		args[handler.ApiInfo.IndexOfRequestBody] = bodyData

	}
	if handler.ApiInfo.ReceiverTypeElem.Kind() == reflect.Ptr {
		handler.ApiInfo.ReceiverTypeElem = handler.ApiInfo.ReceiverTypeElem.Elem()
	}

	context, err := web.CreateHttpContext(handler, w, r)
	if err != nil {
		return err
	}

	args[handler.ApiInfo.IndexOfArg] = context
	if handler.ApiInfo.IndexOfAuthClaimsArg > 0 {
		authType := handler.ApiInfo.Method.Type.In(handler.ApiInfo.IndexOfAuthClaimsArg)
		if authType.Kind() == reflect.Ptr {
			authType = authType.Elem()
		}
		authValue := reflect.New(authType)
		args[handler.ApiInfo.IndexOfAuthClaimsArg] = authValue

	}
	injectorArgs, err := web.LoadInjector(handler, r, w)
	if err != nil {
		return err
	}
	for i, injectIndex := range handler.ApiInfo.IndexOfInjectors {
		args[injectIndex] = injectorArgs[i]
	}

	retArgs, err := web.MethodCall(handler, args)
	if err != nil {
		return err
	}
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

			json.NewEncoder(w).Encode(retArgs[0].Interface())
		}
		// Ví dụ: trả về dạng JSON

	}

	return nil
}
