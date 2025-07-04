package orm

import (
	_ "unvs-orm/internal"
)

// func reflectFromUnsafePtr(ptr unsafe.Pointer, typ reflect.Type) reflect.Value {
// 	return reflect.NewAt(typ, ptr).Elem()
// }
// func (e *Model[T]) Insert() error {
// 	typ := reflect.TypeOf(e.data)
// 	fmt.Println(typ)
// 	return nil

// }
