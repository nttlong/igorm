package htttpserver

import (
	"net/http"
)

func (s *HtttpServer) Middleware(fn func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) *HtttpServer {
	s.mws = append(s.mws, fn)
	return s
}
