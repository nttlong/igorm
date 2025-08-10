package example

import (
	"testing"
	"vapi"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	mt := vapi.GetMethodByName[Media]("File")
	//vapi.Helper.GetAuthClaims((*mt).Type.In(2))
	mtInfo, err := vapi.Helper.GetHandlerInfo(*mt)
	vapi.Controller(func() (*Media, error) {
		return &Media{}, nil
	})
	assert.NoError(t, err)
	t.Log(mtInfo)
}
