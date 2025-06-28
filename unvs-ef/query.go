// Package unvsef provides a type-safe SQL query builder using Go generics.
// It supports multiple SQL dialects, composable expressions, aggregates,
// binary operations, CASE WHEN expressions, JOINs, GROUP BY, HAVING, and more.
package unvsef

import (
	"reflect"
)

type entitiesUtils struct {
}

func (u *entitiesUtils) QueryableFromType(entityType reflect.Type, tableName string) reflect.Value {
	ret := reflect.New(entityType)
	elem := ret.Elem()

	for i := 0; i < elem.NumField(); i++ {
		valField := elem.Field(i)
		ft := entityType.Field(i)

		if ft.Name == "_" {
			continue
		}

		// Đệ quy cho embedded struct
		if ft.Anonymous {
			anonymousValue := u.QueryableFromType(ft.Type, tableName)
			if valField.CanSet() {
				valField.Set(anonymousValue)
			}
			continue
		}

		// Tạo pointer nếu là nil
		if valField.Kind() == reflect.Ptr {
			if valField.IsNil() {
				valField.Set(reflect.New(valField.Type().Elem()))
			}
			valField = valField.Elem()
		}

		// Duyệt field bên trong Field[T] để tìm *dbField
		fieldType := valField.Type()
		for j := 0; j < fieldType.NumField(); j++ {
			subField := fieldType.Field(j)
			if subField.Anonymous && subField.Type == reflect.TypeOf((*DbField)(nil)) {
				dbFieldSlot := valField.Field(j)
				if dbFieldSlot.IsNil() && dbFieldSlot.CanSet() { //<-- kg bao gio vao duoc cho nay
					dbFieldSlot.Set(reflect.ValueOf(&DbField{
						TableName: tableName,
						ColName:   utils.ToSnakeCase(ft.Name),
					}))
				}
				break
			}
		}
	}

	return ret
}

var entityUtils = &entitiesUtils{}

func Queryable[T any]() *T {
	var v T
	typ := reflect.TypeOf(v)
	if typ == nil {
		typ = reflect.TypeOf((*T)(nil)).Elem()
	}
	e := entityUtils.QueryableFromType(typ, utils.TableNameFromStruct(typ))
	ret := e.Interface().(*T)
	return ret
}
