package handlers

import (
	"fmt"
	"wx"
)

func (media *Media) Files(ctx *struct {
	wx.Handler `route:"uri:@/{*FilePath};method:get"`
	FilePath   string
}) error {
	fullFIlePath, err := media.File.GetFilePath(media.FileDirectory, ctx.FilePath)
	if err != nil {
		return err
	}
	//ctx.FilePath = fullFIlePath

	fmt.Println(fullFIlePath)
	return ctx.StreamingFile(fullFIlePath)
	// return nil
}
