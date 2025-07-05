package orm

type BoolField struct {
	*dbField
	left            interface{}
	right           interface{}
	op              string
	val             *bool
	rawText         string
	joinType        string
	alias           map[string]string
	tables          []string
	joinSource      string
	joinSourceAlias string
}

func (f *BoolField) Raw(text string) *BoolField {
	f.rawText = text
	return f
}
func (f *BoolField) Eq(value interface{}) *fieldBinary {
	return &fieldBinary{

		dbField: f.dbField.clone(),
		left:    f,
		right:   value,
		op:      "=",
	}
}
func (f *BoolField) And(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "AND",
	}

}
func (f *BoolField) Or(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "OR",
	}
}
func (f *BoolField) Not() *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		op:      "NOT",
	}
}
func (f *BoolField) Get() *bool {
	return f.val
}
func (f *BoolField) Set(val *bool) {
	f.val = val
}
