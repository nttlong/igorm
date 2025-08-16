package htttpserver

import (
	"reflect"
	"strings"
	"wx/internal"
	swaggers3 "wx/swagger3"
)

func (sb *SwaggerBuild) swagger3GetPaths() *SwaggerBuild {
	ret := map[string]swaggers3.PathItem{}

	for _, h := range HandlerList {
		swaggerUri := strings.TrimPrefix(h.ApiInfo.Uri, "/")

		pathItem := swaggers3.PathItem{}
		pathItemType := reflect.TypeOf(pathItem)

		fieldHttpMethod, ok := pathItemType.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, h.Method)
		})
		if !ok {
			continue
		}

		operation := sb.createOperation(h)
		operationValue := reflect.ValueOf(operation)

		pathItemValue := reflect.ValueOf(&pathItem).Elem() // lấy địa chỉ struct để set

		fieldValue := pathItemValue.FieldByIndex(fieldHttpMethod.Index)
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue.Set(operationValue) // <<--panic: reflect.Set: value of type swaggers3.Operation is not assignable to type *swaggers3.Operation

		} else {
			fieldValue.Set(operationValue.Elem())
		}

		ret["/"+swaggerUri] = pathItem
	}

	sb.swagger.Paths = ret
	return sb
}

func (sb *SwaggerBuild) createOperation(handler WebHandler) *swaggers3.Operation {
	var content map[string]swaggers3.MediaType
	// errType := reflect.TypeOf((*error)(nil)).Elem()
	content = map[string]swaggers3.MediaType{
		"text/plain": {
			Schema: &swaggers3.Schema{
				Type: "string",
			},
		},
	}
	if handler.Method == "POST" {
		content = map[string]swaggers3.MediaType{
			"application/json": {
				Schema: &swaggers3.Schema{
					Type: "object",
				},
			},
		}
	}

	ret := &swaggers3.Operation{
		Tags:       []string{handler.ApiInfo.ReceiverTypeElem.String()},
		Parameters: sb.createParamtersFromUriParams(handler),
		Responses: map[string]swaggers3.Response{
			"200": {
				Description: "OK",
				Content:     content,
			},
			"206": {
				Description: "Partial Content",
				Content:     content,
			},
		},
	}
	if len(handler.ApiInfo.FormUploadFile) > 0 {
		/*
					"requestBody": {
			        "required": true,
			        "content": {
			          "multipart/form-data": {
			            "schema": {
			              "type": "object",
			              "properties": {
			                "Files": {
			                  "type": "array",
			                  "items": {
			                    "type": "string",
			                    "format": "binary"
			                  }
			                }
			              }
			            }
			          }
		*/

		ret.RequestBody = sb.createRequestBodyForUploadFile(handler)
		sb.applySecurity(handler, ret)
		return ret

	}
	if handler.ApiInfo.IndexOfRequestBody > 0 {

		ret.Parameters = append(ret.Parameters, sb.createBodyParameters(handler))

	}
	sb.applySecurity(handler, ret)
	return ret
}
func (sb *SwaggerBuild) createRequestBodyForUploadFile(handler WebHandler) *swaggers3.RequestBody {
	if len(handler.ApiInfo.FormUploadFile) > 0 {
		props := make(map[string]*swaggers3.Schema)

		for _, index := range handler.ApiInfo.FormUploadFile {
			field := handler.ApiInfo.TypeOfRequestBodyElem.Field(index)
			typ := field.Type
			arrayNullable := false
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
				arrayNullable = true

			}

			if typ.Kind() == reflect.Slice {
				// multiple files

				props[field.Name] = &swaggers3.Schema{

					Type:     "array",
					Nullable: arrayNullable,
					Items: &swaggers3.Schema{
						Type:     "string",
						Format:   "binary",
						Nullable: typ.Kind() == reflect.Ptr,
					},
					Description: "select multiple files",
				}
			} else {
				// single file
				props[field.Name] = &swaggers3.Schema{
					Type:     "string",
					Format:   "binary",
					Nullable: arrayNullable,
				}
			}
		}
		for i := 0; i < handler.ApiInfo.TypeOfRequestBodyElem.NumField(); i++ {
			if !internal.Contains(handler.ApiInfo.FormUploadFile, i) {
				field := handler.ApiInfo.TypeOfRequestBodyElem.Field(i)
				fieldType := field.Type
				if fieldType.Kind() == reflect.Ptr {
					fieldType = fieldType.Elem()
				}
				strType := "string"
				if fieldType.Kind() == reflect.Slice {
					strType = "array"
					eleType := fieldType.Elem()
					if eleType.Kind() == reflect.Ptr {
						eleType = eleType.Elem()
					}
					if eleType.Kind() == reflect.Struct {
						strType = "object"
					}
					example := reflect.New(eleType).Interface()
					props[field.Name] = &swaggers3.Schema{
						Type: "array",
						Items: &swaggers3.Schema{
							Type:    strType,
							Example: example,
						},
					}
					continue
				}
				if fieldType.Kind() == reflect.Struct {
					strType = "object"
				}
				example := reflect.New(fieldType).Interface()
				props[field.Name] = &swaggers3.Schema{
					Type:    strType,
					Example: example,
				}

			}
		}
		// Gán vào requestBody thay vì parameters
		ret := &swaggers3.RequestBody{
			Required: true,
			Content: map[string]swaggers3.MediaType{
				"multipart/form-data": {

					Schema: &swaggers3.Schema{
						Type:       "object",
						Properties: props,
					},
				},
			},
		}
		return ret
	}
	return nil

}

func (sb *SwaggerBuild) createParamtersFromUriParams(handler WebHandler) []swaggers3.Parameter {
	ret := []swaggers3.Parameter{}
	if len(handler.ApiInfo.UriParams) > 0 {
		for _, param := range handler.ApiInfo.UriParams {
			ret = append(ret, swaggers3.Parameter{
				Name:     param.Name,
				In:       "path",
				Required: true,
				Schema: &swaggers3.Schema{
					Type: "string",
				},
			})
		}
	}

	return ret

}
func (sb *SwaggerBuild) createBodyParameters(handler WebHandler) swaggers3.Parameter {
	Example := reflect.New(handler.ApiInfo.TypeOfRequestBodyElem).Interface()
	ret := swaggers3.Parameter{
		Name:     "body",
		In:       "body",
		Required: true,
		Schema: &swaggers3.Schema{
			Type: "object",
		},
		Example: Example,
	}
	return ret
}

package htttpserver

import (
	"reflect"
	"strings"
	"wx/internal"
	swaggers3 "wx/swagger3"
)

func (sb *SwaggerBuild) swagger3GetPaths() *SwaggerBuild {
	ret := map[string]swaggers3.PathItem{}

	for _, h := range HandlerList {
		swaggerUri := strings.TrimPrefix(h.ApiInfo.Uri, "/")

		pathItem := swaggers3.PathItem{}
		pathItemType := reflect.TypeOf(pathItem)

		fieldHttpMethod, ok := pathItemType.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, h.Method)
		})
		if !ok {
			continue
		}

		operation := sb.createOperation(h)
		operationValue := reflect.ValueOf(operation)

		pathItemValue := reflect.ValueOf(&pathItem).Elem() // lấy địa chỉ struct để set

		fieldValue := pathItemValue.FieldByIndex(fieldHttpMethod.Index)
		if fieldValue.Kind() == reflect.Ptr {
			fieldValue.Set(operationValue) // <<--panic: reflect.Set: value of type swaggers3.Operation is not assignable to type *swaggers3.Operation

		} else {
			fieldValue.Set(operationValue.Elem())
		}

		ret["/"+swaggerUri] = pathItem
	}

	sb.swagger.Paths = ret
	return sb
}

func (sb *SwaggerBuild) createOperation(handler WebHandler) *swaggers3.Operation {
	var content map[string]swaggers3.MediaType
	// errType := reflect.TypeOf((*error)(nil)).Elem()
	content = map[string]swaggers3.MediaType{
		"text/plain": {
			Schema: &swaggers3.Schema{
				Type: "string",
			},
		},
	}
	if handler.Method == "POST" {
		content = map[string]swaggers3.MediaType{
			"application/json": {
				Schema: &swaggers3.Schema{
					Type: "object",
				},
			},
		}
	}

	ret := &swaggers3.Operation{
		Tags:       []string{handler.ApiInfo.ReceiverTypeElem.String()},
		Parameters: sb.createParamtersFromUriParams(handler),
		Responses: map[string]swaggers3.Response{
			"200": {
				Description: "OK",
				Content:     content,
			},
			"206": {
				Description: "Partial Content",
				Content:     content,
			},
		},
	}
	if len(handler.ApiInfo.FormUploadFile) > 0 {
		/*
					"requestBody": {
			        "required": true,
			        "content": {
			          "multipart/form-data": {
			            "schema": {
			              "type": "object",
			              "properties": {
			                "Files": {
			                  "type": "array",
			                  "items": {
			                    "type": "string",
			                    "format": "binary"
			                  }
			                }
			              }
			            }
			          }
		*/

		ret.RequestBody = sb.createRequestBodyForUploadFile(handler)
		return ret

	}
	if handler.ApiInfo.IndexOfRequestBody > 0 {

		ret.Parameters = append(ret.Parameters, sb.createBodyParameters(handler))
	}
	return ret
}
func (sb *SwaggerBuild) createRequestBodyForUploadFile(handler WebHandler) *swaggers3.RequestBody {
	if len(handler.ApiInfo.FormUploadFile) > 0 {
		props := make(map[string]*swaggers3.Schema)

		for _, index := range handler.ApiInfo.FormUploadFile {
			field := handler.ApiInfo.TypeOfRequestBodyElem.Field(index)
			typ := field.Type
			arrayNullable := false
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
				arrayNullable = true

			}

			if typ.Kind() == reflect.Slice {
				// multiple files

				props[field.Name] = &swaggers3.Schema{

					Type:     "array",
					Nullable: arrayNullable,
					Items: &swaggers3.Schema{
						Type:     "string",
						Format:   "binary",
						Nullable: typ.Kind() == reflect.Ptr,
					},
					Description: "select multiple files",
				}
			} else {
				// single file
				props[field.Name] = &swaggers3.Schema{
					Type:     "string",
					Format:   "binary",
					Nullable: arrayNullable,
				}
			}
		}
		for i := 0; i < handler.ApiInfo.TypeOfRequestBodyElem.NumField(); i++ {
			if !internal.Contains(handler.ApiInfo.FormUploadFile, i) {
				field := handler.ApiInfo.TypeOfRequestBodyElem.Field(i)
				fieldType := field.Type
				if fieldType.Kind() == reflect.Ptr {
					fieldType = fieldType.Elem()
				}
				strType := "string"
				if fieldType.Kind() == reflect.Slice {
					strType = "array"
					eleType := fieldType.Elem()
					if eleType.Kind() == reflect.Ptr {
						eleType = eleType.Elem()
					}
					if eleType.Kind() == reflect.Struct {
						strType = "object"
					}
					example := reflect.New(eleType).Interface()
					props[field.Name] = &swaggers3.Schema{
						Type: "array",
						Items: &swaggers3.Schema{
							Type:    strType,
							Example: example,
						},
					}
					continue
				}
				if fieldType.Kind() == reflect.Struct {
					strType = "object"
				}
				example := reflect.New(fieldType).Interface()
				props[field.Name] = &swaggers3.Schema{
					Type:    strType,
					Example: example,
				}

			}
		}
		// Gán vào requestBody thay vì parameters
		ret := &swaggers3.RequestBody{
			Required: true,
			Content: map[string]swaggers3.MediaType{
				"multipart/form-data": {

					Schema: &swaggers3.Schema{
						Type:       "object",
						Properties: props,
					},
				},
			},
		}
		return ret
	}
	return nil

}

func (sb *SwaggerBuild) createParamtersFromUriParams(handler WebHandler) []swaggers3.Parameter {
	ret := []swaggers3.Parameter{}
	if len(handler.ApiInfo.UriParams) > 0 {
		for _, param := range handler.ApiInfo.UriParams {
			ret = append(ret, swaggers3.Parameter{
				Name:     param.Name,
				In:       "path",
				Required: true,
				Schema: &swaggers3.Schema{
					Type: "string",
				},
			})
		}
	}

	return ret

}
func (sb *SwaggerBuild) createBodyParameters(handler WebHandler) swaggers3.Parameter {
	Example := reflect.New(handler.ApiInfo.TypeOfRequestBodyElem).Interface()
	ret := swaggers3.Parameter{
		Name:     "body",
		In:       "body",
		Required: true,
		Schema: &swaggers3.Schema{
			Type: "object",
		},
		Example: Example,
	}
	return ret
}
