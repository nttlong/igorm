package handlers

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"wx/internal"
	"wx/services"
)

func (h *helperType) FindControllerName(reiverType reflect.Type) string {
	if reiverType.Kind() == reflect.Ptr {
		reiverType = reiverType.Elem()
	}
	if reiverType.Kind() != reflect.Struct {
		return ""
	}
	key := reiverType.String() + "/FindControllerName"
	ret, _ := internal.OnceCall(key, func() (*string, error) {

		for i := 0; i < reiverType.NumField(); i++ {
			field := reiverType.Field(i)
			tags := field.Tag.Get("controller")
			if tags != "" {
				tags = h.ToKebabCase(tags)

				return &tags, nil
			}
		}

		items := strings.Split(reiverType.String(), ".")
		ret := h.ToKebabCase(items[len(items)-1])
		return &ret, nil
		/* find first posistion of  "/controllers/" */

	})
	return *ret
}
func (h *helperType) calculateUrlWithQuery(ret *HandlerInfo) {
	ret.QueryParams = []QueryParam{}

	uri := strings.TrimSuffix(strings.Split(ret.Uri, "?")[0], "/")
	ret.UriQuery = strings.Split(ret.Uri, "?")[1]
	ret.Uri = uri

	//ret.UriHandler = strings.TrimSuffix(strings.Split(uri, "?")[0], "/")
	items := strings.Split(ret.UriQuery, "&")
	for _, x := range items {
		fieldName := strings.Split(x, "=")[1]
		fieldName = strings.TrimPrefix(fieldName, "{")
		fieldName = strings.TrimSuffix(fieldName, "}")
		field, ok := ret.TypeOfArgsElem.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, fieldName)
		})
		if !ok {
			continue
		}
		ret.QueryParams = append(ret.QueryParams, QueryParam{
			Name:       fieldName,
			FieldIndex: field.Index,
		})
	}

}
func (h *helperType) calculateUrl(ret *HandlerInfo) {
	if len(ret.UriParams) > 0 {
		if !strings.Contains(ret.Uri, "{*") {
			ret.RegexUri = h.TemplateToRegex(ret.Uri)
			ret.UriHandler = strings.Split(ret.Uri, "{")[0]
		} else {
			ret.RegexUri = h.convertUrlToRegex(ret.Uri)
			ret.UriHandler = strings.Split(ret.Uri, "{")[0]
		}

		ret.IsRegexHandler = true

	} else {
		ret.RegexUri = h.EscapeSpecialCharsForRegex(ret.Uri)
		if ret.IsRegexHandler {
			ret.UriHandler = ret.Uri + "/"
		} else {
			ret.UriHandler = ret.Uri
		}
	}
}

type initGetHandlerInfo struct {
	val  *HandlerInfo
	err  error
	once sync.Once
}

var cacheGetHandlerInfo sync.Map

func (h *helperType) GetHandlerInfo(method reflect.Method) (*HandlerInfo, error) {
	key := method
	ret, _ := cacheGetHandlerInfo.LoadOrStore(key, &initGetHandlerInfo{})
	init := ret.(*initGetHandlerInfo)
	init.once.Do(func() {
		init.val, init.err = h.getHandlerInfo(method)
	})
	return init.val, init.err

}
func (h *helperType) getHandlerInfo(method reflect.Method) (*HandlerInfo, error) {

	ret := &HandlerInfo{
		IndexOfRequestBody: -1,
	}

	ret.IndexOfInjectors = []int{}
	for i := 1; i < method.Type.NumIn(); i++ {
		typ := method.Type.In(i)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct {
			continue
		}
		if services.ServiceUtils.IsInjector(typ) {
			ret.IndexOfInjectors = append(ret.IndexOfInjectors, i)
			ret.HasInjector = true
			break
		}
	}
	/*
		find an arg handler
	*/
	isHandlerMethod := false
	for i := 1; i < method.Type.NumIn(); i++ {
		if !h.Iscontains(ret.IndexOfInjectors, i) {
			fieldIndex, err := h.FindHandlerFieldIndexFormType(method.Func.Type().In(i))
			if err != nil {
				return nil, err
			}
			if fieldIndex != nil {
				ret.IndexOfArg = i
				ret.FieldIndex = fieldIndex
				isHandlerMethod = true
				break

			}

		}
	}
	if !isHandlerMethod {
		return nil, nil
	}
	ret.IndexOfAuthClaimsArg = -1
	ret.IndexOfAuthClaims = nil

	for i := 0; i < method.Type.NumIn(); i++ {
		if h.Iscontains(ret.IndexOfAuthClaims, i) {
			continue
		}
		indexOfAuthClaimsField := h.GetAuthClaims(method.Type.In(i))
		if indexOfAuthClaimsField != nil {
			ret.IndexOfAuthClaimsArg = i
			ret.IndexOfAuthClaims = indexOfAuthClaimsField
			break
		}
	}

	ret.ReceiverIndex = 0
	ret.ReceiverType = method.Type.In(0)
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
			controllerName := h.FindControllerName(ret.ReceiverTypeElem)
			if ret.Uri != "" && ret.Uri[0] == '/' {
				ret.Uri = strings.Replace(ret.Uri, "@", h.ToKebabCase(method.Name), 1)
			} else {
				typ := method.Type.In(0)
				if typ.Kind() == reflect.Ptr {
					typ = typ.Elem()
				}

				ret.Uri = strings.Replace(ret.Uri, "@", controllerName+"/"+h.ToKebabCase(method.Name), 1)

			}
		} else {
			controllerName := h.FindControllerName(ret.ReceiverTypeElem)
			if ret.Uri == "" {
				ret.Uri = controllerName + "/" + h.ToKebabCase(method.Name)
			} else {
				if ret.Uri[0] == '/' {
					ret.IsAbsUri = true
					ret.Uri = ret.Uri[1:]
				}
				if strings.Contains(ret.Uri, "@") {
					ret.Uri = strings.Replace(ret.Uri, "@", controllerName, 1)
				} else {
					ret.Uri = controllerName + "/" + ret.Uri
				}
				if ret.IsAbsUri {
					ret.Uri = "/" + ret.Uri
				}

			}

		}

		ret.UriParams = h.ExtractUriParams(ret.Uri)
		if strings.Contains(ret.Uri, "?") {
			ret.IsQueryUri = true
		}
		if ret.IsQueryUri {
			h.calculateUrlWithQuery(ret)
		}
		h.calculateUrl(ret)
		if ret.IsQueryUri {
			ret.Uri = ret.Uri + "?" + ret.UriQuery
		}

	}
	if ret.IndexOfRequestBody == -1 {
		for i := 1; i < method.Type.NumIn(); i++ {
			if !h.Iscontains(ret.IndexOfInjectors, i) && i != ret.IndexOfArg && i != ret.IndexOfAuthClaimsArg {
				typ := method.Type.In(i)
				ret.TypeOfRequestBody = typ
				if typ.Kind() == reflect.Ptr {
					typ = typ.Elem()
				}
				if typ.Kind() != reflect.Struct {
					continue
				}

				ret.TypeOfRequestBodyElem = typ
				ret.IndexOfRequestBody = i
				break
			}
		}
	}
	if ret.IndexOfRequestBody != -1 {

		ret.FormUploadFile = h.FindFormUploadInType(ret.TypeOfRequestBodyElem)

	}

	for i := range ret.UriParams {
		if field, ok := ret.TypeOfArgsElem.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, ret.UriParams[i].Name)
		}); ok {
			ret.UriParams[i].FieldIndex = field.Index

		}

	}
	if ret.Uri != "" && ret.Uri[0] == '/' {
		ret.IsAbsUri = true
	}

	// // Thay thế tất cả "/handlers/" thành "/"
	// ret.UriHandler = h.trimHandlers(ret.UriHandler)
	// ret.Uri = h.trimHandlers(ret.Uri)
	// ret.RegexUri = h.trimHandlersRegex(ret.RegexUri)
	// if !ret.IsRegexHandler {
	// 	ret.Uri = strings.TrimSuffix(ret.Uri, "/")
	// 	ret.UriHandler = strings.TrimSuffix(ret.UriHandler, "/")
	// }

	return ret, nil
}
func (h *helperType) delete_trimHandlers(s string) string {
	// isAbs := false
	// if s != "" && s[0] == '/' {
	// 	isAbs = true
	// }
	re := regexp.MustCompile(`/handlers/`)
	s = re.ReplaceAllString("/"+s, "/")
	// s = strings.ReplaceAll(s, "//", "/")
	// if isAbs {
	// 	if s == "" {
	// 		s = "/"
	// 	} else {
	// 		s = "/" + s
	// 	}
	// }
	// if s[len(s)-1] != '/' {
	// 	s = s + "/"
	// }
	return s

}
func (h *helperType) delete_trimHandlersRegex(s string) string {
	isAbs := false
	if s != "" && s[0] == '/' {
		isAbs = true
	}
	re := regexp.MustCompile(`\/\^handlers\/`)
	s = re.ReplaceAllString("/"+s, "/")
	re2 := regexp.MustCompile(`\/handlers\/`)
	s = re2.ReplaceAllString("/"+s, "/")
	if !isAbs {
		for strings.HasPrefix(s, "/") {
			s = strings.TrimPrefix(s, "/")
		}
	}
	// s = strings.ReplaceAll(s, "//", "\\/")
	// if isAbs {
	// 	if s == "" {
	// 		s = "\\/"
	// 	} else {
	// 		s = "\\/" + s
	// 	}
	// }
	return s

}
func (h *helperType) GetReceiverTypeFromMethod(method reflect.Method) (*reflect.Type, error) {
	ret := method.Type.In(0)
	if ret.Kind() == reflect.Ptr {
		return &ret, nil
	}
	return nil, fmt.Errorf("receiver arg of %s is not a point of struct %s", method.Name, ret.String())

}
