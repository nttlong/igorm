package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// helloHandler là handler cho endpoint /hello
// Cần di chuyển các annotation của Swagger cho endpoint này lên đây,
// NGAY TRƯỚC định nghĩa hàm.
// @tags heathz
// @Summary Lấy một lời chào đơn giản
// @Description Trả về chuỗi "Hello World!"
// @Accept json
// @Produce json
// @Success 200 {string} string "Hello World!"
// @Router /hz [get]
func HzHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
