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
		if handlerList[i].apiInfo.UriHandler[0] != '/' {
			url := s.BaseUrl + handlerList[i].apiInfo.UriHandler
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
		/*
			mux.HandleFunc("/a/b/c", func(w http.ResponseWriter, r *http.Request) {...}
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				khi goi request ve luc nao no cung avo day ma kg vao mux.HandleFunc("/a/b/c"
			}


		*/
		s.mux.HandleFunc(h.routePath, func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(h.routePath)
			fmt.Println(r.RequestURI)
			// r.ParseMultipartForm(20 * 1024 * 1024)
			// if r.MultipartForm.File != nil {
			// 	for f := range r.MultipartForm.File {
			// 		fmt.Println(f)
			// 		fmt.Println(r.MultipartForm.File[f])

			// 	}

			// }

		})
	}
	// sort handlerList by len of routePath

}
