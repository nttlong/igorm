package htttpserver

import (
	"fmt"
	"net/http"
	"wx/handlers"
)

var mapRoutes map[string]WebHandler = map[string]WebHandler{}

func (s *HtttpServer) loadController() error {
	for _, x := range handlers.Helper.Routes.UriList {
		fmt.Println("Registering route:", x)
		s.mux.HandleFunc(x, func(w http.ResponseWriter, r *http.Request) {
			route := handlers.Helper.Routes.Data[x]
			data, err := handlers.Helper.ReqExec.Invoke(route.Info, r, w)
			handlers.Helper.ReqExec.ProcesHttp(route.Info, data, err, r, w)

		})

	}
	return nil

}
