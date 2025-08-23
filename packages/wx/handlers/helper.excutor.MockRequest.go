package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	body   *bytes.Buffer
	header map[string]string
	forms  map[string]string
	writer *multipart.Writer
	Req    *http.Request
	Res    httptest.ResponseRecorder
}

func (builder *MockRequestBuilder) Build() (*http.Request, http.ResponseWriter) {

	if builder.writer != nil {
		if err := builder.writer.Close(); err != nil {
			panic(err)
		}
	}
	ret, err := http.NewRequest(builder.method, builder.url, builder.body)
	if err != nil {
		panic(err)
	}
	if builder.writer == nil {
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
	}

	ret.Host = "localhost"

	ret.URL = &url.URL{
		Scheme:  "http",
		Host:    "localhost",
		Path:    "/" + strings.Split(builder.url, "://")[1],
		RawPath: builder.url,
	}
	if builder.writer != nil {
		ret.Header.Set("Content-Type", builder.writer.FormDataContentType())
	}
	// // if builder.body != nil {
	// // 	ret.ContentLength = int64(builder.body.Len())
	// // 	ret.Body.Close()
	// // 	//ret.Body = io.NopCloser(builder.body)
	// // }
	for k, v := range builder.header {
		if k == "Content-Type" {
			continue
		}
		ret.Header.Set(k, v)
	}
	return ret, builder.NewResponse()

}
func (builder *MockRequestBuilder) PostJson(url string, data interface{}) *MockRequestBuilder {
	if builder.body == nil {
		builder.body = new(bytes.Buffer)
	}
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

		builder.body.Write(jsonData)

	}
	return builder

}
func (builder *MockRequestBuilder) Get(url string) *MockRequestBuilder {
	if builder.body == nil {
		builder.body = new(bytes.Buffer)
	}
	builder.method = "GET"
	builder.url = "http://localhost" + url
	if builder.header == nil {
		builder.header = make(map[string]string)
	}
	builder.header["Content-Type"] = "application/json"

	return builder

}
func (builder *MockRequestBuilder) NewResponse() http.ResponseWriter {
	ret := newMockResponseWriter()
	return ret

}
func (reqExec *RequestExecutor) CreateMockRequestBuilder() *MockRequestBuilder {
	return &MockRequestBuilder{}

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
			if field.Type == reflect.TypeOf([]multipart.FileHeader{}) {
				valOfField := valData.Field(i)
				for j := 0; j < valOfField.Len(); j++ {
					builder.AddFile(strings.ToLower(field.Name))
				}
				continue

			}
			if field.Type.Kind() == reflect.Ptr && field.Type.Elem() == reflect.TypeOf([]multipart.FileHeader{}) {
				valOfField := valData.Field(i).Elem()
				for j := 0; j < valOfField.Len(); j++ {
					builder.AddFile(strings.ToLower(field.Name))
				}
				continue

			}
			if field.Type.Kind() == reflect.Ptr && field.Type.Elem() == reflect.TypeOf([]*multipart.FileHeader{}) {
				valOfField := valData.Field(i).Elem()
				for j := 0; j < valOfField.Len(); j++ {
					builder.AddFile(strings.ToLower(field.Name))
				}
				continue

			}
			if field.Type.Kind() == reflect.Slice && field.Type.Elem() == reflect.TypeOf(&multipart.FileHeader{}) {
				valOfField := valData.Field(i)
				for j := 0; j < valOfField.Len(); j++ {
					builder.AddFile(strings.ToLower(field.Name))
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
				builder.AddJsonField(field.Name, bff)
				continue
			}

			if fieldVal.Kind() == reflect.String {

				builder.AddField(field.Name, val.(string))
				continue
			}

		}

	}
	return builder

}
func (builder *MockRequestBuilder) Handler(handler http.HandlerFunc) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	req, _ := builder.Build()
	// gọi handler

	handler.ServeHTTP(rr, req)
	return rr

}

func (builder *MockRequestBuilder) ServerHandler(fn func() (any, error)) {
	rr := httptest.NewRecorder()
	req, _ := builder.Build()

	// var handler http.HandlerFunc
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		builder.Res = *w.(*httptest.ResponseRecorder)
		builder.Req = r
		data, err := fn()

		if err != nil {
			Helper.ReqExec.handlerError(err, r, w)

			// print("Error: ", err.Error(), "\n")
			builder.Res = *w.(*httptest.ResponseRecorder)

			return
		}
		if data != nil {
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		builder.Res = *w.(*httptest.ResponseRecorder)
	})
	handler.ServeHTTP(rr, req)
}
