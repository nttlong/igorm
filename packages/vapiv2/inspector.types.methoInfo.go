package vapi

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"
)

type handlerInfo struct {
	IndexOfArg            int
	TypeOfArgs            reflect.Type
	TypeOfArgsElem        reflect.Type
	FieldIndex            []int
	ReceiverIndex         int
	ReceiverType          reflect.Type
	ReceiverTypeElem      reflect.Type
	Method                reflect.Method
	RouteTags             []string
	Uri                   string
	RegexUri              string
	UriHandler            string
	IsRegexhadler         bool
	UriParams             []uriParam
	IndexOfInjectors      []int
	FormUploadFile        []int
	IndexOfRequestBody    int
	TypeOfRequestBody     reflect.Type
	TypeOfRequestBodyElem reflect.Type
	IndexOfAuthClaimsArg  int
	IndexOfAuthClaims     []int
	HttpMethod            string
}

func (h *helperType) FindHandlerFieldIndexFormType(typ reflect.Type) ([]int, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	key := typ.String() + "/helperType/FindHandlerFieldIndexFormType"
	ret, err := OnceCall(key, func() (*[]int, error) {
		ret, err := h.findHandlerFieldIndexFormType(typ)
		return &ret, err
	})
	return *ret, err

}
func (h *helperType) findHandlerFieldIndexFormType(typ reflect.Type) ([]int, error) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() != reflect.Struct {
			continue
		}
		if fieldType == reflect.TypeOf(Handler{}) {
			return []int{i}, nil
		}
		fieldIndex, err := h.findHandlerFieldIndexFormType(fieldType)
		if err != nil {
			return nil, err
		}
		if fieldIndex != nil {
			return append([]int{i}, fieldIndex...), nil
		}
	}
	return nil, nil
}
func (h *helperType) GetAuthClaims(typ reflect.Type) []int {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}
	if typ == reflect.TypeOf(AuthClaims{}) {
		return []int{}
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldType := field.Type
		ret := h.GetAuthClaims(fieldType)
		if ret != nil {
			return append(field.Index, ret...)
		}
	}
	return nil
}
func (h *helperType) GetHandlerInfo(method reflect.Method) (*handlerInfo, error) {
	receiverIndex := 0
	var ret *handlerInfo
	indexOfRequestBody := -1
	indexOfInjectors := []int{}
	IndexOfAuthClaimsArg := -1
	var IndexOfAuthClaims []int

	for i := 0; i < method.Type.NumIn(); i++ {
		typ := method.Type.In(i)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct {
			continue
		}
		if IndexOfAuthClaimsArg == -1 {
			if ret := h.GetAuthClaims(typ); ret != nil {
				IndexOfAuthClaimsArg = i
				IndexOfAuthClaims = ret
			}
		}
		if typ == reflect.TypeOf(Handler{}) {
			if ret == nil {
				ret = &handlerInfo{
					IndexOfArg: i,
					FieldIndex: []int{},
				}
				continue
			}

		}
		fieldIndex, err := h.FindHandlerFieldIndexFormType(typ)
		if err != nil {
			return nil, err
		}
		if fieldIndex != nil {
			if ret == nil {
				ret = &handlerInfo{
					IndexOfArg: i,
					FieldIndex: fieldIndex,
				}
			}
		}
		if h.IsInjector(typ) {
			indexOfInjectors = append(indexOfInjectors, i)
			continue

		} else if typ.Kind() == reflect.Struct {
			indexOfRequestBody = i

		}
	}
	if ret == nil {
		return nil, nil
	}
	ret.ReceiverIndex = receiverIndex
	ret.ReceiverType = method.Type.In(receiverIndex)
	ret.ReceiverTypeElem = ret.ReceiverType
	if ret.ReceiverType.Kind() == reflect.Ptr {
		ret.ReceiverTypeElem = ret.ReceiverType.Elem()
	}

	ret.Method = method
	ret.HttpMethod = "POST" //<-- defualt is POST
	if ret.IndexOfArg > 0 {
		ret.TypeOfArgs = method.Type.In(ret.IndexOfArg)
		ret.TypeOfArgsElem = ret.TypeOfArgs
		if ret.TypeOfArgs.Kind() == reflect.Ptr {
			ret.TypeOfArgsElem = ret.TypeOfArgs.Elem()
		}
		ret.RouteTags = h.ExtractTags(ret.TypeOfArgsElem, ret.FieldIndex)
		ret.Uri = h.ExtractUriFromTags(ret.RouteTags)
		if HttpMethod := h.ExtractHttpMethodFromTags(ret.RouteTags); HttpMethod != "" {
			ret.HttpMethod = HttpMethod
		}

		if strings.Contains(ret.Uri, "@") {
			ret.Uri = strings.Replace(ret.Uri, "@", h.ToKebabCase(method.Name), 1)
		} else {
			ret.Uri = ret.Uri + "/" + h.ToKebabCase(method.Name)
		}

		ret.UriParams = h.ExtractUriParams(ret.Uri)
		if len(ret.UriParams) > 0 {
			ret.RegexUri = h.TemplateToRegex(ret.Uri)
			ret.UriHandler = strings.Split(ret.Uri, "{")[0] + "/"
			ret.IsRegexhadler = true

		} else {
			ret.RegexUri = h.EscapeSpecialCharsForRegex(ret.Uri)
			ret.UriHandler = ret.Uri + "/"
		}
	}

	ret.IndexOfInjectors = indexOfInjectors

	if indexOfRequestBody != -1 && indexOfRequestBody != ret.IndexOfArg {
		ret.IndexOfRequestBody = indexOfRequestBody
		ret.TypeOfRequestBody = method.Type.In(indexOfRequestBody)
		ret.TypeOfRequestBodyElem = ret.TypeOfRequestBody
		ret.FormUploadFile = h.FindFormUploadInType(ret.TypeOfRequestBodyElem)
	}
	if IndexOfAuthClaimsArg != -1 {
		ret.IndexOfAuthClaimsArg = IndexOfAuthClaimsArg
		ret.IndexOfAuthClaims = IndexOfAuthClaims
	}
	for i := range ret.UriParams {
		if field, ok := ret.TypeOfArgsElem.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, ret.UriParams[i].Name)
		}); ok {
			ret.UriParams[i].FieldIndex = field.Index

		}

	}

	return ret, nil
}
func (h *helperType) ExtractTags(typ reflect.Type, fieldIndex []int) []string {

	if len(fieldIndex) == 0 {
		return nil
	}

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	key := typ.String() + "/helperType/ExtractTags"
	for _, x := range fieldIndex {
		key += "/" + fmt.Sprint(x)
	}
	ret, _ := OnceCall(key, func() (*[]string, error) {
		ret := []string{}
		field := typ.FieldByIndex([]int{fieldIndex[0]})

		ret = append(ret, field.Tag.Get("route"))
		fieldType := field.Type

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		subRet := h.ExtractTags(fieldType, fieldIndex[1:])

		ret = append(ret, subRet...)

		return &ret, nil
	})
	return *ret

}
func (h *helperType) ExtractUriFromTags(tags []string) string {
	key := strings.Join(tags, "**") + "/helperType/ExtractUriFromTags"
	ret, _ := OnceCall(key, func() (*string, error) {
		ret := ""
		for i := len(tags) - 1; i >= 0; i-- {
			tag := tags[i]
			if tag == "" {
				continue
			}
			items := strings.Split(tag, ";")
			for _, item := range items {
				if strings.HasPrefix(item, "uri:") {

					val := item[4:]

					if val != "" {
						if strings.Contains(ret, "@") {
							ret = strings.Replace(ret, "@", val, 1)
						} else {
							ret += "/" + val
						}
					}

				}
			}
		}
		ret = strings.TrimPrefix(strings.TrimSuffix(ret, "/"), "/")
		return &ret, nil
	})
	return *ret

}
func (h *helperType) ExtractHttpMethodFromTags(tags []string) string {
	key := strings.Join(tags, "**") + "/helperType/ExtractHttpMethodFromTags"
	ret, _ := OnceCall(key, func() (*string, error) {
		ret := ""
		for i := len(tags) - 1; i >= 0; i-- {
			tag := tags[i]
			if tag == "" {
				continue
			}
			items := strings.Split(tag, ";")
			for _, item := range items {
				if strings.HasPrefix(item, "method:") {
					ret = strings.ToUpper(item[7:])

				}
			}
		}

		return &ret, nil
	})
	return *ret
}
func (h *helperType) IsInjector(typ reflect.Type) bool {
	key := typ.String() + "/helperType/IsInjector"

	ret, _ := OnceCall(key, func() (*bool, error) {
		typ := typ
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if typ.PkgPath() == reflect.TypeOf(Inject[string]{}).PkgPath() {
			if strings.Contains(typ.Name(), "Inject[") {
				ret := true
				return &ret, nil
			}

		}

		ret := false
		return &ret, nil
	})

	return *ret

}
func (h *helperType) FindFormUploadInType(typ reflect.Type) []int {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct && typ.Kind() != reflect.Slice {
		fmt.Println(typ.String())
		return nil
	}
	if typ == reflect.TypeOf([]multipart.FileHeader{}) {
		return []int{}
	}
	if typ == reflect.TypeOf(multipart.FileHeader{}) {
		return []int{}
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		ret := h.FindFormUploadInType(field.Type)
		if ret != nil {
			return append(field.Index, ret...)
		}
	}
	return nil

}

type uriParam struct {
	Position   int
	Name       string
	FieldIndex []int
}

/*
handlerInfoExtractUriParams extracts all substrings enclosed in curly braces '{}'
from the given URI string, along with their positions (based on segments split by '/').

For example:

Given URI: "abc/{word 1}/abc/dbc/{u2}/{u3}"

The function will return a slice of uriParaim structs:
[

	{Position: 1, Name: "word 1"},
	{Position: 4, Name: "u2"},
	{Position: 5, Name: "u3"},

]

Where:
- Position is the zero-based index of the segment in the URI path split by '/'.
- Name is the trimmed string inside the braces '{}'.

@return []uriParaim - a slice containing extracted parameters with their position and name.
*/
func (h *helperType) ExtractUriParams(uri string) []uriParam {
	key := uri + "/helperType/ExtractUriParams"
	ret, _ := OnceCall(key, func() (*[]uriParam, error) {

		params := []uriParam{}
		segments := h.SplitUriSegments(uri)

		for i, segment := range segments {
			// Check if segment contains a URI parameter enclosed in {}
			name := h.ExtractNameInBraces(segment)
			if name != "" {
				params = append(params, uriParam{
					Position: i,
					Name:     name,
				})
			}
		}

		return &params, nil
	})
	return *ret
}

type helperType struct {
	SpecialCharForRegex string
}

// splitUriSegments splits the URI string by '/', ignoring empty segments.
func (h *helperType) SplitUriSegments(uri string) []string {
	var segments []string
	start := 0

	for i := 0; i < len(uri); i++ {
		if uri[i] == '/' {
			if start < i {
				segments = append(segments, uri[start:i])
			}
			start = i + 1
		}
	}
	// append the last segment if any
	if start < len(uri) {
		segments = append(segments, uri[start:])
	}
	return segments
}

// extractNameInBraces extracts the trimmed content inside the first pair of braces '{}' in the segment.
// Returns empty string if no braces found.
func (h *helperType) ExtractNameInBraces(segment string) string {
	start := -1
	end := -1
	for i, ch := range segment {
		if ch == '{' && start == -1 {
			start = i
		} else if ch == '}' && start != -1 {
			end = i
			break
		}
	}
	if start != -1 && end != -1 && end > start+1 {
		name := segment[start+1 : end]
		return h.TrimSpaces(name)
	}
	return ""
}

// trimSpaces trims leading and trailing spaces from a string.
func (h *helperType) TrimSpaces(s string) string {
	start, end := 0, len(s)-1
	for start <= end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end >= start && (s[end] == ' ' || s[end] == '\t') {
		end--
	}
	if start > end {
		return ""
	}
	return s[start : end+1]
}

// templateToRegex chuyển URI template thành regex pattern string
// lấy các giá trị trong {}
func (h *helperType) TemplateToRegex(template string) string {
	key := template + "/helperType/TemplateToRegex"
	ret, _ := OnceCall(key, func() (*string, error) {
		segments := strings.Split(template, "/")
		regexParts := []string{}
		paramCount := 0
		var escapeRegex = h.EscapeSpecialCharsForRegex
		for _, seg := range segments {
			if seg == "" {
				continue
			}

			var sb strings.Builder
			i := 0
			for i < len(seg) {
				start := strings.Index(seg[i:], "{")
				if start == -1 {
					// No more '{', escape remainder
					sb.WriteString(escapeRegex(seg[i:]))
					break
				}

				start += i
				end := strings.Index(seg[start:], "}")
				if end == -1 {
					// No closing brace, treat literally
					sb.WriteString(escapeRegex(seg[i:]))
					break
				}
				end += start

				// Escape static part before {
				if start > i {
					sb.WriteString(escapeRegex(seg[i:start]))
				}

				// Add capture group for parameter inside {}
				sb.WriteString(`([^/]+)`)
				paramCount++

				// Move index past "}"
				i = end + 1
			}

			regexParts = append(regexParts, sb.String())
		}

		// Join parts with '/'
		regexPattern := "^" + strings.Join(regexParts, "/") + "$"
		return &regexPattern, nil
	})
	return *ret
}
