package fapi

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

type Inject[T any] struct {
	Val T
}

func (i *Inject[T]) Get() T {
	return i.Val
}

type Singleton[T any] struct {
	Val T
}

type AppConfig struct {
	Port int
	Host string
}
type ObjTest struct {
}
type ApiMethodContext struct {
	*Context `route:"method:get;uri:download/{FileName}.mp4"`
	FileName string
}
type TestFile struct {
	Image *multipart.FileHeader
}
type ProfileUpdate struct {
	*Context    `route:"method:post;uri:profile/{ProfileId}"`
	ProfileId   string `json:"profileId"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
	//
	Test        TestFile
	AvatarFiles []*multipart.FileHeader
}
type ApiMethodContextWithFileInput struct {
	*Context `route:"method:get;uri:download/{FileName}.mp4"`
	FileName string
	File     *multipart.FileHeader
}
type ApiPost struct {
	*Context
	FileName string
}
type TenantContext struct {
	*Context `route:"uri:{Tenant}-@function@"`
	Tenant   string
}

// func (obj *ObjTest) Users(c *TenantContext, userId string) {

// }
func (obj *ObjTest) Upload(c *TenantContext, file bytes.Reader) {

}

//	func (obj *ObjTest) ApiMethod1(c *ProfileUpdate, getUser Inject[Singleton[AppConfig]]) struct {
//		Code string
//		Msg  string
//	} {
//
//		panic("Test")
//	}
func Implements(t, iface reflect.Type) bool {
	if t.Kind() != reflect.Ptr && reflect.PtrTo(t).Implements(iface) {
		return true
	}
	return t.Implements(iface)
}

type Itest interface {
	Get() string
}
type MyTest struct {
	Name string
}

func (m *MyTest) Get() string {
	return m.Name

}
func TestApiUtils(t *testing.T) {
	iface := reflect.TypeOf((*Itest)(nil)).Elem()
	mtx := &MyTest{}

	// Giả sử bạn muốn kiểm tra kiểu *os.File có implement không
	implType := reflect.TypeOf(mtx) // hoặc *os.File nếu dùng os.Open

	// Kiểm tra
	if implType.Implements(iface) {
		fmt.Println("YES: implements multipart.File")
	} else {
		fmt.Println("NO: does not implement multipart.File")
	}

	api.CheckTypeIsContextType(reflect.TypeOf(""))
	// ok, idx := api.CheckHasInputFile(reflect.TypeOf(ProfileUpdate{}))
	// fmt.Println(ok)
	// fmt.Print(idx)
	typ := reflect.TypeOf(&ObjTest{})
	mt := api.InspectMethod(typ.Method(0))
	assert.NotEmpty(t, mt)
	assert.Equal(t, "ApiMethod1", mt.method.Name)
	assert.Equal(t, false, mt.IsAbsUri())
	uri := mt.GetRoute()
	fmt.Println(uri)
	assert.Equal(t, true, mt.HasUploadFile())
	assert.Equal(t, true, mt.HasUploadFile())

}
