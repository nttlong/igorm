package vdi

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type metaInfo struct {
	sampleInstanceReflectValue reflect.Value
	mappingFields              map[string]reflect.StructField // Cache injector fields
	sampleFieldValues          map[string]reflect.Value       // Cache sample field values
	initFuncs                  map[string]reflect.Value       // Cache Init functions
}

type containerInfo[T any] struct {
	Instance interface{}
	Value    reflect.Value
	meta     metaInfo
}

func (r *containerInfo[T]) Get() T {
	ret := r.Instance.(T)
	// val := r.contianerVal
	// if val.Kind() == reflect.Ptr {
	// 	val = val.Elem()
	// }

	// v := val.Interface()
	// ret1 := v.(T)

	return ret
}

func (r *containerInfo[T]) init(resolver func(obj T) error) (T, error) {
	typ := reflect.TypeFor[T]()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	_objVal := reflect.New(typ)
	obj := _objVal.Interface().(T)

	err := resolver(obj)

	_objVal = reflect.ValueOf(obj)
	objVal := _objVal
	if objVal.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
	}

	for i := 0; i < objVal.NumField(); i++ {
		field := typ.Field(i)

		if utils.isInjector(field) {

			proValue := objVal.Field(i) // struct embedding
			if proValue.Kind() == reflect.Ptr {
				proValue = proValue.Elem()
			}

			if _, ok := field.Type.FieldByName("Owner"); ok {
				ownerField := proValue.FieldByName("Owner")

				if ownerField.IsValid() && ownerField.CanSet() && ownerField.IsZero() {
					if _objVal.Type().AssignableTo(ownerField.Type()) {
						ownerField.Set(_objVal)
					} else if ownerField.Type().Kind() == reflect.Interface && _objVal.Type().Implements(ownerField.Type()) {
						ownerField.Set(_objVal.Convert(ownerField.Type()))
					}
				}
			}
		} else {
			// r.meta.mappingFields[field.Name] = field
			// r.meta.sampleFieldValues[field.Name] = _objVal.Field(i)
		}
	}

	r.meta.sampleInstanceReflectValue = _objVal

	if err != nil {
		var zero T
		return zero, err
	}
	r.Instance = obj
	r.Value = _objVal

	return obj, nil
}

func (r *containerInfo[T]) GetInitFun(PropertyName string) (*reflect.Value, error) {
	val := reflect.ValueOf(r.Instance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field, ok := val.Type().FieldByName(PropertyName)
	if !ok {
		return nil, fmt.Errorf("field %s not found in type %s", PropertyName, val.Type().String())
	}
	if !utils.isInjector(field) {
		return nil, fmt.Errorf("field %s is not an injector", PropertyName)
	}
	fn := val.FieldByName("Init")

	return &fn, nil

}

// var make sure the containerInfo[T] is only created once for each type T
// var containersCache = sync.Map{}

// var globalMeta : each containerInfo has meta map,
// this value to store reflect.ValueOf(sampleInstance),
var globalMeta = sync.Map{}

func getMetaOfType(typ reflect.Type) *metaInfo {
	key := typ.String()
	if v, ok := globalMeta.Load(key); ok {
		ret := v.(metaInfo)
		return &ret
	}
	return nil
}

// RegisterContainer registers a new container for a given type T.
//
// RegisterContainer ä¸ºæŒ‡å®šç±»åž‹ T æ³¨å†Œä¸€ä¸ªæ–°çš„å®¹å™¨ã€‚
//
// RegisterContainer ã¯æŒ‡å®šã•ã‚ŒãŸåž‹ T ã®ã‚³ãƒ³ãƒ†ãƒŠã‚’ç™»éŒ²ã—ã¾ã™ã€‚
//
// RegisterContainer Ä‘Äƒng kÃ½ container má»›i cho kiá»ƒu T.
type initRegisterContainer struct {
	once     sync.Once
	instance interface{}
	err      error
}

func RegisterContainer[T any](resolver func(svc T) error) (*containerInfo[T], error) {
	key := reflect.TypeFor[T]().String()
	actual, _ := initRegisterContainerCache.LoadOrStore(key, &initRegisterContainer{})
	initContainer := actual.(*initRegisterContainer)
	initContainer.once.Do(func() {
		initContainer.instance, initContainer.err = registerContainer(resolver)
	})
	return initContainer.instance.(*containerInfo[T]), initContainer.err
}
func registerContainer[T any](resolver func(svc T) error) (interface{}, error) {
	var t T

	ret := containerInfo[T]{
		meta: metaInfo{
			mappingFields:              make(map[string]reflect.StructField),
			initFuncs:                  make(map[string]reflect.Value),
			sampleFieldValues:          make(map[string]reflect.Value),
			sampleInstanceReflectValue: reflect.ValueOf(t),
		},
	} // Create new containerInfo instance // åˆ›å»ºæ–°çš„ containerInfo å®žä¾‹ // æ–°ã—ã„ containerInfo ã‚¤ãƒ³ã‚¹ã‚¿ãƒ³ã‚¹ã‚’ä½œæˆ // Táº¡o containerInfo má»›i

	data, err := ret.init(resolver)
	if err != nil {
		return nil, err // Initialization failed // åˆå§‹åŒ–å¤±è´¥ // åˆæœŸåŒ–å¤±æ•— // Khá»Ÿi táº¡o tháº¥t báº¡i
	}
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	//ret.Instance = val.Interface()

	return ret, nil // Return the created container // è¿”å›žåˆ›å»ºçš„å®¹å™¨ // ä½œæˆã•ã‚ŒãŸã‚³ãƒ³ãƒ†ãƒŠã‚’è¿”ã™ // Tráº£ vá» container vá»«a táº¡o
}

type unvsdiUtils struct {
	currentPackage    string
	once              sync.Once
	injectorTypeNames []string
}

func (u *unvsdiUtils) getCurrentPackage() string {
	u.once.Do(func() {
		u.currentPackage = reflect.TypeOf(unvsdiUtils{}).PkgPath()
	})
	return u.currentPackage
}
func (u *unvsdiUtils) isInjector(field reflect.StructField) bool {

	if field.Type.PkgPath() != u.getCurrentPackage() {
		return false
	}
	fieldTypeName := field.Type.String()
	isInjector := false
	for _, injectorTypeName := range u.injectorTypeNames {
		if strings.HasPrefix(fieldTypeName, injectorTypeName+"[") {
			isInjector = true
			break
		}
	}

	return isInjector
}

var utils = &unvsdiUtils{
	injectorTypeNames: []string{ // list of injector type names
		strings.Split(reflect.TypeOf(Singleton[any, any]{}).String(), "[")[0],
		strings.Split(reflect.TypeOf(Scoped[any, any]{}).String(), "[")[0],
		strings.Split(reflect.TypeOf(Transient[any, any]{}).String(), "[")[0],
	},
}
