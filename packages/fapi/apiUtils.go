package fapi

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"
)

type apiUtils struct {
	regexSepecialChars string
}

/*
this function parse tags
example : method:get;uri:{fileName}.mp4 or uri:{fileName}.mp4
*/
func (api *apiUtils) ParseTags(tag string) *apiTagsInfo {
	if tag == "" {
		return nil
	}
	ret := &apiTagsInfo{}
	tags := strings.Split(tag, ";")
	for _, t := range tags {
		if t == "" {
			continue
		}
		if strings.Contains(t, "method:") {
			ret.Method = strings.TrimPrefix(t, "method:")
		}
		if strings.Contains(t, "uri:") {
			ret.Url = strings.TrimPrefix(t, "uri:")

		}
		if strings.Contains(t, "description:") {
			ret.Description = strings.TrimPrefix(t, "description:")

		}

	}
	return ret
}

/*
dectect  if typ is Context or if any fields is Context
and also return tag route
*/
func (api *apiUtils) CheckTypeIsContextType(typ reflect.Type) (bool, string) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ == reflect.TypeOf(Context{}) {
		return true, ""
	}
	if typ.Kind() != reflect.Struct {
		return false, ""
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Type == reflect.TypeOf(Context{}) {
			return true, field.Tag.Get("route")
		} else if field.Anonymous {
			if ok, tags := api.CheckTypeIsContextType(field.Type); ok {
				if tags != "" {
					return ok, tags
				}

				return ok, field.Tag.Get("route")
			}
		}
	}
	return false, ""
}
func (api *apiUtils) InspectMethod(method reflect.Method) *apiMethodInfo {
	var ret *apiMethodInfo
	for i := 1; i < method.Type.NumIn(); i++ {

		paramType := method.Type.In(i)

		if paramType.Kind() == reflect.Ptr {
			paramType = paramType.Elem()
		}
		isContextTyp, tags := api.CheckTypeIsContextType(paramType)
		if isContextTyp {
			tagsInfo := api.ParseTags(tags)
			if tagsInfo == nil {
				tagsInfo = &apiTagsInfo{
					Method: "post",
					Url:    "/" + strings.ToLower(method.Name),
				}
			} else {
				if strings.Contains(tagsInfo.Url, "@function@") {
					tagsInfo.Url = strings.ReplaceAll(tagsInfo.Url, "@function@", strings.ToLower(method.Name))
				}
			}

			ret = &apiMethodInfo{ // yes, this method is server for web api hanlder
				method:             method,
				param:              paramType,
				recieverType:       method.Type.In(0),
				indexOfContextInfo: i,
				indexOfInjectors:   []int{},
				routeTags:          tags,
				tags:               *tagsInfo,
			}
			break
		}

	}
	if ret != nil {
		for i := 1; i < method.Type.NumIn(); i++ {
			if i != ret.indexOfContextInfo {
				paramType := method.Type.In(i)
				if paramType.Kind() == reflect.Ptr {
					paramType = paramType.Elem()
				}
				if paramType.PkgPath() == reflect.TypeOf(apiUtils{}).PkgPath() {
					fmt.Println(paramType.String())
					if strings.Contains(paramType.String(), "Inject[") {
						(*ret).indexOfInjectors = append((*ret).indexOfInjectors, i)
					}

				}

			}
		}
	}
	return ret

}
func (api *apiUtils) EscapeRegExp(str string) string {
	ret := ""
	for _, char := range str {
		if strings.Contains(api.regexSepecialChars, string(char)) {
			ret += "\\"
		}

		ret += string(char)
	}
	return ret

}
func (api *apiUtils) IsTypeUploadFile(t reflect.Type) bool {
	iface := reflect.TypeOf((*multipart.File)(nil)).Elem()
	if reflect.PtrTo(t).Implements(iface) {
		return true
	}
	if t == reflect.TypeOf(multipart.FileHeader{}) {
		return true
	}
	if t == reflect.TypeOf(multipart.Reader{}) {
		return true
	}

	return false
}

/*
Check typ hash upload file and return list of fieldIndex has is fileUpload
*/
func (api *apiUtils) CheckHasInputFile(typ reflect.Type) (bool, [][]int) {
	ret := false
	retIndex := [][]int{}
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()

		check, checkIndex := api.CheckHasInputFile(typ)
		if check {
			retIndex = append(retIndex, checkIndex...)
		}
		ret = ret || check
		if ret {
			return ret, retIndex
		}

	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if api.IsTypeUploadFile(typ) {
		// if typ == reflect.TypeOf(multipart.FileHeader{}) || typ == reflect.TypeOf(multipart.File{}) {

		ret = ret || true
		return ret, retIndex

	}
	if typ.Kind() == reflect.Struct {

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			check, checIndex := api.CheckHasInputFile(field.Type)
			if check {
				retIndex = append(retIndex, field.Index)
				for x := range checIndex {
					retIndex[len(retIndex)-1] = append(retIndex[len(retIndex)-1], x)
				}

			}
			ret = ret || check

		}
	}
	return ret, retIndex

}

var api = &apiUtils{
	regexSepecialChars: "^/\\.*$",
}
