package vapi

import (
	"fmt"
	"mime/multipart"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct {
}

func (a *A) A(h *Handler) {

}
func (a *A) B(h *struct {
	Handler
}) {
}

type TenantHandler struct {
	Tenant  string
	Handler `route:"uri:{Tenant}"`
}
type TenantHandlerGet struct {
	TenantHandler `route:"metho:get"`
}

func (a *A) C(h *TenantHandler) {

}
func (a *A) D(h *struct {
	Auth             *AuthClaims
	UserName         string
	TenantHandlerGet `route:"uri:{fileName}.mp4/@;method:get"`
}, test Inject[string], data *struct {
	Items struct {
		File *[]multipart.FileHeader
	}
}, auth *AuthClaims) {

}
func (a *A) FileUpload(h *Handler, data *struct {
	File *multipart.FileHeader
}) {

}
func (a *A) FileUploadTenant(h *TenantHandler, data *struct {
	Items struct {
		File *multipart.FileHeader
	}
}) {

}
func (a *A) Support() {}
func TestA(t *testing.T) {
	mt := GetMethodByName[A]("A")
	mtInfo, err := inspector.helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	assert.Equal(t, []int{}, mtInfo.FieldIndex)
	assert.Equal(t, 1, mtInfo.IndexOfArg)

}
func TestB(t *testing.T) {
	mt := GetMethodByName[A]("B")
	mtInfo, err := inspector.helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	assert.Equal(t, []int{0}, mtInfo.FieldIndex)
	assert.Equal(t, 1, mtInfo.IndexOfArg)

}
func TestC(t *testing.T) {
	mt := GetMethodByName[A]("C")
	mtInfo, err := inspector.helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	assert.Equal(t, []int{1}, mtInfo.FieldIndex)
	assert.Equal(t, 1, mtInfo.IndexOfArg)

}
func TestD(t *testing.T) {
	mt := GetMethodByName[A]("D")
	mtInfo, err := inspector.helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	fmt.Println(mtInfo.RegexUri)
	testStr := "and/ok.mp4/d"

	re := regexp.MustCompile(mtInfo.RegexUri)
	match := re.MatchString(testStr)
	items := re.FindStringSubmatch(testStr)
	t.Log(items)
	t.Log(match)

	assert.Empty(t, mtInfo)

}
func TestTagsA(t *testing.T) {
	mt := GetMethodByName[A]("D")
	mtInfo, err := inspector.helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	argType := mt.Type.In(mtInfo.IndexOfArg)
	tags := inspector.helper.ExtractTags(argType, mtInfo.FieldIndex)
	
	assert.Equal(t, []string{"", "uri:{Tenant}/*"}, tags)
	uri := inspector.helper.ExtractUriFromTags(tags)
	assert.Equal(t, "/file", uri)
	uriParams := inspector.helper.ExtractUriParams(uri)
	assert.Equal(t, []string{"Tenant"}, uriParams)

}
func TestFileUpload(t *testing.T) {
	mt := GetMethodByName[A]("FileUpload")
	mtInfo, err := inspector.helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	t.Log(mtInfo)
}
func TestFileUploadTenant(t *testing.T) {
	mt := GetMethodByName[A]("FileUploadTenant")
	mtInfo, err := inspector.helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	t.Log(mtInfo)
}
func BenchmarkTestTagsA(t *testing.B) {
	for i := 0; i < t.N; i++ {
		mt := GetMethodByName[A]("D")
		mtInfo, err := inspector.helper.GetHandlerInfo(*mt)
		assert.NoError(t, err)
		argType := mt.Type.In(mtInfo.IndexOfArg)
		tags := inspector.helper.ExtractTags(argType, mtInfo.FieldIndex)

		t.Log(tags)
		uri := inspector.helper.ExtractUriFromTags(tags)
		t.Log(uri)
		uriParams := inspector.helper.ExtractUriParams(uri)
		t.Log(uriParams)

	}

}
