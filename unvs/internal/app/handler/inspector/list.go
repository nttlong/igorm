package inspector

import (
	"dynacall"
	"net/http"
	"reflect"

	_ "dynacall"

	"github.com/labstack/echo/v4"
	_ "unvs.br.auth/roles"
	_ "unvs.br.auth/users"
)

type APIEntry struct {
	CallerPath string
	Args       []interface{}
	Results    []interface{}
}
type InspectorResponse struct {
	APIList []APIEntry `json:"apiList"`
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
		APIList: []APIEntry{},
	}
	for _, entry := range dynacall.GetAllCaller() {
		// Do something with the entry
		args := getInputArgs(entry.Method)
		apiEntry := APIEntry{
			CallerPath: entry.CallerPath,
			Args:       args,
		}
		res.APIList = append(res.APIList, apiEntry)
	}
	return c.JSON(http.StatusCreated, res)
}
func getInputArgs(method reflect.Method) []interface{} {
	args := []interface{}{}
	for i := 1; i < method.Type.NumIn(); i++ {
		args = append(args, reflect.New(method.Type.In(i)).Interface())
	}
	return args
}
