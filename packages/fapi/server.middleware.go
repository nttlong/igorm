package fapi

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func (s *HtttpServer) Middleware(fn func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) *HtttpServer {
	s.mws = append(s.mws, fn)
	return s
}

var Cors = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Cho phép tất cả origin (cẩn thận với sản phẩm thật!)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Nếu là preflight request (OPTIONS), chỉ phản hồi 200
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Gọi tiếp handler chính
	next.ServeHTTP(w, r)
}
var Zip = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		// Client không hỗ trợ gzip
		next.ServeHTTP(w, r)
		return
	}

	// Gửi header báo là đã nén gzip
	w.Header().Set("Content-Encoding", "gzip")

	gz := gzip.NewWriter(w)
	defer gz.Close()

	gzrw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
	next.ServeHTTP(gzrw, r)
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
