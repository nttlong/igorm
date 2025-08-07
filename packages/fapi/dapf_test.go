package fapi

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct{}

type Tenant struct {
	_      string `route:"uri:{Tenant}-{Username}/@.mp4"`
	Tenant string
}
type TenantRequest struct {
	_      string `route:"uri:{Tenant}/@"`
	Tenant string
}

func (u *User) Save(ctx Context, data *struct {
	TenantRequest
	Username string
	Password string
}) {
	panic("implement me")
}
func (u *User) Update(ctx Context, data struct {
	Tenant
	Username string
	Password string
	Avatar   multipart.FileHeader
}) {
	panic("implement me")
}

// -----------------
func (u *User) Delete(ctx Context, data struct {
	UserId string
}) {
	panic("implement me")
}
func getMethodByName[T any](name string) *reflect.Method {
	t := reflect.TypeFor[*T]()
	for i := 0; i < t.NumMethod(); i++ {
		if t.Method(i).Name == name {
			fmt.Println(t.Method(i))
			ret := t.Method(i)
			return &ret
		}
	}
	return nil
}
func TestInspectMethodUpdate(t *testing.T) {
	mt := getMethodByName[User]("Update")
	ret := inspector.InspectorMethod(*mt)
	assert.Equal(t, "Update", mt.Name)

	tags := ret.Tags()
	assert.Equal(t, "uri:{Tenant}-{Username}/@.mp4", tags)
	route := ret.Route()

	uriParams := route.UriParams()
	assert.Equal(t, []string{"tenant", "username"}, uriParams)

	fmt.Println(route.regexUri)
	assert.Equal(t, "{tenant}-{username}/update.mp4", route.uri)
	typ := mt.Type.In(2)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	for i, x := range route.indexOfFieldInUri {
		assert.Equal(t, uriParams[i], strings.ToLower(typ.FieldByIndex(x).Name))
		// assert.Equal(t, typ.FieldByIndex(x).Name), route.indexOfFieldInUri[x])
	}

}
func TestInspectMethodDelete(t *testing.T) {
	mt := getMethodByName[User]("Delete")
	ret := inspector.InspectorMethod(*mt)
	fmt.Println(ret)
	assert.Equal(t, "Delete", mt.Name)
}
func TestInspectMethodSave(t *testing.T) {
	mt := getMethodByName[User]("Save")
	ret := inspector.InspectorMethod(*mt)
	fmt.Println(ret)
	assert.Equal(t, "Save", mt.Name)
}
