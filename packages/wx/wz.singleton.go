package wx

import (
	"reflect"
	"wx/internal"
)

type Global[TInstance any] struct {
	ins    *TInstance
	err    error
	fnInit func() (*TInstance, error)
}

func (g *Global[TInstance]) Ins() (*TInstance, error) {
	key := reflect.TypeFor[TInstance]()
	ret, er := internal.OnceCall(key, func() (*TInstance, error) {
		if g.fnInit == nil {
			retVal, err := Helper.DependNewOnce(reflect.TypeFor[TInstance]())
			if err != nil {
				return nil, err
			}
			return retVal.Interface().(*TInstance), nil
		}

		return g.fnInit()
	})
	return ret, er
}
func (g *Global[TInstance]) Init(fn func() (*TInstance, error)) {

	g.fnInit = fn
}
