package htttpserver

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	vapiErr "wx/errors"
	"wx/handlers"
)

var mapRoutes map[string]WebHandler = map[string]WebHandler{}

func (s *HtttpServer) loadController() error {

	for i := range HandlerList {
		// HandlerList[i].ApiInfo.BaseUrl = s.BaseUrl

		HandlerList[i].Index = i

		if HandlerList[i].ApiInfo.UriHandler == "" || HandlerList[i].ApiInfo.UriHandler[0] != '/' {
			url := s.BaseUrl + "/" + HandlerList[i].ApiInfo.UriHandler
			HandlerList[i].RoutePath = url
			HandlerList[i].RoutePath = strings.ReplaceAll(HandlerList[i].RoutePath, "//", "/")
			HandlerList[i].RoutePath = strings.TrimSuffix(HandlerList[i].RoutePath, "/")
			if HandlerList[i].ApiInfo.IsRegexHandler {
				HandlerList[i].RoutePath += "/"
			}
		} else {

			HandlerList[i].RoutePath = strings.ReplaceAll(HandlerList[i].RoutePath, "//", "/")
			HandlerList[i].RoutePath = strings.TrimSuffix(HandlerList[i].RoutePath, "/")
			if HandlerList[i].ApiInfo.IsRegexHandler {
				HandlerList[i].RoutePath += "/"
			}

		}
		if HandlerList[i].ApiInfo.IsRegexHandler {
			uriRegex := s.BaseUrl + "/"
			uriRegex = handlers.Helper.EscapeSpecialCharsForRegex(uriRegex)
			RegexUri := HandlerList[i].ApiInfo.RegexUri
			RegexUri = strings.TrimPrefix(RegexUri, "^")
			fullRegex := uriRegex + RegexUri
			fullRegex = strings.ReplaceAll(fullRegex, "\\/", "/")
			fullRegex = strings.ReplaceAll(fullRegex, "/", "\\/")

			reg, err := regexp.Compile(fullRegex)
			if err != nil {
				return err
			}
			HandlerList[i].ApiInfo.RegexUriFind = *reg
			if !HandlerList[i].ApiInfo.IsAbsUri {
				HandlerList[i].RoutePath = s.BaseUrl + HandlerList[i].ApiInfo.UriHandler
			}

		} else {
			if !HandlerList[i].ApiInfo.IsAbsUri {
				HandlerList[i].RoutePath = s.BaseUrl + HandlerList[i].ApiInfo.UriHandler
			}

		}
	}
	sort.Slice(HandlerList, func(i, j int) bool {
		return len(HandlerList[i].RoutePath) > len(HandlerList[j].RoutePath) // lớn hơn đứng trước
	})

	for _, h := range HandlerList {
		mapRoutes[h.RoutePath] = h
		fmt.Println(h.RoutePath)
		s.mux.HandleFunc(h.RoutePath, func(w http.ResponseWriter, r *http.Request) {
			err := webHandlerRunner.Exec(h, w, r)
			if err != nil {
				var badReqErr *vapiErr.BadRequestError
				if errors.As(err, &badReqErr) {

					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				var paramMisMatchErr *vapiErr.ParamMissMatchError
				if errors.As(err, &paramMisMatchErr) {

					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				var serviceInitError *vapiErr.ServiceInitError
				if errors.As(err, &serviceInitError) {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				var requireError *vapiErr.RequireError
				if errors.As(err, &requireError) {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		})
	}
	return nil

	// sort handlerList by len of routePath

}
