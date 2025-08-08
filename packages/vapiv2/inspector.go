package vapi

import (
	"fmt"
	"reflect"
	"strings"
)

type inspectors struct {
}

/*
	/*

This function inspects a reflected Method to determine if it is an HTTP handler.
It checks if any of the method's arguments is of an HTTP method type (e.g., HttpGet, HttpPost, etc.).
If such an argument is found, it confirms that the method is an HTTP handler and returns
an inspection info object related to that method.

Examples:

	func (a *A) Func1(ctx *HttpPost, ...)  // returns inspection info of Func1 (recognized as an HTTP handler)
	func (a *A) Func2(ctx *HttpGet, ...)   // returns inspection info of Func2 (recognized as an HTTP handler)
	func (a *A) Func3(a int, b int)        // returns nil (not an HTTP handler)
*/
func (inspector *inspectors) Create(method reflect.Method) (*inspectInfo, error) {
	var ret *inspectInfo
	IndexOfInjectors := []int{}
	IndexOfData := -1 //<-- Assume that the method does not have any post data.
	//Essentially, HTTP GET does not have any post data in the request body
	//step 1- inspect and make inspectInfo if available
	for i := 1; i < method.Type.NumIn(); i++ {
		//	^
		//	|
		//	[Skip 0, 0 index is belong to receier]

		typ := method.Type.In(i)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem() //<--/* ensure typ of in arg i is not ptr*/
		}
		if httpMethod := httpUtilsTypeInstance.InspectHttpMethodFromType(typ); httpMethod != "" {
			/*
			 If typ is HttpMethod or any of its fields are HttpMethod
			*/
			if ret == nil {
				ret = &inspectInfo{
					Args: inspectArs{
						IndexOfHttpMethod: i,
					},
					Route: routeInfo{
						Method: httpMethod,
					},
				}
				continue //<-- ok, has HttpMethod, go next arg
			}

		}
		/*
			detect if arg is injector
		*/
		if httpUtilsTypeInstance.IsInjector(typ) {
			IndexOfInjectors = append(IndexOfInjectors, i)
			continue //<-- ok, this arg is injector, go next arg

		}
		IndexOfData = i

	}
	if ret == nil { //<-- not avaliable return nil
		return nil, nil //<--no error
	}
	ret.Args.IndexOfData = IndexOfData //<-- set index of arg where is Data for hanlder
	// Note: some http handler do not have data post. In this case, ret.Args.IndexOfData will be set  -1
	ret.Args.IndexOfInjectors = IndexOfInjectors
	// end of step 1
	// step 2- build route info

	routeTag := httpUtilsTypeInstance.GetRouteTag(method.Type.In(ret.Args.IndexOfHttpMethod))
	ret.Route.Tags = routeTag
	// step 2.1 -- set up default uri
	/*
		The default URI is constructed by combining the package path of the receiver's type and the method name.

		For example:
		  Given a method:
		    func (a *A) Handler1(...)
		  where the type A is defined in the package "MyPackage",

		The resulting default URI will be:
		  "my-package/handler1"

		Note:
		- The package name is converted to lowercase, and words are separated by hyphens.
		- The method name is converted to lowercase.

		This default URI serves as a standard route path when no custom URI is specified.
	*/
	ret.ReceiverType = method.Type.In(0)
	ret.ReceiverTypeEle = ret.ReceiverType
	if ret.ReceiverTypeEle.Kind() == reflect.Ptr {
		ret.ReceiverTypeEle = ret.ReceiverTypeEle.Elem()
	}

	ret.HttpMethodType = method.Type.In(ret.Args.IndexOfHttpMethod)
	ret.HttpMethodTypeElem = ret.HttpMethodType
	ret.Method = method //<-- store the method for web invoke in the future
	if ret.HttpMethodTypeElem.Kind() == reflect.Ptr {
		ret.HttpMethodTypeElem = ret.HttpMethodTypeElem.Elem()
	}
	placeholders := strings.Split(ret.ReceiverTypeEle.String(), ".")
	placeholders = append(placeholders, method.Name)
	// step 2.2
	for i := 0; i < len(placeholders); i++ {
		placeholders[i] = ToKebabCase(placeholders[i]) //<-- convert to ToKebab Case version
	}
	prefixUri := strings.Join(placeholders, "/")
	ret.Route.Uri = prefixUri
	ret.Route.UriHandler = ret.Route.Uri
	ret.Route.RegexUri = EscapeSpecialCharsForRegex(ret.Route.Uri) //<-- by default RegexUri has same valur of Uri
	ret.Route.UriParams = []string{}
	ret.Route.IndexOfFieldInUri = [][]int{}
	// step 2.3
	if routeTag != "" {
		items := strings.Split(routeTag, ";")
		for _, item := range items {
			if strings.HasPrefix(item, "uri:") {
				ret.Route.Uri = strings.TrimPrefix(item, "uri:")
				ret.Route.Uri = strings.ReplaceAll(ret.Route.Uri, "@", ToKebabCase(method.Name)) //<-- convert to ToKebab Case version
				/*
					Uri has been change
					So "RegexUri" and "UriHandler" need to be changed
				*/
				ret.Route.RegexUri = EscapeSpecialCharsForRegex(ret.Route.Uri)
				ret.Route.UriHandler = ret.Route.Uri

			}
		}

		if strings.Contains(ret.Route.Uri, "{") { //<-- 	if uri look like {placeHolder1}/../{placeHolder n}
			ret.Route.UseRegex = true
			ret.Route.UriHandler = strings.Split(ret.Route.Uri, "{")[0] + "/"
			uriItems := strings.Split(ret.Route.Uri, "{")
			for _, item := range uriItems {
				if strings.Contains(item, "}") {
					placeHolder := strings.Split(item, "}")[0]
					ret.Route.UriParams = append(ret.Route.UriParams, placeHolder)
					ret.Route.RegexUri = strings.ReplaceAll(ret.Route.RegexUri, "{"+placeHolder+"}", "(.*)")
					/*
					 as a commnent in Route.IndexOfFieldInUri
					 The system just scan the first level of ReceiverTypeEle
					*/
					field, found := ret.HttpMethodTypeElem.FieldByNameFunc(func(s string) bool {
						return strings.EqualFold(s, placeHolder)
					})
					if !found { // invalid hanlder set up
						return nil, fmt.Errorf("Invalid handler setup of method %s in type %s", method.Name, ret.ReceiverTypeEle.String())

					}
					ret.Route.IndexOfFieldInUri = append(ret.Route.IndexOfFieldInUri, field.Index)

				}
			}
		}
		/*
		  check if Uri has "@" and start by '/'
		  Note start by '/' that mean uri was handled form root
		*/
		if strings.Contains(ret.Route.Uri, "@") && ret.Route.Uri[0] != '/' {
			ret.Route.Uri = strings.ReplaceAll(ret.Route.Uri, "@", ToKebabCase(prefixUri))
		} else {
			ret.Route.Uri = prefixUri + "/" + ret.Route.Uri
			ret.Route.RegexUri = prefixUri + "/" + ret.Route.RegexUri
			ret.Route.UriHandler = prefixUri + "/" + ret.Route.UriHandler
		}
	}
	// end of step 2
	// buil UirParams

	return ret, nil
}

var inspector = &inspectors{}
