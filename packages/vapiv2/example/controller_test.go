package example

import (
	"testing"
	"vapi"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	mt := vapi.GetMethodByName[Auth]("Oauth")
	server := vapi.NewHtttpServer("/api/v1", 8080, "localhost")
	url, err := vapi.GetUriOfHandler[Auth](server, "Oauth")
	assert.NoError(t, err)
	t.Log(url)
	//vapi.Helper.GetAuthClaims((*mt).Type.In(2))
	mtInfo, err := vapi.Helper.GetHandlerInfo(*mt)
	vapi.Controller(func() (*Media, error) {
		return &Media{}, nil
	})
	assert.NoError(t, err)
	t.Log(mtInfo)
}
