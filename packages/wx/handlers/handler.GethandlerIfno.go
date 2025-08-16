package handlers

import (
	"reflect"
	"strings"
	"wx/services"
)

func (h *helperType) GetHandlerInfo(method reflect.Method) (*HandlerInfo, error) {

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
		if indexOfAuthClaimsField := h.GetAuthClaims(method.Type.In(i)); indexOfAuthClaimsField != nil {
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
			if ret.Uri != "" && ret.Uri[0] == '/' {
				ret.Uri = strings.Replace(ret.Uri, "@", h.ToKebabCase(method.Name), 1)
			} else {
				typ := method.Type.In(0)
				if typ.Kind() == reflect.Ptr {
					typ = typ.Elem()
				}
				fullName := typ.String()
				items := strings.Split(fullName, ".")
				items = append(items, h.ToKebabCase(method.Name))

				for i := range items {
					items[i] = h.ToKebabCase(items[i])
				}

				ret.Uri = strings.Replace(ret.Uri, "@", strings.Join(items, "/"), 1)

			}
		} else {
			if ret.Uri == "" {
				receiverTypeStr := ret.ReceiverTypeElem.String()
				items := strings.Split(receiverTypeStr, ".")
				for i := range items {
					items[i] = h.ToKebabCase(items[i])
				}
				items = append(items, h.ToKebabCase(method.Name))
				ret.Uri = strings.Join(items, "/")
			} else {
				ret.Uri = ret.Uri + "/" + h.ToKebabCase(method.Name)
			}

		}

		ret.UriParams = h.ExtractUriParams(ret.Uri)
		if len(ret.UriParams) > 0 {
			ret.RegexUri = h.TemplateToRegex(ret.Uri)
			ret.UriHandler = strings.Split(ret.Uri, "{")[0]
			ret.IsRegexHandler = true

		} else {
			ret.RegexUri = h.EscapeSpecialCharsForRegex(ret.Uri)
			ret.UriHandler = ret.Uri + "/"
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

		// ret.TypeOfRequestBodyElem = ret.TypeOfRequestBody
		// if ret.TypeOfRequestBody.Kind() == reflect.Ptr {
		// 	ret.TypeOfRequestBodyElem = ret.TypeOfRequestBody.Elem()
		// }

		ret.FormUploadFile = h.FindFormUploadInType(ret.TypeOfRequestBodyElem)

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
