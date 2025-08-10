package vapi

import (
	"mime/multipart"
	"net/http"
	"reflect"
)

type inspectorType struct {
	helper *helperType
}

var inspector *inspectorType
var Helper *helperType

func (*helperType) IsDetectableType(typ reflect.Type, AvailableType ...reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	for _, x := range AvailableType {
		if typ == x {
			return true
		}
	}
	
	return !inspector.helper.IgnoreDetectTypes[typ]
}
func init() {
	inspector = &inspectorType{
		helper: &helperType{
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
				reflect.TypeOf(http.Server{}):   true,
			},
		},
	}
	Helper = inspector.helper
}
