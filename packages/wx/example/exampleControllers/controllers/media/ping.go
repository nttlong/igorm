package media

import "wx"

func (media *Media) Ping(ctx *wx.Handler) (string, error) {
	return "Hello World", nil
}
func (media *Media) Ping2(ctx *struct {
	wx.Handler `route:"test"`
}) (string, error) {
	return "Hello World", nil
}
func (media *Media) Ping3(ctx *struct {
	wx.Handler `route:"test"`
}) (string, error) {
	return "Hello World", nil
}
