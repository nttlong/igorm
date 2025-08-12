package vapi

import (
	"reflect"
	"strings"
	swaggers3 "vapi/swagger3"
)

func (sb *SwaggerBuild) swagger3GetPaths() *SwaggerBuild {
	ret := map[string]swaggers3.PathItem{}

	for _, h := range handlerList {
		swaggerUri := strings.TrimPrefix(h.apiInfo.Uri, "/")

		pathItem := swaggers3.PathItem{}
		pathItemType := reflect.TypeOf(pathItem)

		fieldHttpMethod, ok := pathItemType.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, h.method)
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

func (sb *SwaggerBuild) createOperation(handler webHandler) *swaggers3.Operation {
	ret := &swaggers3.Operation{
		Tags:       []string{handler.apiInfo.ReceiverTypeElem.String()},
		Parameters: sb.createParamtersFromUriParams(handler),
		Responses:  map[string]swaggers3.Response{},
	}
	if len(handler.apiInfo.FormUploadFile) > 0 {
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
	if handler.apiInfo.IndexOfRequestBody > 0 {

		ret.Parameters = append(ret.Parameters, sb.createBodyParameters(handler))
	}
	return ret
}
func (sb *SwaggerBuild) createRequestBodyForUploadFile(handler webHandler) *swaggers3.RequestBody {
	if len(handler.apiInfo.FormUploadFile) > 0 {
		props := make(map[string]*swaggers3.Schema)

		for _, index := range handler.apiInfo.FormUploadFile {
			field := handler.apiInfo.TypeOfRequestBodyElem.Field(index)
			tpy := field.Type
			if tpy.Kind() == reflect.Ptr {
				tpy = tpy.Elem()
			}

			if tpy.Kind() == reflect.Slice {
				// multiple files
				props[field.Name] = &swaggers3.Schema{
					Type: "array",
					Items: &swaggers3.Schema{
						Type:   "string",
						Format: "binary",
					},
					Description: "select multiple files",
				}
			} else {
				// single file
				props[field.Name] = &swaggers3.Schema{
					Type:   "string",
					Format: "binary",
				}
			}
		}
		for i := 0; i < handler.apiInfo.TypeOfRequestBodyElem.NumField(); i++ {
			if !contains(handler.apiInfo.FormUploadFile, i) {
				field := handler.apiInfo.TypeOfRequestBodyElem.Field(i)
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

					props[field.Name] = &swaggers3.Schema{
						Type: "array",
						Items: &swaggers3.Schema{
							Type: strType,
						},
					}
					continue
				}
				if fieldType.Kind() == reflect.Struct {
					strType = "object"
				}
				props[field.Name] = &swaggers3.Schema{
					Type: strType,
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

func (sb *SwaggerBuild) createParamtersFromUriParams(handler webHandler) []swaggers3.Parameter {
	ret := []swaggers3.Parameter{}
	if len(handler.apiInfo.UriParams) > 0 {
		for _, param := range handler.apiInfo.UriParams {
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
func (sb *SwaggerBuild) createBodyParameters(handler webHandler) swaggers3.Parameter {
	Example := reflect.New(handler.apiInfo.TypeOfRequestBodyElem).Interface()
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
