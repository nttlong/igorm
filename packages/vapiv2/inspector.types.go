package vapi

import "reflect"

type inspectArs struct {
	/*
		index of args in method is HttpGet, HttpPost,....
		Example:
			type A struct {
			}
			func (a*A) Update(ctx * HttpPost,dataPost ...) {
										^
										|
									[ Here	]
			}
	*/
	IndexOfHttpMethod int
	/*
		index of args in method is HttpGet, HttpPost,....
		Example:
			type A struct {
			}
			func (a*A) Update(... HttpPost,dataPost ...) {
												^
												|
											[ Here	]
			}

		##Note: only one arg in args of method play as Post Data from client
	*/
	IndexOfData int
	/*
		lits of index of args are Injector
		Example:
			type A strutc {}
			func (a*A) Update(...Httpmethod,data..., auth Inject[Auth],...,Inject[Db])
															^					^
															[here].... and ... [here]

	*/
	IndexOfInjectors []int
}

// ---------------------------------------------------
type routeInfo struct {

	/*
		The rule form at Uir value are:

		1- If no route tags are found, the URI will be the kebab-case version of the method name of the pointer struct
			Exanple:
				func(a*A) MyMethod(....) -> uri: my-method
		2- If route are found , the URI will be route and replace any @ with be the kebab-case version of the method name of the pointer struct
			Example 1;
				func(a*A) MyMethod(ctx struct {
					HttpGet `route:uri"@/file/{fileId}.{fileExt}"`
					FileID string
					FileExt string

				}....) --> my-method/file/{fielId}.{fileExt}
			Examle 2:
				func(a*A) MyMethod(ctx struct {
					HttpGet `route:ur:file/{fileId}.{fileExt}"`
					FileID string
					FileExt string

				}....) --> file/{fielId}.{fileExt}
			Note: The '@' symbol in the route URI represents the kebab-case version of the Go method name
	*/
	Uri string
	/*
		When handling HTTP requests, some URIs require regex matching to detect dynamic values.

		For example:
		  Client sends a GET request to:
		    http://localhost/user/avatar/12345.jpg
		                               ^
		                               |
		                           [User ID]

		In this case:
		  - UriHandler is the static prefix part of the URI:
		      "/user/avatar/"
		    which is used to quickly filter matching requests.

		  - RegexUri is a regex pattern that matches the full path including the dynamic part,
		    for example:
		      \/user\/avatar\/(.*)\.jpg
		    This regex captures the dynamic User ID (e.g., "12345") from the URI.

		Together, UriHandler helps quickly narrow down the relevant paths,
		and RegexUri performs detailed matching and extraction of path variables.
	*/
	UseRegex bool
	ISAbsUri bool

	/*
		By default, RegexUri is derived from the original Uri by escaping all special regex characters.

		This is done by calling the function:
		  EscapeSpecialCharsForRegex(Uri)

		The purpose is to safely convert the Uri string into a regex pattern that matches the literal URI path,
		preventing special characters from being interpreted as regex operators.
	*/
	RegexUri string

	/*
		By default, UriHandler is equal to the full URI.

		In the case of using a regex or placeholders in the URI,
		UriHandler is derived by taking the substring of the URI
		up to (but not including) the first '{' character.

		For example:
		  URI: /a/{name1}/{name2}
		  UriHandler: /a/    // Note: UriHandler must end with '/'

		This allows the UriHandler to represent the static prefix portion of the URI before any path parameters.
	*/
	UriHandler string
	Method     string
	/*
		As usual, value of this field is Uri. In special case such as: Uri has {...}/{...}
		RegexUri will be replace with .*
	*/
	RequestContentType string
	/*
		IndexOfFieldInUri represents the positions of struct fields as they appear in the URI template.
		For example:
		URI template: /test/{Field1}/{Field2}
		Struct:
		  HttpGet  `route:uri"/test/{Field1}/{Field2}"`
		  Field1 string
		  Field2 string
		IndexOfFieldInUri would be a slice of indexes where each index corresponds to
		the position of the field in the struct relative to the URI placeholders.
		In this case:
		  Field1 corresponds to index [1]
		  Field2 corresponds to index [2]
	*/
	IndexOfFieldInUri [][]int
	/*
		UriParams holds the names of the struct fields that correspond to the URI path parameters.
		For example:
		  Given the struct:
		    HttpGet  `route:uri"/test/{Field1}/{Field2}"`
		    Field1 string
		    Field2 string
		  The UriParams slice will be:
		    []string{"Field1", "Field2"}
		These represent the path parameters extracted from the URI template in order.
	*/
	UriParams []string
	Tags      string //<-- orginal route tag of context field
}

// ------------------------------------------------
type inspectInfo struct {
	Args               inspectArs
	Route              routeInfo
	ReceiverTypeEle    reflect.Type
	ReceiverType       reflect.Type
	HttpMethodType     reflect.Type
	HttpMethodTypeElem reflect.Type
	Method             reflect.Method
}
