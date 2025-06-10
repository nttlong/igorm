package caller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"unvs/views"
	_ "unvs/views"

	"unvs/internal/config"

	"github.com/labstack/echo/v4"
)

type CallerHandler struct {
}
type CallerRequest struct {
	ViewId   string                 `json:"viewId"`
	Action   string                 `json:"action"`
	Params   map[string]interface{} `json:"params"`
	Tenant   string                 `json:"tenant"`
	Language string                 `json:"language"`
}
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Response struct for successful account creation
type CallerResponse struct {
	Error   *ErrorResponse `json:"error,omitempty"`
	Results interface{}    `json:"results,omitempty"`
}

// CallerHandler
// @summary CallerHandler
// @description CallerHandler
// @tags caller
// @accept json
// @produce json
// @Param request body CallerRequest true "CallerRequest"
// @router /callers/call [post]
// @Success 201 {object} CallerResponse "Response"
// @Security OAuth2Password
func (h *CallerHandler) Call(c echo.Context) error {

	req := new(CallerRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "INVALID_REQUEST_BODY",
			Message: "Dữ liệu yêu cầu không hợp lệ",
		})
	}
	viewPath := req.ViewId
	action := req.Action
	var JwtInfo *views.JwtDecodeInfo
	if v, ok := views.GetView(viewPath, action); ok {
		// inputTypes := []reflect.Type{}
		if v.IsAuth {
			authHeader := c.Request().Header.Get("Authorization")

			// 2. Kiểm tra xem header có rỗng hay không
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{
					Code:    "UNAUTHORIZED",
					Message: "Missing authorization header",
				})
			}

			// 3. Kiểm tra xem token có phải là Bearer token hay không
			// Chuẩn là "Bearer " (có khoảng trắng ở cuối)
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{
					Code:    "UNAUTHORIZED",
					Message: "Missing authorization header",
				})
			}

			// 4. Trích xuất token bằng cách cắt bỏ tiền tố "Bearer "
			token := strings.TrimPrefix(authHeader, "Bearer ")
			jWTSecretKey := config.GetJWTSecret()
			jwtInfo, err := jwtDecode(token, string(jWTSecretKey))

			if err != nil {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{
					Code:    "UNAUTHORIZED",
					Message: "Invalid token",
				})
			}
			JwtInfo = jwtInfo
		}
		inputValues := []reflect.Value{}

		methodType := v.Method.Type
		for i := 0; i < methodType.NumIn(); i++ {

			paramType := methodType.In(i)
			if paramType.Kind() == reflect.Ptr && i == 0 {
				paramType = paramType.Elem()
			}

			inputData := reflect.New(paramType)
			if i == 0 {
				if JwtInfo != nil {
					if inputData.Kind() == reflect.Ptr {
						sInputData := inputData.Elem()
						fieldClaim := sInputData.FieldByName("Claim")
						fieldClaim.Set(reflect.ValueOf(*JwtInfo))
						/**
										Db       dbx.DBX
						DbTenant dbx.DBXTenant
						Language string
						*/
						fieldDb := sInputData.FieldByName("Db")
						fieldDbTenant := sInputData.FieldByName("DbTenant")
						db := getDbx()
						fieldDb.Set(reflect.ValueOf(db))
						dbTenant, err := getTenantDb(req.Tenant)
						if err != nil {
							return c.JSON(http.StatusUnauthorized, ErrorResponse{
								Code:    "UNAUTHORIZED",
								Message: "Invalid tenant",
							})
						}
						fieldDbTenant.Set(reflect.ValueOf(*dbTenant))
						fieldLanguage := sInputData.FieldByName("Language")
						fieldLanguage.Set(reflect.ValueOf(req.Language))
						fieldContext := sInputData.FieldByName("Context")
						reqContext := c.Request().Context()
						fieldContext.Set(reflect.ValueOf(reqContext))
						// fieldDb.Set(reflect.ValueOf(db))
						// fieldDbTenant.Set(reflect.ValueOf(dbTenant))

						fmt.Println(fieldClaim.Kind())
					}
				}

			} else {
				inputData = inputData.Elem()

			}

			inputValues = append(inputValues, inputData)
			// inputTypes = append(inputTypes, methodType.In(i))
		}
		jsonData, err := json.Marshal(req.Params)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    "INVALID_REQUEST_BODY",
				Message: "Dữ liệu yêu cầu không hợp lệ",
			})
		}
		err = json.Unmarshal(jsonData, inputValues[1].Addr().Interface())
		if err != nil {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Code:    "INVALID_REQUEST_BODY",
				Message: "Dữ liệu yêu cầu không hợp lệ",
			})
		}
		if JwtInfo != nil {
			fx := inputValues[0].Interface()
			fmt.Println(fx)
			fieldClaim := inputValues[0].Elem().FieldByName("Claim")
			fieldClaim.Set(reflect.ValueOf(*JwtInfo))
			// inputValues[0].Set(reflect.ValueOf(JwtInfo))
		}
		results := v.Method.Func.Call(inputValues)
		if len(results) == 2 {
			if results[1].Interface() != nil {
				return c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    "INTERNAL_SERVER_ERROR",
					Message: "Internal server error",
				})
			} else {
				return c.JSON(http.StatusOK, CallerResponse{
					Results: results[0].Interface(),
				})
			}
		}

		return c.JSON(http.StatusOK, CallerResponse{
			Results: results[0].Interface(),
		})
	}
	return c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    "NOT_FOUND",
		Message: "",
	})
}
