package handlers

import (
	"net/http"
	"strings"
	wxErr "wx/errors"
)

func (reqExec *RequestExecutor) Invoke(info HandlerInfo, r *http.Request, w http.ResponseWriter) (interface{}, error) {
	if r.Method != info.HttpMethod {
		return nil, wxErr.NewMethodNotAllowError("method not allowed")
	}

	if r.Header.Get("Content-Type") != "application/json" {

		return reqExec.DoJsonPost(info, r, w)
	}
	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {

		return reqExec.DoFormPost(info, r, w)
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return reqExec.DoFormPost(info, r, w)
	}
	return reqExec.DoJsonPost(info, r, w)

}
