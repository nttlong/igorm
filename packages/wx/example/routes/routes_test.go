package routes

import (
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

type Controller1 struct {
}
type Controller2 struct {
}

func (c1 *Controller1) Method1(ctx *wx.Handler) (any, error) {
	return nil, nil
}

func (c2 *Controller2) Method1(ctx *wx.Handler) (any, error) {
	return nil, nil
}
func (c2 *Controller2) Method2(ctx *struct {
	wx.Handler `route:"@/{FileName};metod:get"`
}) (any, error) {
	return nil, nil
}
func (c2 *Controller2) Method3(ctx *struct {
	wx.Handler `route:"@/{FileName}?field={FileName};metod:get"`
	FileName   string
}) (any, error) {
	return nil, nil
}
func (c2 *Controller2) Method4(ctx *struct {
	wx.Handler `route:"/@/{FileName}?field={FileName};metod:get"`
	FileName   string
}) (any, error) {
	return nil, nil
}
func TestAddRoute(t *testing.T) {
	baseUri := "/api"
	wx.Helper.Routes.Add(baseUri,
		&Controller1{},
		&Controller2{})
	for _, x := range wx.Helper.Routes.UriList {
		if wx.Helper.Routes.Data[x].Info.IsAbsUri {
			assert.Equal(t, wx.Helper.Routes.Data[x].Info.UriHandler, x)
		} else {
			assert.Equal(t, baseUri+"/"+wx.Helper.Routes.Data[x].Info.UriHandler, x)
		}

	}

}
