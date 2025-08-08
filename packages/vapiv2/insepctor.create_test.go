package vapi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
}

func (u *User) Get(ctx HttpGet, id Inject[string]) {

}
func (u *User) Get2(ctx HttpGet, id Inject[string], id2 Inject[string]) {

}
func (u *User) Get3(ctx struct {
	HttpGet   `route:"uri:{userId}/profile/{profileId}"`
	UserId    string
	ProfileId string
}, id string) {

}
func TestCreateMethodGet(t *testing.T) {
	mt := GetMethodByName[User]("Get")
	ret, err := inspector.Create(*mt)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
	assert.Equal(t, "get", ret.Route.Method)
	assert.Equal(t, "vapi/user/get", ret.Route.Uri)
	assert.Equal(t, "vapi/user/get", ret.Route.UriHandler)
	assert.Equal(t, "vapi\\/user\\/get", ret.Route.RegexUri)
	fmt.Println(ret.Route.RegexUri)
	assert.Equal(t, []string{}, ret.Route.UriParams)
	assert.Equal(t, [][]int{}, ret.Route.IndexOfFieldInUri)
	assert.Equal(t, "", ret.Route.Tags)

}
func TestCreateMtheodGet2(t *testing.T) {
	mt := GetMethodByName[User]("Get2")
	ret, err := inspector.Create(*mt)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
	assert.Equal(t, "get", ret.Route.Method)
	assert.Equal(t, "vapi/user/get2", ret.Route.Uri)
	assert.Equal(t, "vapi/user/get2", ret.Route.UriHandler)
	assert.Equal(t, "vapi\\/user\\/get2", ret.Route.RegexUri)
	fmt.Println(ret.Route.RegexUri)
	assert.Equal(t, []string{}, ret.Route.UriParams)
	assert.Equal(t, [][]int{}, ret.Route.IndexOfFieldInUri)
	assert.Equal(t, "", ret.Route.Tags)
}
func TestCreateMtheodGet3(t *testing.T) {
	mt := GetMethodByName[User]("Get3")
	ret, err := inspector.Create(*mt)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
	assert.Equal(t, "get", ret.Route.Method)
	assert.Equal(t, "vapi/user/get3/{userId}/profile/{profileId}", ret.Route.Uri)
}
