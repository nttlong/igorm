package wx

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	httpServer "wx/HtttpServer"
	handler "wx/handlers"
	"wx/internal"
)

type AuthClaims struct {
	handler.AuthClaims
}
type UserClaims struct {
	handler.UserClaims
}

type Handler struct {
	handler.Handler
	schema     string
	rootAbsUrl string
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

type initGetMethodByName struct {
	val *reflect.Method

	once sync.Once
}

var cacheGetMethodByName sync.Map

func GetMethodByName[T any](name string) *reflect.Method {
	key := name + "@" + reflect.TypeFor[T]().String()
	actual, _ := cacheGetMethodByName.LoadOrStore(key, &initGetMethodByName{})
	init := actual.(*initGetMethodByName)
	init.once.Do(func() {
		init.val = getMethodByName[T](name)
	})
	return init.val
}
func getMethodByName[T any](name string) *reflect.Method {

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

type initGetUriOfHandler struct {
	val  string
	err  error
	once sync.Once
}

var cacheGetUriOfHandler sync.Map

func GetUriOfHandler[T any](methodName string) (string, error) {
	key := methodName + "@" + reflect.TypeFor[T]().String() + "@"
	actual, _ := cacheGetUriOfHandler.LoadOrStore(key, &initGetUriOfHandler{})
	init := actual.(*initGetUriOfHandler)
	init.once.Do(func() {

		mt := GetMethodByName[T](methodName)
		if mt == nil {
			init.err = fmt.Errorf("%s of %T was not found", methodName, *new(T))
			return
		}
		mtInfo, err := handler.Helper.GetHandlerInfo(*mt)
		if err != nil {
			init.err = fmt.Errorf("%s of %T cause  error %s", methodName, *new(T), err.Error())
			return
		}
		if mtInfo == nil {
			init.err = fmt.Errorf("%s of %T is not HttpMethod", methodName, *new(T))
			return
		}
		if mtInfo.UriHandler != "" && mtInfo.IsAbsUri {
			init.val = strings.TrimSuffix(mtInfo.UriHandler, "/")
			return
		}
		init.val = "/" + strings.TrimSuffix(mtInfo.UriHandler, "/")
	})
	return init.val, init.err

}

// func init() {
// 	services.GetCheckScopeTypeName = func() string {
// 		return strings.Split(reflect.TypeOf(Scoped[any]{}).String(), "[")[0] + "["
// 	}
// 	services.GetCheckSingletonTypeName = func() string {
// 		return strings.Split(reflect.TypeOf(Singleton[any]{}).String(), "[")[0] + "["

// 	}
// }

var Helper = handler.Helper
var HandlerList = httpServer.HandlerList


