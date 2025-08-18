package wx

import (
	"fmt"
	"reflect"
	"strings"
	httpServer "wx/HtttpServer"
	handler "wx/handlers"
	"wx/internal"

	"wx/services"
)

type AuthClaims struct {
	handler.AuthClaims
}
type UserClaims struct {
	handler.UserClaims
}

type Handler struct {
	handler.Handler
}

type ControllerContext struct {
	httpServer.ContetxService
}

var NewHtttpServer = httpServer.NewHtttpServer

type SwaggerContact struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	URL   string `json:"url,omitempty"`
}
type SwaggerInfo struct {
	Title       string          `json:"title"`
	Description string          `json:"description,omitempty"`
	Version     string          `json:"version"`
	Contact     *SwaggerContact `json:"contact,omitempty"`
}
type SwaggerBuild struct {
	httpServer.SwaggerBuild
}

func (sb *SwaggerBuild) Info(info SwaggerInfo) *SwaggerBuild {
	var contact *httpServer.SwaggerContact
	if info.Contact != nil {
		contact = &httpServer.SwaggerContact{
			Name:  info.Contact.Name,
			Email: info.Contact.Email,
			URL:   info.Contact.URL,
		}
	}
	sb.SwaggerBuild.Info(httpServer.SwaggerInfo{
		Title:       info.Title,
		Description: info.Description,
		Version:     info.Version,
		Contact:     contact,
	})

	return sb
}

/*
Create swagger build

@BaseUri: url of swagger docs
*/
func CreateSwagger(server *httpServer.HtttpServer, BaseUri string) SwaggerBuild {
	if BaseUri[0] != '/' {
		BaseUri = "/" + BaseUri
	}
	return SwaggerBuild{
		httpServer.CreateSwagger(server, BaseUri),
	}

}
func GetMethodByName[T any](name string) *reflect.Method {

	t := reflect.TypeFor[*T]()
	key := t.Elem().String() + "/GetMethodByName/" + name
	ret, _ := internal.OnceCall(key, func() (*reflect.Method, error) {

		for i := 0; i < t.NumMethod(); i++ {
			if t.Method(i).Name == name {
				ret := t.Method(i)
				return &ret, nil
			}
		}
		return nil, nil
	})
	return ret
}
func GetUriOfHandler[T any](server *httpServer.HtttpServer, methodName string) (string, error) {
	mt := GetMethodByName[T](methodName)
	if mt == nil {
		return "", fmt.Errorf("%s of %T was not found", methodName, *new(T))
	}
	mtInfo, err := handler.Helper.GetHandlerInfo(*mt)
	if err != nil {
		return "", fmt.Errorf("%s of %T cause  error %s", methodName, *new(T), err.Error())
	}
	if mtInfo == nil {
		return "", fmt.Errorf("%s of %T is not HttpMethod", methodName, *new(T))
	}
	if mtInfo.Uri != "" && mtInfo.Uri[0] == '/' {
		return mtInfo.Uri, nil
	}
	return server.BaseUrl + "/" + mtInfo.Uri, nil

}
func init() {
	services.GetCheckScopeTypeName = func() string {
		return strings.Split(reflect.TypeOf(Scoped[any]{}).String(), "[")[0] + "["
	}
	services.GetCheckSingletonTypeName = func() string {
		return strings.Split(reflect.TypeOf(Singleton[any]{}).String(), "[")[0] + "["

	}
}

var Helper = handler.Helper
var HandlerList = httpServer.HandlerList
