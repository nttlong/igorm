package fapi

import (
	"reflect"
	"strings"
)

/*
The function inspects the Context struct type and looks for struct tags of the form route:"uri:...".

It parses the route:"uri:myserver/{ServiceName}" tag and converts the URI pattern into a regular expression string stored in regExpUri.

Specifically, it replaces the {ServiceName} path parameter placeholder with a regex capture group, generating 'myserver\/(.*)'.

At the same time, it identifies which struct fields correspond to those parameters in the URI (e.g., the ServiceName field).
It then returns the indices of those fields in the struct type as a way to map the regex capture groups back to struct fields.
According by the tags of arg which base on Context, this function will update
"regExpUri" of apiMethodInfo and return all field Index of typ which have been a part of route uri
Example:

	type MyContext struct {
		*Context `route:"uri:myserver/{ServiceName}"`
											^
		ServiceName string //---------------|
	}
	regExpUri of apiMethodInfo will be 'myserver\/(.*)'

Dựa theo các tag của arg dựa trên Context, hàm này sẽ cập nhật trường regExpUri của apiMethodInfo và trả về tất cả các chỉ số (index) của các trường (fields) trong kiểu typ mà đã trở thành một phần của URI route.
*/
func (apiMethod *apiMethodInfo) ExtractFieldIndex(typ reflect.Type, index []int) []int {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Anonymous {
			index = apiMethod.ExtractFieldIndex(field.Type, index)
			if index != nil {
				return index
			}
		}
		if strings.Contains(apiMethod.tags.Url, "{"+field.Name+"}") {
			index = append(index, field.Index...)
			apiMethod.regExpUri = strings.ReplaceAll(apiMethod.regExpUri, "{"+field.Name+"}", "(.*)")
			return index
		}
	}
	return nil

}

func (apiMethod *apiMethodInfo) GetRoute() string {
	apiMethod.indexOfFieldInUrl = [][]int{}
	apiMethod.routeHandler = apiMethod.tags.Url
	if strings.Contains(apiMethod.tags.Url, "{") {
		apiMethod.routeHandler = strings.Split(apiMethod.tags.Url, "{")[0]
		apiMethod.regExpUri = api.EscapeRegExp(apiMethod.tags.Url)

		if apiMethod.typeOfContextInfo == nil {
			apiMethod.typeOfContextInfo = apiMethod.method.Type.In(apiMethod.indexOfContextInfo)
		}
		contextType := apiMethod.typeOfContextInfo
		if contextType.Kind() == reflect.Ptr {
			contextType = contextType.Elem()
		}
		for i := 0; i < contextType.NumField(); i++ {
			field := contextType.Field(i)
			if field.Anonymous {
				fieldIndex := apiMethod.ExtractFieldIndex(field.Type, field.Index)
				if fieldIndex != nil {
					apiMethod.indexOfFieldInUrl = append(apiMethod.indexOfFieldInUrl, field.Index)
				}
				continue
			}
			if strings.Contains(apiMethod.tags.Url, "{"+field.Name+"}") {
				apiMethod.indexOfFieldInUrl = append(apiMethod.indexOfFieldInUrl, field.Index)
				apiMethod.regExpUri = strings.ReplaceAll(apiMethod.regExpUri, "{"+field.Name+"}", "(.*)")
			}

		}

	}
	if apiMethod.tags.Method == "" {
		apiMethod.tags.Method = "post"
		apiMethod.requestContentType = "application/json"
	}
	apiMethod.httpMethod = strings.ToUpper(apiMethod.tags.Method)
	apiMethod.hasUploadFile = apiMethod.HasUploadFile()
	if apiMethod.hasUploadFile {
		apiMethod.requestContentType = "multipart/form-data"
	}
	return apiMethod.regExpUri
}

/*
This function inspects all the arguments of a method and all the members of each argument's type to check if they contain any multipart.FileHeader or something similar. If it finds any, the function returns true.

	Additionally, the function updates two fields of apiMethodInfo:

	indexOfArgHasFileUpload → a list of argument indices (positions) in the method whose type contains multipart.FileHeader.

	fieldIndexOfFileUpload → for each argument identified above, a list of field indices (within its type) where multipart.FileHeader appears.

Exampe:

	type BaseFiles struct {
		Db *sql.DB
		OfficeFile multipart.FileHeader <- [1]
	}
	type MyFiles struct {

		Description string
		File1 multipart.FileHeader <-[1]
		Code string
		FileCode multipart.FileHeader <-[3]
		BaseFile <-- has file upload when detect this struct, index of field is [4,1]
	}
	And the meyhod looks like


	function Update(Files MyFiles, user User, DesFile multipart.FileHeader)
							^					^
							|					|
						   [0]				   [2]
	indexOfArgHasFileUpload=[0,2] <-- Show that the first and third arg in function is multipart.FileHeader
	However, the infomation in "indexOfArgHasFileUpload" is not enought so the function add more infomation in
	"fieldIndexOfFileUpload"
	fieldIndexOfFileUpload=[([1],[3],[4,1])<-- path to fields of struct ],[// direct in arg]]
*/
func (apiMethod *apiMethodInfo) HasUploadFile() bool {
	ret := false
	for i := 0; i < apiMethod.method.Type.NumIn(); i++ {
		inputType := apiMethod.method.Type.In(i)
		if inputType.Kind() == reflect.Ptr {
			inputType = inputType.Elem()
		}
		print(inputType.String())
		check, checkIndex := api.CheckHasInputFile(inputType)
		if check {
			apiMethod.indexOfArgHasFileUpload = append(apiMethod.indexOfArgHasFileUpload, i)
			apiMethod.fieldIndexOfFileUpload = append(apiMethod.fieldIndexOfFileUpload, checkIndex)

		}

		ret = ret || check

	}

	return ret

}
