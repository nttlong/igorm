package vapi

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Singleton[T any] struct {
	val T
	New func() T
}

func (s *Singleton[T]) Get() T {
	return s.val
}

type TestController struct {
}
type Handler[T any] func(controller *T) error

type Data struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
type FileUpload struct {
	File *multipart.FileHeader
	Data Data
}
type DownloadData struct {
	Accesskey string `api:"url:{accessKey},description:file name,required:true"`
	FileName  string `api:"url:{fileName}.mp4,description:file name,required:true"`
}

func (t *TestController) Method(data FileUpload, ctx HttpContext) {

	fmt.Println(ctx)
}

// func (t *TestController) Donwload(data DownloadData, ctx HttpContext) {

// 	fmt.Println(ctx)
// }

// func (t *TestController) Donwload_Path_mp4(data Data, ctx HttpContext) {

// 	fmt.Println(ctx)
// }
// func (t *TestController) Donwload_Path_all(data Data, ctx HttpContext) {

//		fmt.Println(ctx)
//	}
func TestParseMethod(t *testing.T) {

	typ := reflect.TypeOf(&TestController{})
	list, err := utilsInstance.ParseMethods(typ)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(list))
	for _, item := range list {
		fmt.Println(item.regexpUrl, " ", item.Url, " ", item.HandlerUrl, " ", item.HttpMethod)
	}

}
func TestAddController(t *testing.T) {
	AddController(func() (*TestController, error) {
		return &TestController{}, nil
	})

}
