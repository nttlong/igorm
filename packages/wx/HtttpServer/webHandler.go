package htttpserver

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"wx/handlers"
)

type WebHandler struct {
	RoutePath string
	ApiInfo   handlers.HandlerInfo
	InitFunc  reflect.Value
	Method    string
	Index     int
}

type webHandlerRunnerType struct {
}

var webHandlerRunner = &webHandlerRunnerType{}

func (web *webHandlerRunnerType) Exec(handler WebHandler, w http.ResponseWriter, r *http.Request) error {

	fmt.Printf("exec %s:\n %s\n", handler.RoutePath, r.RequestURI)
	if r.Method != handler.Method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}
	if r.Method == "GET" {
		return web.ExecGet(handler, w, r)
	}
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return err
		}
		return web.ExecFormPost(handler, w, r)
	}
	if strings.HasPrefix(contentType, "multipart/form-data") {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
		return web.ExecFormPost(handler, w, r)
	}

	if strings.HasPrefix(contentType, "application/json") || contentType == "" {
		return web.ExecJson(handler, w, r)
	}
	return nil

}
