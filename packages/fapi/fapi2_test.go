package fapi

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
