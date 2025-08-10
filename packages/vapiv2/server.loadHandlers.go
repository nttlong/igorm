package vapi

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

var mapRoutes map[string]webHandler = map[string]webHandler{}

func (s *HtttpServer) loadController() {

	for i := range handlerList {
		handlerList[i].Index = i
		if handlerList[i].apiInfo.UriHandler[0] != '/' {
			url := s.BaseUrl + "/" + handlerList[i].apiInfo.UriHandler
			handlerList[i].routePath = url
			handlerList[i].routePath = strings.ReplaceAll(handlerList[i].routePath, "//", "/")
			handlerList[i].routePath = strings.TrimSuffix(handlerList[i].routePath, "/")
			if handlerList[i].apiInfo.IsRegexhadler {
				handlerList[i].routePath += "/"
			}
		} else {

			handlerList[i].routePath = strings.ReplaceAll(handlerList[i].routePath, "//", "/")
			handlerList[i].routePath = strings.TrimSuffix(handlerList[i].routePath, "/")
			if handlerList[i].apiInfo.IsRegexhadler {
				handlerList[i].routePath += "/"
			}

		}
	}
	sort.Slice(handlerList, func(i, j int) bool {
		return len(handlerList[i].routePath) > len(handlerList[j].routePath) // lớn hơn đứng trước
	})

	for _, h := range handlerList {
		mapRoutes[h.routePath] = h
		fmt.Println(h.routePath)
		s.mux.HandleFunc(h.routePath, func(w http.ResponseWriter, r *http.Request) {
			err := webHandlerRunner.Exec(h, w, r)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)

				return
			}

		})
	}
	// sort handlerList by len of routePath

}
