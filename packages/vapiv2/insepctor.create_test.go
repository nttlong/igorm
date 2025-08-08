package vapi

import (
	"fmt"
	"regexp"
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
func (u *User) Get4(ctx struct {
	HttpGet   `route:"uri:/avatar/{userId}/profile/{profileId}"`
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
func TestCreateMethodGet3(t *testing.T) {
	mt := GetMethodByName[User]("Get3")
	ret, err := inspector.Create(*mt)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
	assert.Equal(t, "get", ret.Route.Method)
	assert.Equal(t, "vapi/user/get3/{userId}/profile/{profileId}", ret.Route.Uri)
	assert.Equal(t, "vapi/user/get3/", ret.Route.UriHandler)
	assert.Equal(t, "vapi\\/user\\/get3\\/(.*)\\/profile\\/(.*)", ret.Route.RegexUri)
	regexMacth := regexp.MustCompile(ret.Route.RegexUri)
	url := "vapi/user/get3/1234/profile/test-001"
	check := regexMacth.MatchString(url)
	assert.Equal(t, true, check)
	matches := regexMacth.FindStringSubmatch(url)
	assert.Equal(t, "1234", matches[1])
	assert.Equal(t, "test-001", matches[2])
	assert.Equal(t, []string{"userId", "profileId"}, ret.Route.UriParams)

}
func TestCreateMethodGet4(t *testing.T) {
	mt := GetMethodByName[User]("Get4")
	ret, err := inspector.Create(*mt)
	if err != nil {
		t.Error(err)
	}
	t.Log(ret)
	assert.Equal(t, "get", ret.Route.Method)
	assert.Equal(t, true, ret.Route.ISAbsUri)
	assert.Equal(t, "avatar/{userId}/profile/{profileId}", ret.Route.Uri)
	assert.Equal(t, "avatar/", ret.Route.UriHandler)
	fmt.Println(ret.Route.RegexUri)
	assert.Equal(t, []string{"userId", "profileId"}, ret.Route.UriParams)
	assert.Equal(t, [][]int{{1}, {2}}, ret.Route.IndexOfFieldInUri)

}
