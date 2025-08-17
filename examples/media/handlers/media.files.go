package handlers

import (
	"fmt"
	"wx"
)

func (media *Media) Files(ctx *struct {
	wx.Handler `route:"uri:@/{Tenant}/{*FilePath};method:get"`
	FilePath   string
	Tenant     string
}) {
	fmt.Println(ctx.FilePath)
}
