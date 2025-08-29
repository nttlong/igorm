package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	// giả sử wx của bạn import ở đây
	"wx"

	"github.com/go-chi/chi"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Payload test
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

//go test -bench ^BenchmarkWxPostUser$ -benchmem -memprofile mem_wx.out .
//go tool pprof mem_wx.out
// go test -bench ^BenchmarkWxPostUser$ -benchmem -cpuprofile cpu_wx_create_user.out -memprofile mem_wx_create_user.out .
//go tool pprof mem_wx_create_user.out
// go test -bench ^BenchmarkWxPostUser2$ -benchmem -cpuprofile cpu_wx_create_user2.out -memprofile mem_wx_create_user2.out .
// -------- Benchmark với chi --------

// -------- Benchmark với wx --------
func BenchmarkWxPostUser(b *testing.B) {
	methodInfo := wx.GetMethodByName[TestApi]("CreateUser")
	info, err := wx.Helper.GetHandlerInfo(*methodInfo)
	if err != nil {
		b.Fatal(err)
	}

	mockBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	mockBuilder.PostJson(info.UriHandler, User{ID: 1, Name: "abc"})
	req, _ := mockBuilder.Build()

	// api := &TestApi{} // struct giống bạn test với wx
	// router := wx.NewRouter()
	// router.Register(api)

	// body, _ := json.Marshal(User{ID: 1, Name: "abc"})
	// req := httptest.NewRequest("POST", "/user", bytes.NewReader(body))
	// req.Header.Set("Content-Type", "application/json")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mockBuilder.HandlerRequest(req, func(w http.ResponseWriter, r *http.Request) {
			wx.Helper.ReqExec.DoJsonPost(*info, r, w)
		})
	}
}
func BenchmarkWxPostUser2(b *testing.B) {
	methodInfo := wx.GetMethodByName[TestApi]("CreateUser2")
	info, err := wx.Helper.GetHandlerInfo(*methodInfo)
	if err != nil {
		b.Fatal(err)
	}

	mockBuilder := wx.Helper.ReqExec.CreateMockRequestBuilder()
	mockBuilder.PostJson(info.UriHandler, User{ID: 1, Name: "abc"})
	req, _ := mockBuilder.Build()

	// api := &TestApi{} // struct giống bạn test với wx
	// router := wx.NewRouter()
	// router.Register(api)

	// body, _ := json.Marshal(User{ID: 1, Name: "abc"})
	// req := httptest.NewRequest("POST", "/user", bytes.NewReader(body))
	// req.Header.Set("Content-Type", "application/json")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mockBuilder.HandlerRequest(req, func(w http.ResponseWriter, r *http.Request) {
			wx.Helper.ReqExec.DoJsonPost(*info, r, w)
		})
	}
}
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var body UserInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := fmt.Sprintf("User %s, %d tuổi, sống ở %s",
		body.Name, body.Age, body.Address.City)
	w.Write([]byte(result))
}
func BenchmarkHttpNetPurePostUser(b *testing.B) {
	// Tạo ServeMux thuần net/http
	mux := http.NewServeMux()
	mux.HandleFunc("/user", CreateUserHandler)

	// Request mẫu
	body, _ := json.Marshal(UserInput{
		Name: "abc",
		Age:  30,
		Address: Address{
			City: "Hanoi",
		},
	})
	req := httptest.NewRequest("POST", "/user", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		mux.ServeHTTP(rr, req)
	}
}
func BenchmarkWxPostUserNew(b *testing.B) {
	methodInfo := wx.GetMethodByName[TestApi]("Hello")
	info, _ := wx.Helper.GetHandlerInfo(*methodInfo)
	var handler = func(w http.ResponseWriter, r *http.Request) {
		wx.Helper.ReqExec.Invoke(*info, r, w)
	}
	// Tạo ServeMux thuần net/http
	mux := http.NewServeMux()
	mux.HandleFunc("/"+info.UriHandler, handler)

	// Request mẫu
	body, _ := json.Marshal(UserInput{
		Name: "abc",
		Age:  30,
		Address: Address{
			City: "Hanoi",
		},
	})
	req := httptest.NewRequest("GET", "/"+info.UriHandler+"dasda/dasd", bytes.NewReader(body))
	//req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	for i := 0; i < 100; i++ {

		mux.ServeHTTP(rr, req)
	}
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		mux.ServeHTTP(rr, req)
	}
}
func TestWxUpload(t *testing.T) {
	h, err := wx.Helper.Mock.AddHandler(reflect.TypeOf(TestApi{}), "Upload")
	assert.NoError(t, err)
	assert.NotNil(t, h)
	req := wx.Helper.Mock.NewFormRequest()
	req.AddPhysicalFile("File", `D:\code\go\news2\igorm\comparr-framework\wx_hello\bm_test.go`)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.Handler)
	r, err := req.GetRequest("POST", h.Url)
	assert.NoError(t, err)
	handler.ServeHTTP(rr, r)
}
func BenchmarkWxUpload(b *testing.B) {
	h, _ := wx.Helper.Mock.AddHandler(reflect.TypeOf(TestApi{}), "Upload")

	formReqBuilder := wx.Helper.Mock.NewFormRequest()
	formReqBuilder.AddPhysicalFile("File", `D:\code\go\news2\igorm\comparr-framework\wx_hello\bm_test.go`)
	rr := httptest.NewRecorder()

	req, _ := formReqBuilder.GetRequest("POST", h.Url)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		http.HandlerFunc(h.Handler).ServeHTTP(rr, req)
	}
}
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Giới hạn dung lượng tối đa 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	// Lấy file từ form field "file"
	file, _, err := r.FormFile("File") // bein r nay la CHI no wrapp lai hau dung nguyen goc cua http/net

	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()
}

func TestChiUploadFile(t *testing.T) {
	router := chi.NewRouter()
	router.Post("/upload", UploadHandler)

	formReqBuilder := wx.Helper.Mock.NewFormRequest()
	formReqBuilder.AddPhysicalFile("File", `D:\code\go\news2\igorm\comparr-framework\wx_hello\bm_test.go`)
	rr := httptest.NewRecorder()

	req, _ := formReqBuilder.GetRequest("POST", "/upload")
	router.ServeHTTP(rr, req)
}
func BenchmarkChiUploadFile(b *testing.B) {
	router := chi.NewRouter()
	router.Post("/upload", UploadHandler)

	formReqBuilder := wx.Helper.Mock.NewFormRequest()
	formReqBuilder.AddPhysicalFile("File", `D:\code\go\news2\igorm\comparr-framework\wx_hello\bm_test.go`)
	rr := httptest.NewRecorder()

	req, _ := formReqBuilder.GetRequest("POST", "/upload")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(rr, req)
	}
}
func FiberUploadHandler(c *fiber.Ctx) error {
	// Lấy file từ form field "File"
	fileHeader, err := c.FormFile("File") // "File" là tên field
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Error retrieving file: %v", err))
	}

	// Mở file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Error opening file: %v", err))
	}
	defer file.Close()

	// TODO: xử lý file ở đây (ví dụ lưu vào disk hoặc đọc nội dung)
	// ví dụ: fileHeader.Save("./uploads/" + fileHeader.Filename)

	return c.SendString("File uploaded successfully")
}
func TestFiberUploadHandler(t *testing.T) {
	app := fiber.New()
	app.Post("/upload", FiberUploadHandler)

	// Tạo buffer để chứa form file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Tạo field "File" với nội dung fake
	part, err := writer.CreateFormFile("File", `D:\code\go\news2\igorm\comparr-framework\wx_hello\bm_test.go`)
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("Hello World")) // nội dung file

	writer.Close()

	// Tạo request HTTP fake
	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Gọi Fiber app
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}
func BenchmarkFiberUploadHandler(b *testing.B) {
	app := fiber.New()
	app.Post("/upload", FiberUploadHandler)

	// Tạo buffer để chứa form file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Tạo field "File" với nội dung fake
	part, _ := writer.CreateFormFile("File", `D:\code\go\news2\igorm\comparr-framework\wx_hello\bm_test.go`)

	part.Write([]byte("Hello World")) // nội dung file

	writer.Close()

	// Tạo request HTTP fake
	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	b.ReportAllocs()
	b.ResetTimer()

	// Gọi Fiber app
	for i := 0; i < b.N; i++ {
		app.Test(req)
	}

}
