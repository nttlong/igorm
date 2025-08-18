package controllers

import (
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

func TestController1Post1(t *testing.T) {
	mt := wx.GetMethodByName[Controller1]("Post1")
	assert.NotEmpty(t, *mt)
	info, err := wx.Helper.GetHandlerInfo(*mt)
	assert.Nil(t, err)
	assert.NotEmpty(t, info)
	assert.Equal(t, "controller1/post1", info.Uri)
	assert.Equal(t, false, info.IsRegexHandler)
	assert.Equal(t, "controller1/post1", info.UriHandler)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, "POST", info.HttpMethod)
	assert.Equal(t, 0, len(info.UriParams))
	assert.Equal(t, 0, len(info.IndexOfInjectors))
	assert.Equal(t, false, info.HasInjector)
	assert.Equal(t, 0, len(info.FormUploadFile))
	assert.Equal(t, -1, info.IndexOfRequestBody)
	assert.Equal(t, 0, len(info.IndexOfAuthClaims))
	assert.Equal(t, -1, info.IndexOfAuthClaimsArg)

	mt2 := wx.GetMethodByName[Controller1]("NoPost")
	assert.NotEmpty(t, *mt2)
	info2, err := wx.Helper.GetHandlerInfo(*mt2)
	assert.Nil(t, err)
	assert.Empty(t, info2)

}
