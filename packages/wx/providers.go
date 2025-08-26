package wx

import (
	"fmt"
	"net/http"
	"reflect"
	"sync"
)

type ProviderContext struct {
	Req *http.Request
	Res http.ResponseWriter
}
type Provider[TImplement any, TInterface any] struct {
	Provider *TInterface
	Context  *ProviderContext
}
type registerEntry struct {
	implemnetType reflect.Type
	fn            reflect.Value
}

var cacheProviderRegister map[reflect.Type]registerEntry = map[reflect.Type]registerEntry{}

func (p *Provider[TImplement, TInterface]) Register(fn func(svc *TImplement) (TInterface, error)) {
	cacheProviderRegister[reflect.TypeFor[TInterface]()] = registerEntry{
		implemnetType: reflect.TypeFor[TImplement](),
		fn:            reflect.ValueOf(fn),
	}
}

type initNewOnce[TImplement any, TInterface any] struct {
	once     sync.Once
	instance TInterface
	err      error
}

var cachNewOnce sync.Map

func (p *Provider[TImplement, TInterface]) NewOnce() (TInterface, error) {
	actual, _ := cachNewOnce.LoadOrStore(reflect.TypeFor[TInterface](), &initNewOnce[TImplement, TInterface]{})
	once := actual.(*initNewOnce[TImplement, TInterface])
	once.once.Do(func() {
		once.instance, once.err = p.New()
	})
	return once.instance, once.err

}
func (p *Provider[TImplement, TInterface]) New() (TInterface, error) {
	entry, ok := cacheProviderRegister[reflect.TypeFor[TInterface]()]
	if !ok {
		//var ret TInterface
		ret := new(TInterface)

		return *ret, fmt.Errorf("%s not register", reflect.TypeFor[TInterface]().String())
	}
	implementValue := reflect.New(entry.implemnetType)
	ret := entry.fn.Call([]reflect.Value{implementValue})
	if ret[1].IsNil() {
		return ret[0].Interface().(TInterface), nil
	}
	return implementValue.Interface().(TInterface), ret[1].Interface().(error)

}
