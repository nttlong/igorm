package vapi

import (
	"net/http"
	"sync"
)

type Handler struct {
	Res     http.ResponseWriter
	Req     *http.Request
	BaseUrl string
}

type initInspectHttpMethodFromType struct {
	val  string
	once sync.Once
}

var cacheInspectHttpMethodFromType sync.Map
