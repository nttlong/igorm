package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	wxErr "wx/errors"
)

func (reqExec *RequestExecutor) Invoke(info HandlerInfo, r *http.Request, w http.ResponseWriter) (any, error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		// In ra thông báo panic
	// 		// In ra stack trace
	// 		buf := make([]byte, 1<<16)
	// 		runtime.Stack(buf, true) // Lấy toàn bộ stack trace
	// 		fmt.Printf("Stack trace:\n%s\n", buf)
	// 	}
	// }()

	if r.Method != info.HttpMethod {
		return nil, wxErr.NewMethodNotAllowError("method not allowed")
	}

	if r.Method == "GET" {
		return reqExec.DoGet(info, r, w)
	}
	if r.Header.Get("Content-Type") == "application/json" {

		return reqExec.DoJsonPost(info, r, w)
	}
	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {

		return reqExec.DoFormPost(info, r, w)
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data; boundary=") {
		return reqExec.DoFormPost(info, r, w)
	}
	return reqExec.DoJsonPost(info, r, w)

}
func (reqExec *RequestExecutor) handlerError(err error, r *http.Request, w http.ResponseWriter) {

	isShowError := false
	switch err.(type) {
	case *wxErr.UriParamParseError:
		isShowError = true
		http.Error(w, "Bad Request", http.StatusBadRequest)
	case *wxErr.RegexUriNotMatchError:
		if r.URL.Path[len(r.URL.Path)-1] == '/' {
			newPath := strings.TrimSuffix(r.URL.Path, "/")
			http.Redirect(w, r, newPath, http.StatusMovedPermanently) // 301
		} else {
			isShowError = true
			http.Error(w, "not found", http.StatusNotFound)
		}
	case *wxErr.UriParamConvertError:
		isShowError = true
		http.Error(w, "Bad Request", http.StatusBadRequest)
	case *wxErr.MethodNotAllowError:
		isShowError = true
		http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
	case *wxErr.ServiceInitError:
		isShowError = true
		http.Error(w, "Server error", http.StatusInternalServerError)

	default:
		isShowError = true
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
	if isShowError {
		fmt.Printf("Error: %s\n", r.URL.RequestURI())
		fmt.Printf("Error: %v\n", err)
	}
}

func (reqExec *RequestExecutor) ProcesHttp(info HandlerInfo, data interface{}, previousErr error, r *http.Request, w http.ResponseWriter) {
	if previousErr != nil {
		reqExec.handlerError(previousErr, r, w)
		return
	}
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
}
