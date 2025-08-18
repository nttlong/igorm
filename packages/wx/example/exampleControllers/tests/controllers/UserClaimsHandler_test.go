package controllers

import (
	"reflect"
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

func TestUserClaimsHandler(t *testing.T) {
	type test1 struct {
		Test string
		wx.UserClaims
	}

	type test2 struct {
		Test string
	}
	ret1 := wx.Helper.GetUserClaims(reflect.TypeOf(test1{}))

	assert.Equal(t, [][]int{{1, 0}}, ret1)
	ret := wx.Helper.GetUserClaims(reflect.TypeOf(wx.UserClaims{}))

	assert.Equal(t, [][]int{{0}}, ret)
	ret = wx.Helper.GetUserClaims(reflect.TypeOf(&wx.UserClaims{}))

	assert.Equal(t, [][]int{{0}}, ret)

	ret = wx.Helper.GetUserClaims(reflect.TypeOf(test2{}))

	assert.Equal(t, [][]int{}, ret)
	type test3 struct {
		Test string
		UserClaimsHandler
	}
	ret = wx.Helper.GetUserClaims(reflect.TypeOf(test3{}))

	assert.Equal(t, [][]int{{1, 0, 0}}, ret)
	type test4 struct {
		Test string
		Data struct {
			UserClaimsHandler
		}
	}
	ret = wx.Helper.GetUserClaims(reflect.TypeOf(test4{}))

	assert.Equal(t, [][]int{{1, 0, 0, 0}}, ret)
	type RescursiveStruc1 struct {
		*RescursiveStruc1
		UserClaimsHandler
	}
	ret = wx.Helper.GetUserClaims(reflect.TypeOf(RescursiveStruc1{}))

	assert.Equal(t, [][]int{{1, 0, 0}}, ret)
}
func TestGetUserClaimsOnMethod(t *testing.T) {

}
