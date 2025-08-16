package wx

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"wx/services"
)

type ServiceContext struct {
	Req *http.Request
	Res http.ResponseWriter
	Ctx context.Context // dùng để truyền dữ liệu giữa middleware
}

// Lấy URL tuyệt đối
func (s *ServiceContext) GetAbsUrl() string {
	if s.Req.URL == nil {
		return ""
	}
	scheme := "http"
	if s.Req.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + s.Req.Host + s.Req.URL.RequestURI()
}

// Lấy query string value
func (s *ServiceContext) Query(key string) string {
	return s.Req.URL.Query().Get(key)
}

// Lấy tất cả query params
func (s *ServiceContext) QueryParams() url.Values {
	return s.Req.URL.Query()
}

// Lấy header
func (s *ServiceContext) Header(key string) string {
	return s.Req.Header.Get(key)
}

// Lấy cookie
func (s *ServiceContext) Cookie(name string) (*http.Cookie, error) {
	return s.Req.Cookie(name)
}

// Đọc body dạng []byte
func (s *ServiceContext) BodyBytes() ([]byte, error) {
	return io.ReadAll(s.Req.Body)
}

// Parse JSON body vào struct
func (s *ServiceContext) BindJSON(v any) error {
	return json.NewDecoder(s.Req.Body).Decode(v)
}

// Gửi JSON response
func (s *ServiceContext) JSON(status int, v any) error {
	s.Res.Header().Set("Content-Type", "application/json")
	s.Res.WriteHeader(status)
	return json.NewEncoder(s.Res).Encode(v)
}

// Gửi text response
func (s *ServiceContext) Text(status int, text string) {
	s.Res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	s.Res.WriteHeader(status)
	_, _ = s.Res.Write([]byte(text))
}

// Set header response
func (s *ServiceContext) SetHeader(key, value string) {
	s.Res.Header().Set(key, value)
}

// Set status code
func (s *ServiceContext) Status(code int) {
	s.Res.WriteHeader(code)
}

func init() {
	services.NewServiceContext = func(req *http.Request, res http.ResponseWriter) interface{} {
		return NewServiceContext(req, res)

	}

}
