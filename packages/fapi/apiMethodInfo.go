package fapi

import "reflect"

type apiMethodInfo struct {
	method       reflect.Method
	param        reflect.Type
	recieverType reflect.Type
	/*
					If a method has an argument of type Context.
		If a method of a struct has any argument with the type Context or an embedded struct that contains Context, then the function RegisterController will use that method as an HTTP handler in the web server.
		The function indexOfContextInfo determines the position of the Context argument or the argument type that embeds Context.


					Nếu phương thức của struct có tham số kiểu Context hoặc thậm chí kiểu struct mà embed (nhúng) Context bên trong, thì hàm RegisterController sẽ sử dụng phương thức đó làm HttpHandler trong web server.

					Trường indexOfContextInfo dùng để xác định vị trí (index) của tham số kiểu Context hoặc tham số kiểu struct có embed Context trong danh sách tham số của phương thức đó.

	*/
	indexOfContextInfo int
	typeOfContextInfo  reflect.Type
	/*
	 list of index of args are injectors
	*/
	indexOfInjectors []int
	routeTags        string
	tags             apiTagsInfo
	/*
		if tags.Url look like abc/{Param1}/{Param2}/../{ParamN}
		and method look like
		MyApi(podData PostDataContext)
		and PostDataContext is
		type PostDataContext struct {
			*Context
			Param1 string
			...
			ParamN string

		}
		This field is
		[[1],[2],...[n]]

	*/
	indexOfFieldInUrl  [][]int
	regExpUri          string
	routeHandler       string
	httpMethod         string
	hasUploadFile      bool
	requestContentType string
	/*
		list of index of func args where contains file upload
		Example:
			tup FX struct { F multipart.FileHeader }
			typp A struct { FX} <-- this struct have embeded field with FileUploade
			func B ( a A, b int, c multipart.FileHeader)
					   ^			^
					   |			|
					   |			[2] <-- this args is FileUpload
					   [0]<-- A have a field is FileUpload
			indexOfArgHasFileUpload -->{0,2}
	*/
	indexOfArgHasFileUpload []int
	/*
		list of index of func args where contains file upload
		Example:
			tup FX struct {
				F multipart.FileHeader <--[0]
			}
			typp A struct {
				Code string
				FX <--[1]

			} <-- this struct have embeded field with FileUploade
			func B ( a A, b int, c multipart.FileHeader)
					   ^			^
					   |			|
					   |			[2] <-- this args is FileUpload
					   [0]<-- A have a field is FileUpload
			fieldIndexOfFileUpload -->[({1,0}),(<Direct in args>)]
	*/
	fieldIndexOfFileUpload [][][]int
	ResponseContentType    string
}

func (am *apiMethodInfo) IsAbsUri() bool {
	return am.tags.Url[0] == '/'
}
