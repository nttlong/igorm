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
	IsSlug     bool
	FieldIndex []int
}
type QueryParam struct {
	Name       string
	FieldIndex []int
}
type HandlerInfo struct {
	//BaseUrl        string
	IndexOfArg     int
	TypeOfArgs     reflect.Type
	TypeOfArgsElem reflect.Type
	FieldIndex     []int
	ReceiverIndex  int

	ReceiverType     reflect.Type
	ReceiverTypeElem reflect.Type

	Method    reflect.Method
	RouteTags []string
	/*
		Master uri
	*/
	Uri      string
	UriQuery string
	/*
		if Uri start with '/' IsAbsUrl
	*/
	IsAbsUri     bool
	IsQueryUri   bool
	QueryParams  []QueryParam
	RegexUri     string
	RegexUriFind regexp.Regexp
	/*
		This value will be used in http Handler
	*/
	UriHandler string
	/*
		If current http handler will use regex for router
	*/
	IsRegexHandler bool
	/*
		List of Uriparam
		Exmaple:
			<controller name>/<handler name>/{param1}/.../{param2}/{*slug}
	*/
	UriParams []uriParam
	/*
		 index of arg in handler is injector
		 Example:
			func(c*conroller) (h*wx.Handler, S3 wx.Depen[S2Utils],googelDrive ws.Depen[Gg], user wx.UserClaims)
															^				 ^
															[2]				[3]
			value will be [2,3]-----------------------------------------------------
	*/
	IndexOfInjectors []int
	/*
		if handler has inject

		Example:
			func(c*conroller) (h*wx.Handler, inject1 wx.Depen[S2Utils], user wx.UserClaims)
	*/
	HasInjector bool
	/*

		This field is index of all fields are multipart.HeadeFile (at first level of struct)
		Example:
			func(c*MyController) (ctx *wx.Handler, data struct {
				...
				File1 multipart.HeadeFile <--field inxe on Strutc is 10
				...
				Files[] multipart.HeadeFile <--field inxe on Strutc is 13
				....
			})
			FormUploadFile=[10,13]

	*/
	
	FormUploadFile []int
	/*
			If the handler has a request body argument and is neither wx.Handler nor inject or auth,
		    the value of this field is the index of that argument.
			Default value is -1 and request body is always struct
	*/
	IndexOfRequestBody int
	/*
		Type of arg is request body in Ptr
	*/
	TypeOfRequestBody reflect.Type
	/*
		Type of arg is request body in Struct
		Note: reques body is always struct
	*/
	TypeOfRequestBodyElem reflect.Type
	/*
		index of arg is UserClaims
		Example:
			func(c*conroller) (h*wx.Handler, inject1 wx.Depen[S2Utils], user wx.UserClaims)
			IndexOfAuthClaimsArg=3
			if handler of controller has any field is wx.UserClaims, looks like:

			type MyController struct {]
			type UserClaimsHandler struct {
				User wx.UserClaims
				...
			}
			func(c*MyController) (h*UserClaimsHandler, inject1 wx.Depen[S2Utils])
			IndexOfAuthClaimsArg=1

	*/
	IndexOfAuthClaimsArg int
	/*
		if handler of controller has any field is wx.UserClaims, looks like:

		type MyController struct {]
		type UserClaimsHandler struct {
			...
			User wx.UserClaims <-FieldIndex is 15
			...
		}
		func(c*MyController) (h*UserClaimsHandler, inject1 wx.Depen[S2Utils])
		IndexOfAuthClaims=[15]
	*/
	IndexOfAuthClaims []int
	/*
		Http method of handler
	*/
	HttpMethod string
}
