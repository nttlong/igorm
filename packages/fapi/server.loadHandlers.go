package fapi

import (
	"fmt"
	"net/http"
)

var mapRoutes map[string]handlerInfo

func (s *HtttpServer) loadController() {
	for _, h := range handlerList {
		mapRoutes[s.BaseUrl+h.routePath] = h
		s.mux.HandleFunc(s.BaseUrl+h.routePath, func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(s.BaseUrl + h.routePath)
			fmt.Println(r.RequestURI)
		})
	}

}
func init() {
	mapRoutes = map[string]handlerInfo{}
}
