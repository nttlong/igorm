package vapi

import (
	"encoding/json"
	"mime"
	"reflect"
	"strings"
	"vapi/swaggers"
)

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
func createSwaggerConsumes(web webHandler) []string {
	ret := []string{}
	if web.apiInfo.FormUploadFile != nil {
		ret = append(ret, "multipart/form-data")
	} else {
		ret = append(ret, "application/json")
	}
	return ret

}
func createSwaggerProduces(web webHandler) []string {
	ret := []string{}
	if strings.Contains(web.apiInfo.Uri, ".") {
		extFile := "." + strings.Split(web.apiInfo.Uri, ".")[1]
		mimeType := mime.TypeByExtension(extFile)
		ret = append(ret, mimeType)
	} else {
		ret = append(ret, "application/json")
	}
	return ret
}
func createSwaggerParameters(web webHandler) []swaggers.Parameter {
	ret := []swaggers.Parameter{}
	if len(web.apiInfo.UriParams) > 0 {
		for _, param := range web.apiInfo.UriParams {
			if strings.Contains(web.apiInfo.Uri, "?") {
				ret = append(ret, swaggers.Parameter{

					Name:     param.Name,
					In:       "query",
					Required: true,
					Type:     "string",
				})

			} else {
				ret = append(ret, swaggers.Parameter{

					Name:     param.Name,
					In:       "path",
					Required: true,
					Type:     "string",
				})
			}
		}
	}
	if web.apiInfo.FormUploadFile != nil {
		for _, index := range web.apiInfo.FormUploadFile {
			field := web.apiInfo.TypeOfRequestBodyElem.Field(index)
			tpy := field.Type
			if tpy.Kind() == reflect.Ptr {
				tpy = tpy.Elem()
			}
			if tpy.Kind() == reflect.Slice {
				ret = append(ret, swaggers.Parameter{

					Name:        field.Name,
					In:          "formData",
					Required:    false,
					Type:        "file",
					Description: "select multiple files",
				})
			} else {
				ret = append(ret, swaggers.Parameter{

					Name:     field.Name,
					In:       "formData",
					Required: true,
					Type:     "file",
				})
			}
		}
		for i := 0; i < web.apiInfo.TypeOfRequestBodyElem.NumField(); i++ {
			if !contains(web.apiInfo.FormUploadFile, i) {
				field := web.apiInfo.TypeOfRequestBodyElem.Field(i)
				typ := field.Type
				isRequire := true
				if typ.Kind() == reflect.Ptr {
					typ = typ.Elem()
					isRequire = false
				}
				if typ.Kind() == reflect.Struct {
					ins := reflect.New(typ).Interface()
					desc, _ := json.MarshalIndent(ins, " ", "  ")

					ret = append(ret, swaggers.Parameter{

						Name:        field.Name,
						In:          "formData",
						Required:    isRequire,
						Type:        "object",
						Description: string(desc),
					})

				} else {
					ret = append(ret, swaggers.Parameter{

						Name:     field.Name,
						In:       "formData",
						Required: isRequire,
						Type:     "string",
					})
				}

			}
		}

	}
	return ret
}
func createSwaggerOperation(web webHandler) *swaggers.Operation {
	//panic("unimplemented, at file packages/fapi/Swagger.createSwaggerOperation.go")
	ret := &swaggers.Operation{
		Consumes: createSwaggerConsumes(web),
		Produces: createSwaggerProduces(web),

		Responses: map[string]swaggers.Response{
			"200": {
				Description: "OK",
			},
			"400": {
				Description: "Bad Request",
			},
			"401": {
				Description: "Unauthorized",
			},
		},
		Parameters: []swaggers.Parameter{},
		Security:   []map[string][]string{},
	}
	if len(web.apiInfo.IndexOfAuthClaims) > 0 {
		ret.Parameters = createSwaggerParameters(web)
		ret.Security = append(ret.Security, map[string][]string{})
		ret.Security[0]["OAuth2Password"] = []string{}
	}
	return ret
}

func LoadHandlerInfo(s *swaggers.Swagger) {

	// panic("unimplemented, at file packages/fapi/Swagger.Load.handlerInfos.go")
	for _, handler := range handlerList {
		op := createSwaggerOperation(handler)
		op.Tags = append(op.Tags, handler.apiInfo.ReceiverTypeElem.String())
		pathItem := swaggers.PathItem{}

		pathItemVal := reflect.ValueOf(&pathItem).Elem()

		fieldInfo := pathItemVal.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, handler.method)

		})
		fieldInfo.Set(reflect.ValueOf(op))
		if handler.apiInfo.Uri[0] == '/' {
			s.Paths[handler.apiInfo.Uri] = pathItem
		} else {
			s.Paths["/"+handler.apiInfo.Uri] = pathItem
		}

		// fmt.Println(handler.SwaggerRoute)

	}
}
