package vapi

import (
	"sync"
)

type Handler struct {
}
type initInspectHttpMethodFromType struct {
	val  string
	once sync.Once
}

var cacheInspectHttpMethodFromType sync.Map
