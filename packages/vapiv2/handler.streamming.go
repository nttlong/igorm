package vapi

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (ctx *Handler) StreamingFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		http.Error(ctx.Res, "File not found", http.StatusNotFound)
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(ctx.Res, "Cannot read file info", http.StatusInternalServerError)
		return err
	}

	fileSize := stat.Size()

	// Xác định MIME type chuẩn cho media
	ext := strings.ToLower(filepath.Ext(fileName))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		// fallback: detect by reading first 512 bytes
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		contentType = http.DetectContentType(buf[:n])
		file.Seek(0, io.SeekStart)
	}

	ctx.Res.Header().Set("Content-Type", contentType)
	ctx.Res.Header().Set("Cache-Control", "public, max-age=86400") // cache 1 ngày
	ctx.Res.Header().Set("Accept-Ranges", "bytes")

	rangeHeader := ctx.Req.Header.Get("Range")
	if rangeHeader == "" {
		// Không yêu cầu partial → trả nguyên file
		ctx.Res.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
		ctx.Res.WriteHeader(http.StatusOK)
		_, err = io.Copy(ctx.Res, file)
		return err
	}

	// Có yêu cầu Range → trả partial content
	start, end, err := parseRange(rangeHeader, fileSize)
	if err != nil {
		http.Error(ctx.Res, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
		return err
	}

	ctx.Res.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
	ctx.Res.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	ctx.Res.WriteHeader(http.StatusPartialContent)

	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = io.CopyN(ctx.Res, file, end-start+1)
	return err
}

func (ctx *Handler) StreamingFile2(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		http.Error(ctx.Res, "File not found", http.StatusNotFound)
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(ctx.Res, "Cannot read file info", http.StatusInternalServerError)
		return err
	}

	fileSize := stat.Size()
	ctx.Res.Header().Set("Content-Type", mime.TypeByExtension(fileName))
	ctx.Res.Header().Set("Cache-Control", "public, max-age=86400")
	ctx.Res.Header().Set("Accept-Ranges", "bytes")

	rangeHeader := ctx.Req.Header.Get("Range")
	if rangeHeader == "" {
		// Không có yêu cầu tải từng phần → trả full file
		ctx.Res.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
		ctx.Res.WriteHeader(http.StatusPartialContent)
		_, err = io.Copy(ctx.Res, file)
		return err
	}

	// Có yêu cầu tải một phần file
	start, end, err := parseRange(rangeHeader, fileSize)
	if err != nil {
		http.Error(ctx.Res, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
		return err
	}

	// Set header cho partial content
	ctx.Res.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
	ctx.Res.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	ctx.Res.WriteHeader(http.StatusPartialContent)

	// Seek tới vị trí start
	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		return err
	}

	// Copy đúng số byte cần gửi
	_, err = io.CopyN(ctx.Res, file, end-start+1)
	return err
}

// parseRange nhận "Range: bytes=start-end" và trả về start, end
func parseRange(rangeHeader string, fileSize int64) (int64, int64, error) {
	var start, end int64
	n, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
	if err != nil {
		return 0, fileSize - 1, nil
	}
	if n == 1 {
		// Chỉ có start → end là cuối file
		end = fileSize - 1
	}
	if start > end || end >= fileSize {
		return 0, 0, fmt.Errorf("invalid range")
	}
	return start, end, nil
}
