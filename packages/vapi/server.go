package vapi

import (
	"fmt"
	"net/http"
	"regexp"
)

type HtttpServer struct {
	// Port is the port the server will listen on.
	Port int
	// Host is the host the server will listen on.
	Host string
	// Handler is the HTTP handler for the server.
	handler http.Handler
	// server is the underlying http.Server.
	server           *http.Server
	mux              *http.ServeMux
	onValidateToken  func(token string) TokenValidationResponse
	onGetAccessToken func(w http.ResponseWriter, r *http.Request)
	publicUrl        []*regexp.Regexp
	urlPassAuth      map[string]bool
	urlAuth          string
}

func NewHtttpServer(port int, host string) *HtttpServer {
	return &HtttpServer{
		Port: port,
		Host: host,

		mux:         http.NewServeMux(),
		publicUrl:   []*regexp.Regexp{},
		urlPassAuth: map[string]bool{},
	}

}
func (s *HtttpServer) Start() error {
	if s.handler == nil {
		s.handler = s.mux
	}

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.Host, s.Port), s.handler)
	if err != nil {
		return err
	}
	return nil
}

func Start(auth OAuth2Password, mux *http.ServeMux, bind, port string) error {
	fmt.Printf("server will be start at %s:%s is ok\n", bind, port)
	fmt.Printf("Swagger UI is running at http://%s:%s/swagger/index.html\n", bind, port)
	var handler http.Handler

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		mux.ServeHTTP(w, r)

	})
	if auth != nil {
		mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
			auth.OnGetAccessToken(w, r)
		})
		handler = AuthMiddleware(handler, auth.OnValidateToken)
	}
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", bind, port), handler)
	if err != nil {
		fmt.Printf("start server at %s:%s failed", bind, port)
		return err
	}

	return nil
}
