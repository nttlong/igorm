package handlers

import (
	"mime/multipart"
	"net/http"
	"reflect"
)

type helperType struct {
	SpecialCharForRegex string
	IgnoreDetectTypes   map[reflect.Type]bool
	PrefixGenericDepen  string
	ErrorType           reflect.Type
	ReqExec             *RequestExecutor
	Routes              *RouteTypes
	Services            *ServiceType
}

func (h *helperType) Iscontains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

var Helper = &helperType{
	SpecialCharForRegex: "/\\?.$%^*-+",
	IgnoreDetectTypes: map[reflect.Type]bool{
		reflect.TypeOf(int(0)):                 true,
		reflect.TypeOf(int8(0)):                true,
		reflect.TypeOf(int16(0)):               true,
		reflect.TypeOf(int32(0)):               true,
		reflect.TypeOf(int64(0)):               true,
		reflect.TypeOf(uint(0)):                true,
		reflect.TypeOf(uint8(0)):               true,
		reflect.TypeOf(uint16(0)):              true,
		reflect.TypeOf(uint32(0)):              true,
		reflect.TypeOf(uint64(0)):              true,
		reflect.TypeOf(float32(0)):             true,
		reflect.TypeOf(float64(0)):             true,
		reflect.TypeOf(string("")):             true,
		reflect.TypeOf(bool(false)):            true,
		reflect.TypeOf(nil):                    true,
		reflect.TypeOf([]uint8{}):              true,
		reflect.TypeOf([]byte{}):               true,
		reflect.TypeOf(multipart.FileHeader{}): true,

		//----------------------------------------------------------
		reflect.TypeOf([]int{}):     true,
		reflect.TypeOf([]int8{}):    true,
		reflect.TypeOf([]int16{}):   true,
		reflect.TypeOf([]int32{}):   true,
		reflect.TypeOf([]int64{}):   true,
		reflect.TypeOf([]uint{}):    true,
		reflect.TypeOf([]uint8{}):   true,
		reflect.TypeOf([]uint16{}):  true,
		reflect.TypeOf([]uint32{}):  true,
		reflect.TypeOf([]uint64{}):  true,
		reflect.TypeOf([]float32{}): true,
		reflect.TypeOf([]float64{}): true,
		reflect.TypeOf([]string{}):  true,
		reflect.TypeOf([]bool{}):    true,
		//--------------------------------------
		reflect.TypeOf(http.Request{}):  true,
		reflect.TypeOf(http.Response{}): true,
		reflect.TypeOf(http.Cookie{}):   true,
		reflect.TypeOf(http.Cookie{}):   true,
		reflect.TypeOf(http.Client{}):   true,
		//		reflect.TypeOf(http.Server{}):   true,
	},
	ErrorType: reflect.TypeOf((*error)(nil)).Elem(),
	ReqExec:   &RequestExecutor{},
	Routes: &RouteTypes{
		Data:    map[string]RouteItem{},
		UriList: []string{},
	},
	Services: &ServiceType{},
}
