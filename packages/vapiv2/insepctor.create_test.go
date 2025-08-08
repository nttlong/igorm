package vapi

import "testing"

type User struct {
}

func (u *User) Get(ctx Handler, id Inject[string]) {

}
func (u *User) Get2(ctx Handler, id Inject[string], id2 Inject[string]) {

}
func (u *User) Get3(ctx struct {
	Handler   `route:"uri:{userId}/profile/{profileId}"`
	UserId    string
	ProfileId string
}, id string) {

}

func (u *User) Get4(ctx struct {
	Handler   `route:"uri:/avatar/{userId}/profile/{profileId}"`
	UserId    string
	ProfileId string
}, id string) {

}

type TenantGet struct {
	Handler  `route:"uri:{tenantId}/@"`
	TenantId string
}

func (u *User) GetInfo(ctx TenantGet) {

}
func (u *User) GetUser(ctx struct {
	TenantGet `route:"uri:*/{userId}"`
	TenantId  string
}) {
}
func TestGet(t *testing.T) {
	
}
