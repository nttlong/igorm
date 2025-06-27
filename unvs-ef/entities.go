package unvsef

import (
	"reflect"
)

// Entity creates an instance of struct T and auto-populates any fields
// named "TableName" and "ColName" with the corresponding struct and field names.
//
// This is useful for initializing DbField[TTable, TField] fields in a model,
// so that the table and column names can be inferred via reflection without manual assignment.
//
// Requirements:
// - Each field inside struct T must be a struct that contains fields named "TableName" and "ColName".
// - Those inner fields must be settable (exported and addressable).
func Entity[T any]() T {
	var v T

	// Get the type name of T to use as table name
	typ := reflect.TypeOf(v)

	ret := EntityFromType(typ)
	return ret.Interface().(T)
}
func EntityFromType(typ reflect.Type) reflect.Value {
	val := reflect.New(typ).Elem()

	for i := 0; i < typ.NumField(); i++ {
		// Locate and set the "TableName" field inside each struct field
		tableNameField := val.Field(i).FieldByName("TableName")
		if tableNameField.IsValid() && tableNameField.CanSet() {
			tableNameField.SetString(utils.TableNameFromStruct(typ))
		}

		// Locate and set the "ColName" field inside each struct field
		columnNameField := val.Field(i).FieldByName("ColName")
		if columnNameField.IsValid() && columnNameField.CanSet() {
			columnNameField.SetString(utils.ToSnakeCase(typ.Field(i).Name))
		}
	}
	return val
}
