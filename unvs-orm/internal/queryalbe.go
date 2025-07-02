package internal

// Package unvsef provides a type-safe SQL query builder using Go generics.
// It supports multiple SQL dialects, composable expressions, aggregates,
// binary operations, CASE WHEN expressions, JOINs, GROUP BY, HAVING, and more.

import (
	"reflect"
	"strings"
	"sync"
)

type entitiesUtils struct {
	cacheGetAllFields sync.Map
	FieldResolver     func(tableName string, colName string, field reflect.StructField) reflect.Value
}

/*
Create a new instance of the queryable type and set the DbField for each field.

# Note: This function return a pointer to the new instance of the queryable type.
*/

// func (u *entitiesUtils) getAllFields(entityType reflect.Type) map[string]reflect.StructField {
// 	if fields, ok := u.cacheGetAllFields.Load(entityType); ok {
// 		return fields.(map[string]reflect.StructField)
// 	}

// 	fields := make(map[string]reflect.StructField)
// 	for i := 0; i < entityType.NumField(); i++ {
// 		ft := entityType.Field(i)
// 		fields[ft.Name] = ft
// 	}

//		u.cacheGetAllFields.Store(entityType, fields)
//		return fields
//	}
func (u *entitiesUtils) QueryableFromType(entityType reflect.Type, tableName string) reflect.Value {

	ret := reflect.New(entityType)
	elem := ret.Elem()
	mapField := utils.GetMetaInfo(entityType)

	for colName, fieldTags := range mapField[tableName] {

		valField := elem.FieldByName(fieldTags.Field.Name)
		ft, ok := entityType.FieldByName(fieldTags.Field.Name)
		if !ok {
			continue
		}

		ftType := ft.Type
		if ftType.Kind() == reflect.Ptr {
			ftType = ftType.Elem()
		}

		if strings.HasPrefix(ftType.String(), utils.entityTypeName) {

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

		// Duyệt field bên trong Field[T] để tìm *DbField

		fv := u.FieldResolver(tableName, colName, ft)

		valField.Set(fv)

	}

	return ret
}

var entityUtils = &entitiesUtils{}

var EntityUtils = entityUtils
