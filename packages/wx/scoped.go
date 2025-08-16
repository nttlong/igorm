package wx

import (
	"fmt"
	"net/http"
	"reflect"
)

type Scoped[T any] struct {
	ins  *T
	err  error
	Ctx  *ServiceContext
	init func(ctx *ServiceContext) (*T, error)
}

func (t *Scoped[T]) Init(fn func(ctx *ServiceContext) (*T, error)) {
	t.init = fn
}
func (t *Scoped[T]) GetInstance() (*T, error) {
	if t.init == nil {
		return nil, fmt.Errorf("%s not initialized,please call Init() of %s first", reflect.TypeOf(t).String(), reflect.TypeOf(t).String())
	}
	r, err := t.init(t.Ctx)
	return r, err

}

func NewServiceContext(req *http.Request, res http.ResponseWriter) *ServiceContext {
	return &ServiceContext{
		Req: req,
		Res: res,
	}

}
