package vapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

// Handler trả về JSON
func helloHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "hello world",
	})
}

func Benchmark_HelloWorldEcho(b *testing.B) {
	// Khởi tạo Echo và route
	e := echo.New()
	e.GET("/hello", helloHandler)

	// Tạo request sẵn
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	// Reset timer để bắt đầu benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Gọi handler trực tiếp
		if err := helloHandler(c); err != nil {
			b.Error(err)
		}
	}
}
