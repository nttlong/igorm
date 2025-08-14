package example

import (
	"testing"
	"vapi"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	mt := vapi.GetMethodByName[Media]("ListOfFiles")
	// server := vapi.NewHtttpServer("/api/v1", 8080, "localhost")
	// url, err := vapi.GetUriOfHandler[Auth](server, "Oauth")
	// assert.NoError(t, err)
	// t.Log(url)
	//vapi.Helper.GetAuthClaims((*mt).Type.In(2))
	// for i := 0; i < (*mt).Type.NumIn(); i++ {
	// 	ok := vapi.Helper.IsInjector((*mt).Type.In(i))
	// 	t.Log(ok)

	// 	// c, err := vapi.Helper.FindHandlerFieldIndexFormType((*mt).Type.In(3))
	// 	// t.Log(c, err)
	// }
	// fmt.Println((*mt).Type.In(3).String())
	// ok := vapi.Helper.IsInjector((*mt).Type.In(3))
	// t.Log(ok)
	vapi.Helper.FindHandlerFieldIndexFormType((*mt).Func.Type().In(1))
	mtInfo, err := vapi.Helper.GetHandlerInfo(*mt)
	vapi.Controller(func() (*Media, error) {
		return &Media{}, nil
	})
	assert.NoError(t, err)
	t.Log(mtInfo)
}
