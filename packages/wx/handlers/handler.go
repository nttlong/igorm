package handlers

import (
	"net/http"
	"reflect"
	"regexp"
)

type Handler struct {
	Res     http.ResponseWriter
	Req     *http.Request
	BaseUrl string
}
type uriParam struct {
	Position   int
	Name       string
	FieldIndex []int
}
type HandlerInfo struct {
	BaseUrl        string
	IndexOfArg     int
	TypeOfArgs     reflect.Type
	TypeOfArgsElem reflect.Type
	FieldIndex     []int
	ReceiverIndex  int

	ReceiverType          reflect.Type
	ReceiverTypeElem      reflect.Type
	Method                reflect.Method
	RouteTags             []string
	Uri                   string
	RegexUri              string
	RegexUriFind          regexp.Regexp
	UriHandler            string
	IsRegexHandler        bool
	UriParams             []uriParam
	IndexOfInjectors      []int
	HasInjector           bool
	FormUploadFile        []int
	IndexOfRequestBody    int
	TypeOfRequestBody     reflect.Type
	TypeOfRequestBodyElem reflect.Type
	IndexOfAuthClaimsArg  int
	IndexOfAuthClaims     []int
	HttpMethod            string
}
