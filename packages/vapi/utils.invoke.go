package vapi

import (
	"fmt"
	"net/http"
	"reflect"
)

func (u *utils) AddInitHandler(typ reflect.Type, fn interface{}) {

}
func (u *utils) Invoke(routePath string, w http.ResponseWriter, r *http.Request) {
	fmt.Printf("exec %s\n", routePath)
	w.Header().Set("Content-Type", "application/json")
	if methodInfo, ok := u.mapUrlMethodInfo[routePath]; ok {
		fmt.Println(methodInfo)

	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}
