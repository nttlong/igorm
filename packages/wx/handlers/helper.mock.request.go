package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type mockRequest struct {
	Body   *bytes.Buffer
	Writer *multipart.Writer
	req    *http.Request
	files  []*os.File
	Err    error
}

// AddField thêm field text
func (mr *mockRequest) AddField(key, value string) error {
	return mr.Writer.WriteField(key, value)
}

// AddPhysicalFile thêm file thật từ disk
func (mr *mockRequest) AddPhysicalFile(fieldName, pathToFile string) error {
	file, err := os.Open(pathToFile)
	if err != nil {
		return err
	}
	mr.files = append(mr.files, file)

	part, err := mr.Writer.CreateFormFile(fieldName, file.Name())
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)
	return err
}

// Build hoàn tất multipart và tạo http.Request
func (mr *mockRequest) GetRequest(method, url string) (*http.Request, error) {
	// đóng writer để kết thúc multipart boundary
	err := mr.Writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, mr.Body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mr.Writer.FormDataContentType())
	mr.req = req
	return req, nil
}

// Close giải phóng file đã mở
func (mr *mockRequest) Close() {
	for _, f := range mr.files {
		f.Close()
	}
}
