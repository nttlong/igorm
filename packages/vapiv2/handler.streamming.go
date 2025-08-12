package vapi

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

	ext := strings.ToLower(filepath.Ext(fileName))
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		contentType = http.DetectContentType(buf[:n])
		// Reset offset sau khi đọc
		_, _ = file.Seek(0, io.SeekStart)
	}

	ctx.Res.Header().Set("Content-Type", contentType)
	ctx.Res.Header().Set("Cache-Control", "public, max-age=86400")
	ctx.Res.Header().Set("Accept-Ranges", "bytes")

	rangeHeader := ctx.Req.Header.Get("Range")
	flusher, _ := ctx.Res.(http.Flusher)

	if rangeHeader == "" {
		// Trả toàn bộ file, status 200
		ctx.Res.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
		ctx.Res.WriteHeader(http.StatusOK)

		buf := make([]byte, 32*1024) // 32KB
		for {
			n, err := file.Read(buf)
			if n > 0 {
				_, writeErr := ctx.Res.Write(buf[:n])
				if writeErr != nil {
					return writeErr
				}
				if flusher != nil {
					flusher.Flush()
				}
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
		}
		return nil
	}

	// Parse Range header
	start, end, err := parseRange(rangeHeader, fileSize)
	if err != nil {
		http.Error(ctx.Res, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
		return err
	}

	contentLength := end - start + 1

	ctx.Res.Header().Set("Content-Length", fmt.Sprintf("%d", contentLength))
	ctx.Res.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	ctx.Res.WriteHeader(http.StatusPartialContent)

	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		return err
	}

	buf := make([]byte, 32*1024) // 32KB
	remaining := contentLength
	for remaining > 0 {
		readSize := int64(len(buf))
		if remaining < readSize {
			readSize = remaining
		}

		n, err := file.Read(buf[:readSize])
		if n > 0 {
			_, writeErr := ctx.Res.Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			if flusher != nil {
				flusher.Flush()
			}
			remaining -= int64(n)
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}

/*
parseRange parses a Range header of the form "bytes=start-end" and returns
the start and end byte positions as int64.

Note:
- The Range header might be like "bytes=0-" (from start to end of file)
- Or "bytes=12345-" (from position 12345 to end of file)
- The function validates the range according to the fileSize.
- It returns an error if the range is invalid or malformed.

Parameters:
- rangeHeader: the value of the Range header from the HTTP request.
- fileSize: the total size of the file in bytes.

Returns:
- start: the starting byte position of the requested range.
- end: the ending byte position of the requested range.
- error: non-nil if the range is invalid or cannot be parsed.
*/
func parseRange(rangeHeader string, fileSize int64) (int64, int64, error) {
	const prefix = "bytes="
	// Check if the header starts with "bytes="
	if !strings.HasPrefix(rangeHeader, prefix) {
		return 0, 0, fmt.Errorf("invalid range header")
	}

	// Remove the "bytes=" prefix
	r := strings.TrimPrefix(rangeHeader, prefix)
	// Split into start and end parts, e.g. "123-456" -> ["123", "456"]
	items := strings.SplitN(r, "-", 2)
	if len(items) != 2 {
		return 0, 0, fmt.Errorf("invalid range format")
	}

	// Parse start position
	start, errStart := strconv.ParseInt(strings.TrimSpace(items[0]), 10, 64)
	// Parse end position
	end, errEnd := strconv.ParseInt(strings.TrimSpace(items[1]), 10, 64)

	if errStart != nil {
		// If start is not a valid number, return error
		// (Note: suffix ranges like "bytes=-500" are not supported here)
		return 0, 0, fmt.Errorf("invalid start range")
	}

	if errEnd != nil {
		// If end is missing or invalid (e.g. "bytes=12345-"),
		// treat it as the last byte of the file
		end = fileSize - 1
	}

	// Validate range boundaries
	if start < 0 || start > end || end >= fileSize {
		return 0, 0, fmt.Errorf("invalid range values")
	}

	return start, end, nil
}
