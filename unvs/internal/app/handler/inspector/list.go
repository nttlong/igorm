package inspector

import (
	"dynacall"
	"net/http"

	_ "dynacall"

	"github.com/labstack/echo/v4"
	_ "unvs.br.auth/roles"
	_ "unvs.br.auth/users"
)

type InspectorResponse struct {
	APIList []string `json:"apiList"`
}
type InspectorHandler struct {
}

// helloHandler là handler cho endpoint /hello
// Cần di chuyển các annotation của Swagger cho endpoint này lên đây,
// NGAY TRƯỚC định nghĩa hàm.
// @tags System
// @Summary Query all api action and domain
// @Description Query all api action and domain
// @Accept json
// @Produce json
// @Success 200 {object} InspectorResponse
// @Router /inspector/list [post]
func (h *InspectorHandler) List(c echo.Context) error {
	res := &InspectorResponse{
		APIList: []string{},
	}
	for _, entry := range dynacall.GetAllCaller() {
		// Do something with the entry
		res.APIList = append(res.APIList, entry.CallerPath)
	}
	return c.JSON(http.StatusCreated, res)
}
