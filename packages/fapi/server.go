package fapi

import (
	"net/http"
)

type HtttpServer struct {

	// Port is the port the server will listen on.
	Port int
	// Host is the host the server will listen on.
	Host string
	// Handler is the HTTP handler for the server.
	handler http.Handler
	// server is the underlying http.Server.
	server *http.Server
	mux    *http.ServeMux

	mws []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
