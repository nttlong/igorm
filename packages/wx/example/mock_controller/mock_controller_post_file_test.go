package mockcontroller

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"testing"
	"wx"
	"wx/handlers"

	"github.com/stretchr/testify/assert"
)

type FileController struct {
}

func (fc *FileController) Post(ctx *wx.Handler, data struct {
	File multipart.FileHeader
}) (interface{}, error) {
	return nil, nil

}

// CreateMockUploadRequest tạo mock request upload file
func CreateMockUploadRequest(url, fieldName, fileName string, fileContent []byte) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Thêm file giả lập
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, err
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return nil, err
	}

	// Nếu cần, thêm các field khác
	_ = writer.WriteField("description", "test file upload")

	// Close writer để finalize body
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	// Thiết lập header
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestHalderFileControllerPost(t *testing.T) {
	mt := wx.GetMethodByName[FileController]("Post")

	mtInfo, err := handlers.Helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()

	requestBuild.PostForm("/api/"+mtInfo.UriHandler, struct {
		File *multipart.FileHeader
	}{
		File: &multipart.FileHeader{},
	})

	requestBuild.Handler(func(w http.ResponseWriter, r *http.Request) {
		ret, err := wx.Helper.ReqExec.DoFormPost(*mtInfo, r, w)
		assert.NoError(t, err)
		t.Log(ret)

	})

}
func BenchmarkHalderFileControllerPost(t *testing.B) {
	mt := wx.GetMethodByName[FileController]("Post")
	requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()
	mtInfo, _ := handlers.Helper.GetHandlerInfo(*mt)
	requestBuild.PostForm("/api/"+mtInfo.UriHandler, struct {
		File *multipart.FileHeader
	}{
		File: &multipart.FileHeader{},
	})
	//t.ResetTimer()
	t.Run("test", func(t *testing.B) {
		for i := 0; i < t.N; i++ {

			requestBuild.Handler(func(w http.ResponseWriter, r *http.Request) {
				wx.Helper.ReqExec.DoFormPost(*mtInfo, r, w)
			})
		}
	})
	// t.RunParallel(func(pb *testing.PB) {
	// 	for pb.Next() {
	// 		requestBuild.Handler(func(w http.ResponseWriter, r *http.Request) {
	// 			wx.Helper.ReqExec.DoFormPost(*mtInfo, r, w)
	// 		})
	// 	}

	// })

}

type MultiFilePostBody struct {
	File []multipart.FileHeader
}

func (fc *FileController) PostFiles(c *wx.Handler, body MultiFilePostBody) (interface{}, error) {
	return nil, nil
}

func TestHalderFileControllerPostFiles(t *testing.T) {
	mt := wx.GetMethodByName[FileController]("PostFiles")
	mtInfo, err := handlers.Helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()

	requestBuild.PostForm("/api/"+mtInfo.UriHandler, MultiFilePostBody{
		File: make([]multipart.FileHeader, 2),
	})

	requestBuild.Handler(func(w http.ResponseWriter, r *http.Request) {
		ret, err := wx.Helper.ReqExec.DoFormPost(*mtInfo, r, w)
		assert.NoError(t, err)
		t.Log(ret)

	})

}

type MultiFilePostBodyPtr struct {
	File       *[]multipart.FileHeader
	Files      []*multipart.FileHeader
	FolderName string
}

func (fc *FileController) PostFilesPtr(c *wx.Handler, body MultiFilePostBodyPtr) (any, error) {
	return nil, nil
}
func TestHalderFileControllerPostFilesPtr(t *testing.T) {
	mt := wx.GetMethodByName[FileController]("PostFilesPtr")
	mtInfo, err := handlers.Helper.GetHandlerInfo(*mt)
	assert.NoError(t, err)
	requestBuild := wx.Helper.ReqExec.CreateMockRequestBuilder()
	files := make([]multipart.FileHeader, 2)
	requestBuild.PostForm("/api/"+mtInfo.UriHandler, MultiFilePostBodyPtr{
		File:  &files,
		Files: make([]*multipart.FileHeader, 2),
	})

	requestBuild.Handler(func(w http.ResponseWriter, r *http.Request) {
		ret, err := wx.Helper.ReqExec.DoFormPost(*mtInfo, r, w)
		assert.NoError(t, err)
		t.Log(ret)

	})

}
