package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
)

type MockRequestBuilder struct {
	method string
	url    string
	body   io.Reader
	header map[string]string
	forms  map[string]string
}

func (builder *MockRequestBuilder) Build() (*http.Request, http.ResponseWriter) {
	ret, err := http.NewRequest(builder.method, builder.url, builder.body)
	if err != nil {
		panic(err)
	}
	if builder.header != nil {
		for k, v := range builder.header {
			ret.Header.Add(k, v)
		}
	}
	if builder.forms != nil {
		for k, v := range builder.forms {
			ret.Form.Add(k, v)
		}
	}

	ret.Host = "localhost"

	ret.URL = &url.URL{
		Scheme:  "http",
		Host:    "localhost",
		Path:    "/" + strings.Split(builder.url, "://")[1],
		RawPath: builder.url,
	}
	return ret, builder.NewResponse()

}
func (builder *MockRequestBuilder) PostJson(url string, data interface{}) *MockRequestBuilder {

	builder.method = "POST"
	builder.url = "http://localhost" + url
	if builder.header == nil {
		builder.header = make(map[string]string)
	}
	builder.header["Content-Type"] = "application/json"
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		bodyData := strings.NewReader(string(jsonData))
		builder.body = bodyData
	}
	return builder

}
func (builder *MockRequestBuilder) NewResponse() http.ResponseWriter {
	ret := newMockResponseWriter()
	return ret

}
func (reqExec *RequestExecutor) CreateMockRequestBuilder() *MockRequestBuilder {
	return &MockRequestBuilder{}

}
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

func (builder *MockRequestBuilder) PostForm(url string, data interface{}) *MockRequestBuilder {

	builder.method = "POST"
	builder.url = "http://localhost" + url
	if builder.header == nil {
		builder.header = make(map[string]string)
	}
	builder.header["Content-Type"] = "multipart/form-data"
	if data != nil {
		typData := reflect.TypeOf(data)
		valData := reflect.ValueOf(data)
		if typData.Kind() == reflect.Ptr {
			typData = typData.Elem()
			valData = valData.Elem()
		}
		for i := 0; i < typData.NumField(); i++ {

			field := typData.Field(i)
			if field.Type == reflect.TypeOf(multipart.FileHeader{}) {
				fileBody, contentType, err := CreateMockUploadRequest(strings.ToLower(field.Name), fmt.Sprintf("file %d", i), []byte("mock content of file"))
				if err != nil {
					panic(err)
				}
				builder.body = fileBody
				builder.header["Content-Type"] = contentType
				continue

			}
			if field.Type == reflect.TypeOf(&multipart.FileHeader{}) {
				fieldVal := valData.Field(i)
				var nilFile *multipart.FileHeader = nil

				if fieldVal != reflect.ValueOf(nilFile) {

					fileBody, contentType, err := CreateMockUploadRequest(strings.ToLower(field.Name), fmt.Sprintf("file %d", i), []byte("mock content of file"))
					if err != nil {
						panic(err)
					}
					builder.header["Content-Type"] = contentType
					builder.body = fileBody
				}

				continue

			}
			fieldVal := valData.Field(i)
			if fieldVal.Kind() == reflect.Ptr {
				{
					fieldVal = fieldVal.Elem()
				}
			}
			val := fieldVal.Interface()
			if val == nil {
				continue
			}
			if fieldVal.Kind() == reflect.Struct {
				bff, err := json.Marshal(val)
				if err != nil {
					panic(err)
				}
				builder.forms[field.Name] = string(bff)
				continue
			}
			if fieldVal.Kind() == reflect.String {
				builder.forms[field.Name] = val.(string)
				continue
			}

		}
		// jsonData, err := json.Marshal(data)
		// if err != nil {
		// 	panic(err)
		// }
		// bodyData := strings.NewReader(string(jsonData))
		// builder.body = bodyData
	}
	return builder

}
func (builder *MockRequestBuilder) Handler(handler http.HandlerFunc) {
	rr := httptest.NewRecorder()
	req, _ := builder.Build()
	// gọi handler

	handler.ServeHTTP(rr, req)
}
