package vapi

import "net/http"

func (server *HtttpServer) Oauth() {
	handler := server.handler
	if handler == nil {
		handler = server.mux
	}

	server.mux.HandleFunc("/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		if server.onGetAccessToken == nil {
			panic("server.onGetAccessToken is nil, please call server.SetOnGetAccessToken")
		}
		server.onGetAccessToken(w, r)
	})

}
func (server *HtttpServer) OnGetAccessToken(onGetAccessToken func(w http.ResponseWriter, r *http.Request)) {
	server.onGetAccessToken = onGetAccessToken

}
