package dbx

import (
	"fmt"
	"reflect"
)

func (q QrBuilder[T]) Delete() error {
	if q.where == "" {
		return fmt.Errorf("where clause is required")
	}
	entityType, err := newEntityType(reflect.TypeFor[T]())
	if err != nil {
		return err
	}
	sqlDelete := "DELETE FROM " + entityType.TableName + " WHERE " + q.where
	_, err = q.dbx.Exec(sqlDelete)
	return err
}
