/*
This file is part of dynacall.
Declares the function to create dynamic struct from a list of reflect.Type.
*/
package dynacall

import (
	"fmt"
	"reflect"
)

// CreateDynamicStruct tạo một struct động từ danh sách reflect.Type
//
// # Tạo một struct động từ danh sách reflect.Type
//
// Create a dynamic struct from a list of reflect.Type.
//
// 创建一个动态结构，从reflect.Type列表中创建。
//
// reflect.Typeのリストから動的構造体を作成する\n
//
// return structType: dynamic struct type, structInstance: dynamic struct instance, fields: fields of dynamic struct, error
func CreateDynamicStruct(fieldTypes []reflect.Type) (
	reflect.Type, // dynamic struct type
	reflect.Value, // dynamic struct instance
	[]reflect.StructField, // fields of dynamic struct
	error,
) {
	// Tạo slice của reflect.StructField
	// Create a slice of reflect.StructField
	//创建一个reflect.StructField切片
	//reflect.StructFieldのスライスを作成する
	fields := make([]reflect.StructField, len(fieldTypes))
	for i, fieldType := range fieldTypes {
		// Đặt tên trường là Field1, Field2, ...
		fieldName := fmt.Sprintf("Field%d", i+1)
		fields[i] = reflect.StructField{
			Name: fieldName,
			Type: fieldType,
		}
	}

	// Tạo kiểu struct động
	structType := reflect.StructOf(fields)

	// Tạo instance của struct
	structInstance := reflect.New(structType).Elem()

	return structType, structInstance, fields, nil
}
