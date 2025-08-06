package fapi

import (
	"fmt"
	"net/http"
)

func (s *HtttpServer) loadController() {
	for _, h := range handlerList {
		s.mux.HandleFunc(s.BaseUrl+h.routePath, func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.RequestURI)
		})
	}
}
