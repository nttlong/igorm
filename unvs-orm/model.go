package orm

import (
	"reflect"

	internal "unvs-orm/internal"
)

type Model[T any] struct {
	TenantDb *internal.TenantDb
	meta     map[string]internal.FieldTag
	A        string
}

func (e *Model[T]) New() T {
	entityType := reflect.TypeFor[T]()
	tableName := internal.Utils.TableNameFromStruct(entityType)

	retVal := internal.EntityUtils.QueryableFromType(entityType, tableName, nil, e.TenantDb)
	return retVal.Elem().Interface().(T)

}
func Queryable[T any](tenantDb *internal.TenantDb) *T {
	var v T
	typ := reflect.TypeOf(v)
	if typ == nil {
		typ = reflect.TypeOf((*T)(nil)).Elem()
	}

	e := internal.EntityUtils.QueryableFromType(typ, internal.Utils.TableNameFromStruct(typ), nil, tenantDb)
	ret := e.Interface().(*T)
	return ret
}
