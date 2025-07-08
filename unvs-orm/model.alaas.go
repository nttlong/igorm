package orm

import (
	"reflect"
	"sync"
)

var proxyAliasCache sync.Map

func (m *Model[T]) Alias(aliasName string) T {
	// check proxy alias cache
	if alias, ok := proxyAliasCache.Load(aliasName); ok {
		return alias.(T)
	}

	typ := reflect.TypeFor[T]()
	var modelVal *reflect.Value = nil
	queryableVal := EntityUtils.QueryableFromTypeNoCache(typ, m.TableName+"*"+aliasName, modelVal)
	ret := queryableVal.Interface().(*T)
	proxyAliasCache.Store(aliasName, *ret)
	return *ret
}
