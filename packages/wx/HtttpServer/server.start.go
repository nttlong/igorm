package htttpserver

import (
	"fmt"
	"net/http"
)

func (s *HtttpServer) Start() error {
	// Đăng ký các handler vào mux
	s.loadController()
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

	addr := fmt.Sprintf("%s:%d", s.Bind, s.Port)
	fmt.Println("Server listening at", addr)
	return http.ListenAndServe(addr, s.handler)
}
