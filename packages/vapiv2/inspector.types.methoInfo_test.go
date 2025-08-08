package vapi

import (
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
	Handler `route:"uri:{Tenant}/@"`
}
type TenantHandlerGet struct {
	TenantHandler `route:"metho:get"`
}

func (a *A) C(h *TenantHandler) {

}
func (a *A) D(h *struct {
	UserName         string
	TenantHandlerGet `route:"uri:media/{fileName}.mp4"`
}) {

}
func TestA(t *testing.T) {
	mt := GetMethodByName[A]("A")
	mtInfo, err := handlerInfoFromMethod(*mt)
	assert.NoError(t, err)
	assert.Equal(t, []int{}, mtInfo.FieldIndex)
	assert.Equal(t, 1, mtInfo.IndexOfArg)

}
func TestB(t *testing.T) {
	mt := GetMethodByName[A]("B")
	mtInfo, err := handlerInfoFromMethod(*mt)
	assert.NoError(t, err)
	assert.Equal(t, []int{0}, mtInfo.FieldIndex)
	assert.Equal(t, 1, mtInfo.IndexOfArg)

}
func TestC(t *testing.T) {
	mt := GetMethodByName[A]("C")
	mtInfo, err := handlerInfoFromMethod(*mt)
	assert.NoError(t, err)
	assert.Equal(t, []int{1}, mtInfo.FieldIndex)
	assert.Equal(t, 1, mtInfo.IndexOfArg)

}
func TestD(t *testing.T) {
	mt := GetMethodByName[A]("D")
	mtInfo, err := handlerInfoFromMethod(*mt)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 1}, mtInfo.FieldIndex)
	assert.Equal(t, 1, mtInfo.IndexOfArg)

}
func TestTagsA(t *testing.T) {
	mt := GetMethodByName[A]("D")
	mtInfo, err := handlerInfoFromMethod(*mt)
	assert.NoError(t, err)
	argType := mt.Type.In(mtInfo.IndexOfArg)
	tags := handlerInfoExtractTags(argType, mtInfo.FieldIndex)

	assert.Equal(t, []string{"", "uri:{Tenant}/*"}, tags)
	uri := handlerInfoExtractUriFromTags(tags)
	assert.Equal(t, "/file", uri)

}
