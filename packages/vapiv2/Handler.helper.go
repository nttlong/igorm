package vapi

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
)

func (ctx *Handler) WriteFile(fileName string) error {

	// Mở file
	file, err := os.Open(fileName)
	if err != nil {
		http.Error(ctx.Res, "File not found", http.StatusNotFound)
		return err
	}
	defer file.Close()

	// Lấy thông tin file để set header
	stat, err := file.Stat()
	if err != nil {
		http.Error(ctx.Res, "Can not read file", http.StatusInternalServerError)
		return err
	}

	// Xác định Content-Type (image/png, image/jpeg...)

	ctx.Res.Header().Set("Content-Type", mime.TypeByExtension(fileName)) // hoặc dùng http.DetectContentType
	ctx.Res.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	ctx.Res.Header().Set("Cache-Control", "public, max-age=86400") // cache 1 ngày (86400 giây)
	// khai báo MIME type

	ctx.Res.WriteHeader(http.StatusOK)

	// Ghi nội dung file xuống trình duyệt
	_, err = io.Copy(ctx.Res, file)
	if err != nil {
		return err
	}

	return nil
}
