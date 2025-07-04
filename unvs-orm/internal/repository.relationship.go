package internal

import (
	"fmt"
	"reflect"
)

type RelationshipRegister struct {
	owner      *Base
	err        error
	fromTable  string
	fromFields []string
	toTable    string
	toFields   []string
}

func (r *RelationshipRegister) From(fields ...interface{}) *RelationshipRegister {
	for _, field := range fields {
		valField := reflect.ValueOf(field)
		tableNameField := valField.FieldByName("TableName")
		tableName := ""
		fieldName := ""
		if tableNameField.IsValid() {
			tableName = tableNameField.String()
		} else {
			r.err = fmt.Errorf("TableName not found in field at RelationshipRegister.From '%s'", field)
			r.owner.Err = r.err
			return r
		}
		if r.fromTable == "" {
			r.fromTable = tableName
		}

		fieldNameField := valField.FieldByName("ColName")

		if fieldNameField.IsValid() {
			fieldName = fieldNameField.String()
		} else {
			r.err = fmt.Errorf("TableName not found in field in field at RelationshipRegister.From '%s'", field)
			return r
		}
		if r.fromTable != "" && r.fromTable != tableName {

			r.owner.Err = fmt.Errorf("%s must belong to the same table name '%s' but got '%s'", fieldName, r.fromTable, tableName)
			return r
		}

		r.fromFields = append(r.fromFields, fieldName)
		r.fromTable = tableName

	}
	return r

	// panic("not implemented")
}
func (r *RelationshipRegister) To(fields ...interface{}) error {
	if r.err != nil {
		return r.err
	}
	for _, field := range fields {
		valField := reflect.ValueOf(field)
		tableNameField := valField.FieldByName("TableName")
		tableName := ""
		fieldName := ""
		if tableNameField.IsValid() {
			tableName = tableNameField.String()
		} else {
			r.owner.Err = fmt.Errorf("TableName not found in field at RelationshipRegister.To '%s'", field)
			return r.owner.Err
		}
		fieldNameField := valField.FieldByName("ColName")
		if fieldNameField.IsValid() {
			fieldName = fieldNameField.String()
		} else {
			r.owner.Err = fmt.Errorf("TableName not found in field at RelationshipRegister.To '%s'", field)
			return r.owner.Err
		}
		if r.toTable == "" {
			r.toTable = tableName
		}
		if r.toTable != tableName {
			r.err = fmt.Errorf("%s must belong to the same table name '%s' but got '%s'", fieldName, r.fromTable, tableName)
			return r.err
		}

	}
	// check the balance of num of field in from and to
	if len(r.fromFields) != len(r.toFields) {
		r.err = fmt.Errorf("number of fields were not balance, from %d to %d", len(r.fromFields), len(r.toFields))
		return r.err
	}

	return nil

}
