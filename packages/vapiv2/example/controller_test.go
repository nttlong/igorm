package example

import (
	"testing"
	"vapi"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	mt := vapi.GetMethodByName[Media]("Upload")
	vapi.Helper.FindFormUploadInType((*mt).Type.In(2))
	mtInfo, err := vapi.Helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	t.Log(mtInfo)
}
