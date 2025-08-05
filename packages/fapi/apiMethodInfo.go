package fapi

import "reflect"

type apiMethodInfo struct {
	method       reflect.Method
	param        reflect.Type
	recieverType reflect.Type
	/*
		Inspect that method has arg with typ is Context
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
	indexOfFieldInUrl       [][]int
	regExpUri               string
	routeHandler            string
	httpMethod              string
	hasUploadFile           bool
	requestContentType      string
	indexOfArgHasFileUpload []int
	fieldIndexOfFileUpload  [][][]int
}

func (am *apiMethodInfo) IsAbsUri() bool {
	return am.tags.Url[0] == '/'
}
