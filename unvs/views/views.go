package views

import (
	"dbx"
	"reflect"
	"strings"
	"sync"
	"unicode"

	"github.com/golang-jwt/jwt/v4"
)

var viewsCache = sync.Map{}

type BaseView struct {
	ViewPath string
	IsAuth   bool
	Claim    JwtDecodeInfo
	Db       dbx.DBX
	DbTenant dbx.DBXTenant
	Language string
}
type ViewCacheType struct {
	Path     string
	Method   reflect.Method
	ViewType reflect.Type
	IsAuth   bool
}
type JwtDecodeInfo struct {
	// Ví dụ về các claim riêng của bạn
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	TenantID string `json:"tenant_id,omitempty"` // Nếu bạn có multi-tenant

	jwt.RegisteredClaims // Nhúng RegisteredClaims để có các claim chuẩn như Exp, Iss, Sub, etc.
}

func AddView(viewlist ...interface{}) {
	for _, view := range viewlist {
		//_viewTyp := reflect.TypeOf(view)
		viewTyp := reflect.TypeOf(view)

		if viewTyp.Kind() == reflect.Ptr {
			viewTyp = viewTyp.Elem()
		}
		viewPath := ""
		IsAuth := false
		for i := 0; i < viewTyp.NumField(); i++ {
			field := viewTyp.Field(i)
			if field.Anonymous {
				v := reflect.ValueOf(view).Elem().Field(i).Interface()
				if obj, ok := v.(BaseView); ok {

					viewPath = obj.ViewPath
					IsAuth = obj.IsAuth
					break

				}
				//fmt.Println(field.Name)

			}
		}
		if viewPath == "" {
			continue
		} else {
			for i := 0; i < reflect.TypeOf(view).NumMethod(); i++ {

				fullViewPath := viewPath + "/" + reflect.TypeOf(view).Method(i).Name
				fullViewPath = strings.ToLower(fullViewPath)
				mt := reflect.TypeOf(view).Method(i)
				mtName := mt.Name
				if !unicode.IsUpper(rune(mtName[0])) {
					continue
				}
				viewsCache.Store(
					fullViewPath,
					&ViewCacheType{
						Path:     viewPath,
						Method:   reflect.TypeOf(view).Method(i),
						ViewType: viewTyp,
						IsAuth:   IsAuth,
					},
				)
			}
		}

	}

}
func GetView(viewPath string, method string) (*ViewCacheType, bool) {
	fullViewPath := viewPath + "/" + method
	fullViewPath = strings.ToLower(fullViewPath)
	if v, ok := viewsCache.Load(fullViewPath); ok {
		return v.(*ViewCacheType), true
	}
	return nil, false

}
