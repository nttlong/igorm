package handlers

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	wxErrors "wx/errors"
)

func (reqExec *RequestExecutor) GetFormValueOld(handlerInfo HandlerInfo, r *http.Request) (*reflect.Value, error) {
	var bodyDataRet reflect.Value
	var bodyType reflect.Type
	if handlerInfo.IsFormPost {
		bodyDataRet = reflect.New(handlerInfo.FormPostTypeEle)
		bodyType = handlerInfo.FormPostTypeEle
	} else {
		bodyDataRet = reflect.New(handlerInfo.TypeOfRequestBodyElem)
		bodyType = handlerInfo.TypeOfRequestBodyElem
	}
	bodyData := bodyDataRet.Elem()
	var formData map[string][]string
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary=") {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			return nil, wxErrors.NewFileParseError("error parsing multipart form", err)
		}
		formData = r.MultipartForm.Value
	} else {
		err := r.ParseForm()
		if err != nil {
			return nil, wxErrors.NewFileParseError("error parsing form", err)
		}
		formData = r.Form
	}

	//scan all post files
	if r.MultipartForm != nil && len(r.MultipartForm.File) > 0 {
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {

			return nil, wxErrors.NewFileParseError("error parsing multipart form", err)
		}

		for key, values := range r.MultipartForm.File {
			field := reqExec.GetFieldByName(handlerInfo.TypeOfRequestBodyElem, key)

			if field == nil {
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
	for key, values := range formData {

		field := reqExec.GetFieldByName(bodyType, key)
		if field == nil {
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
	if handlerInfo.IsFormPost {
		retVal := reflect.New(handlerInfo.TypeOfRequestBodyElem)
		retVal.Elem().FieldByName("Data").Set(bodyDataRet)
		return &retVal, nil

	}

	return &bodyDataRet, nil

}
func (reqExec *RequestExecutor) GetFormValue(handlerInfo HandlerInfo, r *http.Request) (*reflect.Value, error) {
	var target reflect.Value
	var targetType reflect.Type
	var ret reflect.Value

	if handlerInfo.IsFormPost {
		ret = reflect.New(handlerInfo.FormPostTypeEle)
		target = ret.Elem()
		targetType = handlerInfo.FormPostTypeEle
	} else {
		ret = reflect.New(handlerInfo.TypeOfRequestBodyElem)

		target = ret.Elem()
		targetType = handlerInfo.TypeOfRequestBodyElem
	}

	var formData map[string][]string
	var files map[string][]*multipart.FileHeader

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// r.ParseMultipartForm(10 << 20)
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			return nil, wxErrors.NewFileParseError("error parsing multipart form", err)
		}
		formData = r.MultipartForm.Value
		files = r.MultipartForm.File //<-- wx tu dong lay file theo kieu nay
	} else {
		if err := r.ParseForm(); err != nil {
			return nil, wxErrors.NewFileParseError("error parsing form", err)
		}
		formData = r.Form
	}

	// set form fields
	for key, values := range formData {
		field := reqExec.GetFieldByName(targetType, key)
		if field == nil || len(values) == 0 {
			continue
		}
		fv := target.FieldByIndex(field.Index)

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(values[0])
		case reflect.Slice:
			if fv.Type().Elem().Kind() == reflect.String {
				fv.Set(reflect.ValueOf(values))
			}
		case reflect.Ptr:
			elemKind := fv.Type().Elem().Kind()
			if elemKind == reflect.String {
				ptr := reflect.New(fv.Type().Elem())
				ptr.Elem().SetString(values[0])
				fv.Set(ptr)
			} else if elemKind == reflect.Struct {
				ptr := reflect.New(fv.Type().Elem())
				if err := json.Unmarshal([]byte(values[0]), ptr.Interface()); err != nil {
					return nil, err
				}
				fv.Set(ptr)
			}
		case reflect.Struct:
			if err := json.Unmarshal([]byte(values[0]), fv.Addr().Interface()); err != nil {
				return nil, err
			}
		}
	}

	// set files
	for key, fhArr := range files {
		field := reqExec.GetFieldByName(targetType, key)
		if field == nil || len(fhArr) == 0 {
			continue
		}
		fv := target.FieldByIndex(field.Index)
		ft := fv.Type()

		switch {
		case ft == reflect.TypeOf(&multipart.FileHeader{}):
			fv.Set(reflect.ValueOf(fhArr[0]))
		case ft == reflect.TypeOf([]*multipart.FileHeader{}):
			fv.Set(reflect.ValueOf(fhArr))
		case ft == reflect.TypeOf([]multipart.FileHeader{}):
			slice := make([]multipart.FileHeader, len(fhArr))
			for i, f := range fhArr {
				slice[i] = *f
			}
			fv.Set(reflect.ValueOf(slice))
		}
	}

	// special wrapper for FormPost
	if handlerInfo.IsFormPost {
		retVal := reflect.New(handlerInfo.TypeOfRequestBodyElem)
		retVal.Elem().FieldByName("Data").Set(target)
		return &retVal, nil
	}

	return &ret, nil
}
