package handlers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	wxErrors "wx/errors"
)

// type formBodyItem struct {
// 	IndexFields [][]int
// 	Value       interface{}
// 	IsRequire   bool
// }

func (reqExec *RequestExecutor) GetFormValue(handlerInfo HandlerInfo, r *http.Request) (*reflect.Value, error) {
	bodyDataRet := reflect.New(handlerInfo.TypeOfRequestBodyElem)
	bodyData := bodyDataRet.Elem()
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return nil, wxErrors.NewFileParseError("error parsing multipart form", err)
	}

	//scan all post files
	if r.MultipartForm != nil && len(r.MultipartForm.File) > 0 {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			return nil, wxErrors.NewFileParseError("error parsing multipart form", err)
		}

		for key, values := range r.MultipartForm.File {
			field, ok := handlerInfo.TypeOfRequestBodyElem.FieldByNameFunc(func(s string) bool {
				return strings.EqualFold(s, key)
			})
			if !ok {
				continue
			}

			fileFieldSet := bodyData.FieldByIndex(field.Index)
			if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Slice {
				eleType := field.Type.Elem().Elem()
				if eleType == reflect.TypeOf(&multipart.FileHeader{}) { //<--*[]*multipart.FileHeader
					fileFieldSet.Set(reflect.ValueOf(values))
				} else if eleType == reflect.TypeOf(multipart.FileHeader{}) { //<--*[]multipart.File
					files := make([]multipart.FileHeader, len(values))
					for i, v := range values {
						files[i] = *v
					}
					fileFieldSet := reflect.New(fileFieldSet.Type().Elem())
					fileFieldSet.Elem().Set(reflect.ValueOf(files))

					//fileFieldSet.Set(vPtr)

				}
			}
			if field.Type.Kind() == reflect.Slice {
				eleType := field.Type.Elem()
				if eleType == reflect.TypeOf(&multipart.FileHeader{}) { //<--[]*multipart.FileHeader

					fileFieldSet.Set(reflect.ValueOf(values))
				} else if eleType == reflect.TypeOf(multipart.FileHeader{}) { //<--[]multipart.File
					files := make([]multipart.FileHeader, len(values))
					for i, v := range values {
						files[i] = *v
					}
					fileFieldSet.Set(reflect.ValueOf(files))
				}

			}
			if field.Type == reflect.TypeOf(&multipart.FileHeader{}) {
				fileFieldSet.Set(reflect.ValueOf(values[0]))
			}
			if field.Type == reflect.TypeOf(multipart.FileHeader{}) {
				fileFieldSet.Set(reflect.ValueOf(*values[0]))
			}

		}
	}
	for key, values := range r.PostForm {
		field, ok := handlerInfo.TypeOfRequestBodyElem.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, key)
		})
		if !ok {
			continue
		}
		fileFieldSet := bodyData.FieldByIndex(field.Index)
		if fileFieldSet.Kind() == reflect.Ptr {
			eleType := fileFieldSet.Type().Elem()
			if eleType.Kind() == reflect.Slice {
				fileFieldSet.Set(reflect.ValueOf(values))
			} else if eleType.Kind() == reflect.String {
				fileFieldSet.Set(reflect.ValueOf(values).Elem())
			} else if eleType.Kind() == reflect.Struct {
				value := reflect.New(eleType)
				data := value.Interface()
				err := json.Unmarshal([]byte(values[0]), data)
				if err != nil {
					return nil, err
				}
				fileFieldSet.Set(value)
			}

			continue
		}
		if fileFieldSet.Kind() == reflect.Slice {
			eleType := fileFieldSet.Type().Elem()
			if eleType.Kind() == reflect.Ptr {
				fileFieldSet.Set(reflect.ValueOf(values))
			} else {
				fileFieldSet.Set(reflect.ValueOf(values).Elem())
			}
			continue
		}
		if fileFieldSet.Kind() == reflect.String {
			fileFieldSet.Set(reflect.ValueOf(values[0]))
			continue
		}
		if fileFieldSet.Kind() == reflect.Struct {
			value := reflect.New(fileFieldSet.Type())
			data := value.Interface()
			err := json.Unmarshal([]byte(values[0]), data)
			if err != nil {
				return nil, err
			}
			fileFieldSet.Set(value.Elem())
			continue
		}
		//panic("not implete at file packages\\wx\\handlers\\helper.excutor.DoPostForm.go")
	}

	return &bodyDataRet, nil

}
func (reqExec *RequestExecutor) DoFormPost(handlerInfo HandlerInfo, r *http.Request, w http.ResponseWriter) (interface{}, error) {
	ctlValue, err := reqExec.CreateControllerValue(handlerInfo)
	if err != nil {
		return nil, wxErrors.NewServiceInitError(err.Error())
	}
	controllerValue := *ctlValue

	args := make([]reflect.Value, handlerInfo.Method.Func.Type().NumIn())
	args[0] = controllerValue
	args[handlerInfo.IndexOfArg] = reflect.New(handlerInfo.TypeOfArgsElem)
	if handlerInfo.IndexOfRequestBody != -1 {
		bodyValue, err := reqExec.GetFormValue(handlerInfo, r)
		if err != nil {
			return nil, err
		}
		if args[handlerInfo.IndexOfRequestBody].Kind() == reflect.Ptr {
			args[handlerInfo.IndexOfRequestBody] = *bodyValue
		} else {
			args[handlerInfo.IndexOfRequestBody] = (*bodyValue).Elem()
		}

	}
	if handlerInfo.IndexOfAuthClaimsArg != -1 {
		AuthClaimsType := handlerInfo.Method.Type.In(handlerInfo.IndexOfAuthClaimsArg)
		AuthClaimsValue, err := Helper.DepenAuthCreate(AuthClaimsType, r, w)
		if err != nil {
			return nil, err
		}
		if AuthClaimsType.Kind() == reflect.Ptr {
			args[handlerInfo.IndexOfAuthClaimsArg] = *AuthClaimsValue
		} else {
			args[handlerInfo.IndexOfAuthClaimsArg] = (*AuthClaimsValue).Elem()
		}

	}
	err = reqExec.LoadInjectorInjectServiceToArgs(handlerInfo, r, w, args)
	if err != nil {
		return nil, err
	}
	err = reqExec.LoadInjectorsToArgs(handlerInfo, r, w, args)
	if err != nil {
		return nil, err
	}
	//reqExec.CreateHandler(handlerInfo)
	rets := handlerInfo.Method.Func.Call(args)
	if len(rets) > 0 {
		if err, ok := rets[len(rets)-1].Interface().(error); ok {
			return nil, err
		}
	}
	return rets[0].Interface(), nil

}
