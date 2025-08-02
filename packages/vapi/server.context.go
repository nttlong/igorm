package vapi

import "net/http"

type HttpContext struct {
	// Request is the raw http.Request
	Request *http.Request
	// Response is the raw http.ResponseWriter
	Response http.ResponseWriter
	// PathParams holds the path parameters extracted from the URL.
	PathParams map[string]string
	// Claims holds the claims extracted from the JWT token.

}

func newHttpContext(request *http.Request, response http.ResponseWriter) *HttpContext {
	return &HttpContext{
		Request:  request,
		Response: response,
	}

}
