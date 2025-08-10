package vapi

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

var mapRoutes map[string]webHandler = map[string]webHandler{}

func (s *HtttpServer) loadController() error {

	for i := range handlerList {
		handlerList[i].Index = i

		if handlerList[i].apiInfo.UriHandler == "" || handlerList[i].apiInfo.UriHandler[0] != '/' {
			url := s.BaseUrl + "/" + handlerList[i].apiInfo.UriHandler
			handlerList[i].routePath = url
			handlerList[i].routePath = strings.ReplaceAll(handlerList[i].routePath, "//", "/")
			handlerList[i].routePath = strings.TrimSuffix(handlerList[i].routePath, "/")
			if handlerList[i].apiInfo.IsRegexHandler {
				handlerList[i].routePath += "/"
			}
		} else {

			handlerList[i].routePath = strings.ReplaceAll(handlerList[i].routePath, "//", "/")
			handlerList[i].routePath = strings.TrimSuffix(handlerList[i].routePath, "/")
			if handlerList[i].apiInfo.IsRegexHandler {
				handlerList[i].routePath += "/"
			}

		}
		if handlerList[i].apiInfo.IsRegexHandler {
			uriRegex := s.BaseUrl + "/"
			uriRegex = inspector.helper.EscapeSpecialCharsForRegex(uriRegex)
			RegexUri := handlerList[i].apiInfo.RegexUri
			RegexUri = strings.TrimPrefix(RegexUri, "^")
			fullRegex := uriRegex + RegexUri
			fullRegex = strings.ReplaceAll(fullRegex, "\\/", "/")
			fullRegex = strings.ReplaceAll(fullRegex, "/", "\\/")

			reg, err := regexp.Compile(fullRegex)
			if err != nil {
				return err
			}
			handlerList[i].apiInfo.RegexUriFind = *reg

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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		})
	}
	return nil

	// sort handlerList by len of routePath

}
