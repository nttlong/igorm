package vapi

import (
	"encoding/json"
	"fmt"
	"mime"
	"mime/multipart"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type urlParamInfo struct {
	Name       string
	Type       string
	FieldIndex []int
}
type formParamInfo struct {
	Name       string
	Type       string
	FieldIndex []int
}
type objectParamInfo struct {
	Name       string
	Type       string
	ParamIndex int
	ParamType  reflect.Type
}
type methodInfo struct {
	/*
		A part of the function name in struct
		Example:
			func(st * testStructt) Handle_Hello_Get
		The value of field is Hello
	*/
	MethodName string //<-- any method of struct like Handler<Something>Post the method name is Somthing

	/*
		A part of the function name in struct
		Example:
			func(st * testStructt) Handle_Hello_Get
		The value of field is Get
	*/
	HttpMethod string
	/*
		A part of the function name in struct
		Example:
			func(st * testStructt) TestAPI_Handle_Hello_Get
		The value of field is TestAPI
	*/
	Tag    string
	Method reflect.Method
	/*
		is a first type of method.In
	*/
	ReceiverTypeOfInstance reflect.Type
	/*
		is a second type of method.In
	*/
	ReceiverType reflect.Type

	/*
		ReceiverType.String() replace "." with "/" and
		Is a combination of ReceiverType.String() replace "." with "/" and MethodName, all to lower
	*/
	PathOfMethod          string
	Url                   string
	MasterUrl             string
	HandlerUrl            string
	IndexOfHttpContextArg int
	IndexfUserClaimsArg   int
	IndexDataArg          int
	Description           string
	SwaggerInputType      string
	regexpUrl             *regexp.Regexp
	Priority              int
	Parameters            []interface{}
	ResponseMimeType      string
	RequestMimeType       string
	Ext                   string
}
type utils struct {
	mapGoTypeToSwaggerType map[reflect.Type]string
	rootPackage            string
	handler                []*methodInfo
	mapUrlMethodInfo       map[string]*methodInfo
}

func (u *utils) GetType(typ reflect.Type) string {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if ret, ok := u.mapGoTypeToSwaggerType[typ]; ok {
		return ret
	} else {
		return "object"
	}

}

// ParseMethodName parses the method name to extract MethodName, HttpMethod, Tag, and constructs PathOfMethod.
// It uses regex to match patterns like Handle_<MethodName>_<HttpMethod> or <Tag>_Handle_<MethodName>_<HttpMethod>.
func (u *utils) ParseMethodName(method reflect.Method) (*methodInfo, error) {
	info := methodInfo{
		IndexOfHttpContextArg:  -1,
		IndexfUserClaimsArg:    -1,
		Method:                 method,
		ReceiverType:           method.Type.In(0), // The first parameter is the receiver type
		Parameters:             []interface{}{},
		ReceiverTypeOfInstance: method.Type.In(0).Elem(),
	}

	info.Tag = info.ReceiverType.Elem().String()
	for i := 1; i < method.Type.NumIn(); i++ {
		typ := method.Type.In(i)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if typ == reflect.TypeOf(HttpContext{}) {
			info.IndexOfHttpContextArg = i
		} else if typ == reflect.TypeOf(UserClaims{}) {
			info.IndexfUserClaimsArg = i
		} else {
			info.IndexDataArg = i
		}

	}
	if info.IndexOfHttpContextArg == -1 {
		return nil, nil
	}
	// if info.IndexDataArg > -1 {

	// }

	path := info.ReceiverType.Elem().String()
	path = strings.ReplaceAll(path, ".", "/")
	path = strings.ToLower(path)
	urlMethod := strings.ToLower(method.Name)

	if info.IndexDataArg > -1 {
		regexpUrl := path + "/" + strings.ToLower(method.Name)
		info.MasterUrl = "/" + regexpUrl + "/"
		urlsInfos, defaultMethod := u.GetUrlParamInfo(method.Type.In(info.IndexDataArg), []int{})
		info.Url = "/" + path + "/" + urlMethod
		//info.Url = "/" + path + "/" + urlMethod
		info.HandlerUrl = "/" + path
		for _, urlsInfo := range urlsInfos {
			if urlsInfo.url == "" && urlsInfo.form == "" {
				info.HttpMethod = defaultMethod
				info.RequestMimeType = "application/json"
				info.Parameters = append(info.Parameters, formParamInfo{
					Name:       urlsInfo.name,
					Type:       urlsInfo.fieldType,
					FieldIndex: urlsInfo.fieldIndex,
				})

				continue
			}
			if urlsInfo.form != "" {
				info.HttpMethod = defaultMethod
				info.RequestMimeType = "multipart/form-data"
				info.Parameters = append(info.Parameters, formParamInfo{
					Name:       urlsInfo.form,
					Type:       urlsInfo.fieldType,
					FieldIndex: urlsInfo.fieldIndex,
				})

				continue
			}
			if urlsInfo.url != "" {
				info.HttpMethod = "GET"

				argsName := strings.Split(urlsInfo.url, "{")[1]
				argsName = strings.Split(argsName, "}")[0]

				if strings.Contains(urlsInfo.url, ".") {
					regexpUrl += "/.*." + strings.Split(urlsInfo.url, ".")[1]
					ext := strings.Split(urlsInfo.url, ".")[1]
					info.HandlerUrl += "/:" + argsName + "." + ext
					info.Url += "/{" + argsName + "}" + "." + ext
					info.Priority = 0
					info.Ext = ext
					info.ResponseMimeType = mime.TypeByExtension("." + ext)
				} else {
					regexpUrl += "/.*"
					info.HandlerUrl += "/:" + argsName
					info.Url += "/{" + argsName + "}"
					info.Priority = 1
				}

				info.Parameters = append(info.Parameters, urlParamInfo{
					Name:       argsName,
					Type:       "string",
					FieldIndex: urlsInfo.fieldIndex,
				})
				continue
			}

		}
		info.HttpMethod = defaultMethod
		fmt.Println(`^/` + regexpUrl + `$`)
		info.regexpUrl = regexp.MustCompile(`^/` + regexpUrl + `$`)

	}
	if len(info.Parameters) == 0 {
		info.HttpMethod = "POST"
		inputDataType := method.Type.In(info.IndexDataArg)
		if inputDataType.Kind() == reflect.Ptr {
			inputDataType = inputDataType.Elem()

		}

		sample := reflect.New(inputDataType).Elem().Interface()
		jsonData, err := json.MarshalIndent(sample, " ", " ")
		if err != nil {
			return nil, err
		}
		info.Description = string(jsonData)
		if inputDataType.Kind() == reflect.Struct {
			info.SwaggerInputType = "object"

		} else {
			info.SwaggerInputType = u.GetType(inputDataType)
		}
		info.ResponseMimeType = mime.TypeByExtension(".json")

		info.Url = "/" + path + "/" + strings.ToLower(method.Name)
		info.MasterUrl = info.Url
		info.regexpUrl = regexp.MustCompile(`^/` + path + "/" + strings.ToLower(method.Name) + `/?$`)
		info.HandlerUrl = "/" + path
		info.Priority = 4
		info.RequestMimeType = "application/json"
		info.Parameters = append(info.Parameters, objectParamInfo{
			Name:       "data",
			Type:       "object",
			ParamIndex: info.IndexDataArg,
			ParamType:  method.Type.In(info.IndexDataArg),
		})

	}

	return &info, nil
}

type apiTagInfo struct {
	url         string
	description string
	fieldIndex  []int
}

func (u *utils) GetApiTag(field reflect.StructField, preFixFieldIndex []int) *apiTagInfo {
	tag := field.Tag.Get("api")
	if tag == "" {
		return nil
	}
	items := strings.Split(tag, ",")
	ret := apiTagInfo{}
	for _, item := range items {
		lowerItem := strings.ToLower(item)
		if strings.HasPrefix(lowerItem, "url:") {
			ret.url = item[4:]

		}
		if strings.HasPrefix(lowerItem, "description:") {
			ret.description = item[12:]
		}

	}
	ret.fieldIndex = append(preFixFieldIndex, field.Index...)
	return &ret
}

type itemGetUrlParamInfo struct {
	url        string
	form       string
	name       string
	fieldIndex []int
	fieldType  string
}

func (u *utils) GetUrlParamInfo(typ reflect.Type, fieldIndex []int) ([]itemGetUrlParamInfo, string) {
	ret := []itemGetUrlParamInfo{}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	itemType := ""
	defaultMethod := "GET"
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fielType := field.Type
		if fielType.Kind() == reflect.Ptr {
			fielType = fielType.Elem()
		}

		if field.Anonymous {
			subRet, httpMethod := u.GetUrlParamInfo(field.Type, field.Index)
			defaultMethod = httpMethod
			ret = append(ret, subRet...)
		} else if fielType == reflect.TypeOf(multipart.FileHeader{}) {
			if itemType == "" {
				itemType = "form"
			}

			ret = append(ret, itemGetUrlParamInfo{
				form:       field.Name,
				fieldIndex: fieldIndex,
				fieldType:  "file",
			})
		} else {
			tagInfo := u.GetApiTag(field, fieldIndex)
			if tagInfo != nil {
				if tagInfo.url != "" {

					ret = append(ret, itemGetUrlParamInfo{
						url:        tagInfo.url,
						fieldIndex: tagInfo.fieldIndex,
					})

				}

			} else {
				fieldType := field.Type
				if fieldType.Kind() == reflect.Ptr {
					fieldType = fieldType.Elem()
				}
				if fieldType == reflect.TypeOf(uuid.UUID{}) || fieldType == reflect.TypeOf(time.Time{}) {

					ret = append(ret, itemGetUrlParamInfo{
						name:       field.Name,
						fieldIndex: append(fieldIndex, field.Index...),
						fieldType:  "string",
					})
				} else if fieldType.Kind() == reflect.Struct {
					defaultMethod = "POST"
					ret = append(ret, itemGetUrlParamInfo{
						name:       field.Name,
						fieldIndex: append(fieldIndex, field.Index...),
						fieldType:  "object",
					})
				} else {
					ret = append(ret, itemGetUrlParamInfo{
						name:       field.Name,
						fieldIndex: append(fieldIndex, field.Index...),
						fieldType:  "string",
					})
				}
			}

		}
	}
	if itemType == "" && defaultMethod == "POST" {
		return ret, "POST"
	} else {
		for i := range ret {

			if itemType == "url" {
				if ret[i].url == "" {
					ret[i].url = ret[i].name
				}

				continue
			}
			if itemType == "form" {
				if ret[i].form == "" {
					ret[i].form = ret[i].name

				}
				if ret[i].fieldType == "" {
					ret[i].fieldType = "form"

				}
				continue
			}
		}
	}
	return ret, defaultMethod

}

func (u *utils) ParseMethods(typ reflect.Type) ([]*methodInfo, error) {
	ret := []*methodInfo{}
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		info, err := u.ParseMethodName(method)
		if err != nil {
			return nil, err
		}
		if info != nil {
			ret = append(ret, info)
		}
	}
	//sort ret by priority
	for i := 0; i < len(ret); i++ {
		for j := i + 1; j < len(ret); j++ {
			if ret[i].Priority > ret[j].Priority {
				ret[i], ret[j] = ret[j], ret[i]
			}
		}
	}
	return ret, nil

}
func (u *utils) GetOperation(info *methodInfo) *Operation {

	ret := Operation{
		Tags:       []string{info.Tag},
		Summary:    "",
		Consumes:   []string{info.RequestMimeType},
		Parameters: []Parameter{},
		Responses: map[string]Response{
			"200": {
				Description: "",
				Schema:      nil,
			},
		},

		Produces: []string{info.RequestMimeType},
	}
	for _, param := range info.Parameters {
		switch v := param.(type) {
		case urlParamInfo:
			urlParamItem := param.(urlParamInfo)
			ret.Parameters = append(ret.Parameters, Parameter{
				Type:        urlParamItem.Type,
				Description: v.Name,
				Name:        v.Name,
				In:          "path",
				Required:    true,
			})
		case formParamInfo:
			formParamItem := param.(formParamInfo)
			ret.Parameters = append(ret.Parameters, Parameter{
				Type:        formParamItem.Type,
				Description: v.Name,
				Name:        v.Name,
				In:          "formData",
				Required:    false,
			})
		case objectParamInfo:
			objectParamItem := param.(objectParamInfo)
			typ := objectParamItem.ParamType
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			sample := reflect.New(typ).Elem().Interface()
			jsonData, err := json.MarshalIndent(sample, " ", " ")
			if err != nil {
				panic(err)
			}

			ret.Parameters = append(ret.Parameters, Parameter{
				Type:        "object",
				Description: string(jsonData),
				Name:        v.Name,
				In:          "body",
				Required:    true,
			})

		}
	}
	if info.IndexfUserClaimsArg != -1 {

		ret.Security = []map[string][]string{
			{
				"OAuth2Password": {},
			},
		}
	}

	return &ret
}
func (u *utils) SortHandlers() {
	for i := 0; i < len(u.handler); i++ {
		for j := i + 1; j < len(u.handler); j++ {
			if u.handler[i].Priority < u.handler[j].Priority {
				u.handler[i], u.handler[j] = u.handler[j], u.handler[i]
			}
		}
	}

}

var utilsInstance *utils = &utils{
	mapGoTypeToSwaggerType: map[reflect.Type]string{
		reflect.TypeOf(int(0)):      "integer",
		reflect.TypeOf(int8(0)):     "integer",
		reflect.TypeOf(int16(0)):    "integer",
		reflect.TypeOf(int32(0)):    "integer",
		reflect.TypeOf(int64(0)):    "integer",
		reflect.TypeOf(uint(0)):     "integer",
		reflect.TypeOf(uint8(0)):    "integer",
		reflect.TypeOf(uint16(0)):   "integer",
		reflect.TypeOf(uint32(0)):   "integer",
		reflect.TypeOf(uint64(0)):   "integer",
		reflect.TypeOf(float32(0)):  "number",
		reflect.TypeOf(float64(0)):  "number",
		reflect.TypeOf(uuid.UUID{}): "string",

		reflect.TypeOf(""):   "string",
		reflect.TypeOf(true): "boolean",
	},
	handler:     []*methodInfo{},
	rootPackage: reflect.TypeOf(utils{}).PkgPath(),
}
