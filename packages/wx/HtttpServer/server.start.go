package htttpserver

import (
	"fmt"
	"net/http"
	"time"
)

var baseUrlOfServer string

func (s *HtttpServer) Start() error {
	baseUrlOfServer = s.BaseUrl
	// Đăng ký các handler vào mux
	err := s.loadController()
	if err != nil {
		return err
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

	addr := fmt.Sprintf("%s:%s", s.Bind, s.Port)
	// fmt.Println("Server listening at", addr)
	// return http.ListenAndServe(addr, s.handler)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.handler,
		ReadTimeout:  10 * time.Second, // Giới hạn đọc request
		WriteTimeout: 10 * time.Second, // Giới hạn ghi response
		IdleTimeout:  60 * time.Second, // Cho keep-alive
	}

	fmt.Println("Server listening at", addr)
	return s.server.ListenAndServe()
}
