package fapi

import (
	"net/http"
)

type HtttpServer struct {

	// Port is the port the server will listen on.
	Port int
	// BaseUrl is the base URL of the server.
	BaseUrl string
	// Host is the host the server will listen on.
	Bind string
	// Handler is the HTTP handler for the server.
	handler http.Handler
	// server is the underlying http.Server.
	server *http.Server
	mux    *http.ServeMux

	mws []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

func NewHtttpServer(baseUrl string, port int, bind string) *HtttpServer {
	mux := http.NewServeMux()
	return &HtttpServer{
		Port:    port,
		Bind:    bind,
		BaseUrl: baseUrl,

		mux: mux,
		mws: []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){},
	}

}
