package example

import (
	"testing"
	"wx"
	"wx/handlers"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	mt := wx.GetMethodByName[Media]("ListOfFiles")

	handlers.Helper.FindHandlerFieldIndexFormType((*mt).Func.Type().In(1))
	mtInfo, err := handlers.Helper.GetHandlerInfo(*mt)
	wx.LoadController(func() (*Media, error) {
		return &Media{}, nil
	})
	assert.NoError(t, err)
	t.Log(mtInfo)
}
