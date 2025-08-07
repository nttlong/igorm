package fapi

import (
	"fapi/swaggers"
	"fmt"
	"reflect"
)

func (h *apiMethodInfo) CreatePathItemFromhandlerInfo() *swaggers.Operation {

	ret := swaggers.Operation{
		Consumes: []string{h.GetRequestContentType()},
		Produces: []string{h.GetResponseContentType()},

		Summary:    h.tags.Description,
		Parameters: []swaggers.Parameter{},
		Responses:  map[string]swaggers.Response{},
		Tags:       []string{},
	}
	ret.Parameters = h.CreateSwaggersParameter()
	return &ret

}
func (h *apiMethodInfo) CreateSwaggersParameter() []swaggers.Parameter {
	ret := []swaggers.Parameter{}
	for _, x := range h.indexOfFieldInUrl {
		contextType := h.method.Type.In(h.indexOfContextInfo)
		if contextType.Kind() == reflect.Ptr {
			contextType = contextType.Elem()
		}
		fmt.Println(contextType.String())
		field := contextType.FieldByIndex(x)

		ret = append(ret, swaggers.Parameter{
			Type: "string",
			Name: field.Name,
			In:   "path",
		})

	}
	for i, x := range h.indexOfArgHasFileUpload {
		containerOfArg := h.method.Type.In(x)
		if containerOfArg.Kind() == reflect.Ptr {
			containerOfArg = containerOfArg.Elem()
		}
		if len(h.fieldIndexOfFileUpload) == 0 {
			ret = append(ret, swaggers.Parameter{
				Type: "file",
				Name: fmt.Sprintf("File%d", x),
				In:   "formData",
			})
		} else {
			if len(h.fieldIndexOfFileUpload[i]) == 0 {
				ret = append(ret, swaggers.Parameter{
					Type: "file",
					Name: fmt.Sprintf("File%d", x),
					In:   "formData",
				})
			}
			for _, y := range h.fieldIndexOfFileUpload[i] {
				if len(y) == 0 {
					ret = append(ret, swaggers.Parameter{
						Type: "file",
						Name: fmt.Sprintf("File%d", x),
						In:   "formBody",
					})
				} else {
					field := containerOfArg.FieldByIndex(y)
					ret = append(ret, swaggers.Parameter{
						Type: "file",
						Name: field.Name,
						In:   "formBody",
					})
				}
			}
			for _, y := range h.fieldIndex[i] {
				field := containerOfArg.FieldByIndex(y)

				ret = append(ret, swaggers.Parameter{
					Type: "string",
					Name: field.Name,
					In:   "formData",
				})
			}
		}

	}
	return ret

}
