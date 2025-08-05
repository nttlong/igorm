package vapi

import (
	"fmt"
	"net/http"
	"regexp"
)

func (u *utils) MapHandler(relUrl string, info *methodInfo) {
	if u.mapUrlMethodInfo == nil {
		u.mapUrlMethodInfo = make(map[string]*methodInfo)
	}
	u.mapUrlMethodInfo[relUrl] = info
}
func (s *HtttpServer) Start() error {
	// Đăng ký các handler vào mux
	for _, h := range utilsInstance.handler {

		regContent := h.regexpUrl.String()
		regContent = "^" + swaggerData.BasePath + regContent
		h.regexpUrl = regexp.MustCompile(regContent)
		routePath := swaggerData.BasePath + h.MasterUrl
		utilsInstance.MapHandler(routePath, h)
		s.mux.HandleFunc(routePath, func(w http.ResponseWriter, r *http.Request) {
			utilsInstance.Invoke(routePath, w, r)
		})
	}

	// handler cuối cùng gọi mux
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mux.ServeHTTP(w, r)
	})

	// Gắn middleware vào handler chain
	for i := len(s.mws) - 1; i >= 0; i-- {
		mw := s.mws[i]
		next := final
		final = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mw(w, r, next.ServeHTTP)
		})
	}

	s.handler = final

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	fmt.Println("Server listening at", addr)
	return http.ListenAndServe(addr, s.handler)
}
