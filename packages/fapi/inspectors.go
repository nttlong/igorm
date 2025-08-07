package fapi

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"
)

type Injector struct {
}
type inspectors struct {
}
type inspectorRoute struct {
	uri                string
	path               string
	method             string
	regexUri           string
	requestContentType string
	indexOfFieldInUri  [][]int
	uriParams          []string
}
type inspectorMethod struct {
	/*
		Index of any arg is Context or embeded by Context
	*/
	IndexOfContext   int
	IndexOfPostData  int
	postDataType     reflect.Type
	IndexOfData      int
	IndexOfInjectors []int
	Method           reflect.Method
	tags             *string
	route            *inspectorRoute
}

func (route *inspectorRoute) UriParams() []string {
	if route.uriParams != nil {
		return route.uriParams
	}
	route.uriParams = []string{}
	items := strings.Split(route.uri, "{")
	for _, item := range items {
		if strings.Contains(item, "}") {
			route.uriParams = append(route.uriParams, strings.Split(item, "}")[0])
		}
	}
	return route.uriParams

}
func (route *inspectorRoute) DetectUri(typ reflect.Type, parentIndex []int) [][]int {
	ret := [][]int{}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		fmt.Println(typ.Name())
		return ret
	}
	for _, param := range route.UriParams() {
		if field, ok := typ.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, param)
		}); ok {
			ret = append(ret, append(parentIndex, field.Index...))
		}
	}

	return ret

}
func (m *inspectorMethod) PostDataType() reflect.Type {
	if m.IndexOfData == -1 {
		return nil
	}
	if m.postDataType == nil {
		m.postDataType = m.Method.Type.In(m.IndexOfData)
		if m.postDataType.Kind() == reflect.Ptr {
			m.postDataType = m.postDataType.Elem()
		}
	}
	return m.postDataType

}
func replaceSpecialCharInRegex(s string) string {
	specialChar := "/\\?.$%^*"
	ret := ""
	for _, x := range s {
		if strings.Contains(specialChar, string(x)) {
			ret += "\\"
		}
		ret += string(x)

	}
	return ret
}

func (m *inspectorMethod) Route() inspectorRoute {
	if m.route != nil {
		return *m.route
	}
	tag := m.Tags()
	if tag == "" {
		m.route = &inspectorRoute{
			uri:                strings.ToLower(m.Method.Name),
			method:             "POST",
			regexUri:           strings.ToLower(m.Method.Name),
			requestContentType: "application/json",
		}
	} else {
		uri := ""
		method := "POST"
		requestContentType := "application/json"
		tags := strings.Split(tag, ";")
		for _, tag := range tags {
			if strings.HasPrefix(tag, "uri:") {
				uri = strings.TrimPrefix(tag, "uri:")
				uri = strings.ToLower(uri)
				uri = strings.ReplaceAll(uri, "@", strings.ToLower(m.Method.Name))

			}
			if strings.HasPrefix(tag, "method:") {
				method = strings.ToUpper(strings.TrimPrefix(tag, "method:"))
			}

		}
		for i := 0; i < m.PostDataType().NumField(); i++ {
			typ := m.PostDataType().Field(i).Type
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			if typ == reflect.TypeOf(multipart.FileHeader{}) {
				method = "POST"
				requestContentType = "application/x-www-form-urlencoded"

			}

		}
		m.route = &inspectorRoute{
			uri:                uri,
			method:             method,
			regexUri:           strings.ToLower(m.Method.Name),
			requestContentType: requestContentType,
		}

	}
	m.route.indexOfFieldInUri = [][]int{}
	if strings.Contains(m.route.uri, "{") {
		m.route.path = strings.Split(m.route.uri, "{")[0] + "/"
	} else {
		m.route.path = m.route.uri
	}
	retExgUri := replaceSpecialCharInRegex(m.route.uri)

	for _, paramName := range m.route.UriParams() {
		if field, ok := m.PostDataType().FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, paramName)
		}); ok {
			retExgUri = strings.ReplaceAll(retExgUri, "{"+paramName+"}", "(.*)")
			fieldType := field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
			if fieldType.Kind() == reflect.Struct {
				detectResult := m.route.DetectUri(fieldType, field.Index)
				for _, x := range detectResult {
					m.route.indexOfFieldInUri = append(m.route.indexOfFieldInUri, x)
				}
			} else {
				m.route.indexOfFieldInUri = append(m.route.indexOfFieldInUri, field.Index)
			}

		}
		// m.route.indexOfFieldInUri = append(m.route.indexOfFieldInUri, m.Route().DetectUri(m.PostDataType(), nil)...)
	}
	m.route.regexUri = retExgUri
	return *m.route

}
func (m *inspectorMethod) getTagsFormType(typ reflect.Type) string {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return ""
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == "_" {
			t := field.Tag.Get("route")
			if t != "" {
				return t
			}
		}
		if field.Anonymous {
			return m.getTagsFormType(field.Type)
		}
	}
	return ""

}
func (m *inspectorMethod) Tags() string {
	if m.tags != nil {
		return *m.tags
	}
	if m.PostDataType() == nil {
		t := ""
		m.tags = &t

	} else {
		for i := 0; i < m.PostDataType().NumField(); i++ {
			field := m.PostDataType().Field(i)
			if field.Name == "_" {
				t := field.Tag.Get("route")
				if t != "" {
					m.tags = &t
					return *m.tags
				}
			} else {
				t := m.getTagsFormType(field.Type)
				if t != "" {
					m.tags = &t
					return *m.tags
				}
			}

		}
	}
	if m.tags == nil {
		t := ""
		m.tags = &t
	}
	return *m.tags
}
func (inspector *inspectors) IsPostData(currentIndexOfArgs int, typ reflect.Type) int {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return -1
	}

	return currentIndexOfArgs
}

func (inspector *inspectors) IsInjectorType(currentIndexOfArgs int, typ reflect.Type) int {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if reflect.TypeOf(Injector{}) == typ {
		return currentIndexOfArgs
	}
	return -1
}

func (inspector *inspectors) CreateFormTypeOfArgs(currentIndexOfContext int, typ reflect.Type) *inspectorMethod {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	ret := &inspectorMethod{
		IndexOfContext:   -1,
		IndexOfInjectors: []int{},
	}

	if reflect.TypeOf(Context{}) == typ {
		ret.IndexOfContext = currentIndexOfContext
		return ret
	}
	for i := 0; i < typ.NumField(); i++ {
		if ret := inspector.CreateFormTypeOfArgs(currentIndexOfContext, typ.Field(i).Type); ret != nil {
			return ret
		}
	}
	return nil

}
func (inspector *inspectors) ContextMethod(mt reflect.Method) *inspectorMethod {
	var ret *inspectorMethod
	IndexOfInjectors := []int{}
	IndexOfPostData := -1

	for i := 1; i < mt.Type.NumIn(); i++ {
		if ret == nil {
			ret = inspector.CreateFormTypeOfArgs(i, mt.Type.In(i))
			if ret != nil {
				continue
			}
		}
		if retInspect := inspector.IsInjectorType(i, mt.Type.In(i)); retInspect != -1 {
			IndexOfInjectors = append(IndexOfInjectors, retInspect)
			continue
		}
		if retInspect := inspector.IsPostData(i, mt.Type.In(i)); retInspect != -1 {
			IndexOfPostData = retInspect
			continue
		}

	}
	if ret == nil {
		return nil
	}
	ret.Method = mt
	ret.IndexOfInjectors = IndexOfInjectors
	ret.IndexOfData = IndexOfPostData

	return ret
}
func (inspector *inspectors) InspectorMethod(mt reflect.Method) *inspectorMethod {
	ret := inspector.ContextMethod(mt)
	ret.Route() //<-- init route
	if ret == nil {
		return nil
	}
	return ret

}

var inspector = &inspectors{}
