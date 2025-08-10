package vapi

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type webHandler struct {
	routePath string
	apiInfo   handlerInfo
	initFunc  reflect.Value
	method    string
	Index     int
}

var handlerList []webHandler = []webHandler{}

type webHandlerRunnerType struct {
}

var webHandlerRunner = &webHandlerRunnerType{}

func (web *webHandlerRunnerType) Exec(handler webHandler, w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("exec %s:\n %s\n", handler.routePath, r.RequestURI)
	if r.Method != handler.method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
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

	if strings.HasPrefix(contentType, "application/json") {
		return web.ExecJson(handler, w, r)
	}
	return nil

}
