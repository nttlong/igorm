package handlers

import (
	"bytes"
	"mime/multipart"
)

func CreateMockUploadRequest(fieldName, fileName string, fileContent []byte) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Thêm file giả lập
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, "", err
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return nil, "", err
	}

	// Nếu cần, thêm các field khác
	_ = writer.WriteField("description", "test file upload")

	// Close writer để finalize body
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	// Thiết lập header
	//req.Header.Set("Content-Type", writer.FormDataContentType())
	return body, writer.FormDataContentType(), nil
}
func (builder *MockRequestBuilder) AddFile(fileName string) {
	if builder.body == nil {
		builder.body = new(bytes.Buffer)
	}
	if builder.writer == nil {
		builder.writer = multipart.NewWriter(builder.body)
	}

	part, _ := builder.writer.CreateFormFile(fileName, fileName)

	part.Write([]byte("Mock content of file"))

}
func (builder *MockRequestBuilder) AddField(fieldName string, value string) {
	if builder.body == nil {
		builder.body = new(bytes.Buffer)
	}
	if builder.writer == nil {
		builder.writer = multipart.NewWriter(builder.body)
	}

	part, _ := builder.writer.CreateFormField(fieldName)

	part.Write([]byte(value))

}
func (builder *MockRequestBuilder) AddJsonField(fieldName string, value []byte) {
	if builder.body == nil {
		builder.body = new(bytes.Buffer)
	}
	if builder.writer == nil {
		builder.writer = multipart.NewWriter(builder.body)
	}

	part, _ := builder.writer.CreateFormField(fieldName)

	part.Write(value)

}
