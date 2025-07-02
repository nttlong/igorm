package unvsef

import "reflect"

type Entity[T any] struct {
	Meta map[string]FieldTag
}

func (e *Entity[T]) New() T {
	entityType := reflect.TypeFor[T]()
	tableName := utils.TableNameFromStruct(entityType)
	retVal := entityUtils.QueryableFromType(entityType, tableName)
	return retVal.Elem().Interface().(T)

}
