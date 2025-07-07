package orm

type BoolField1 struct {
	JoinExpr *JoinExpr
	expr     string
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
type rawTextField struct {
	rawText string
}
type BoolField struct {
	UnderField interface{}
}

func (f *BoolField) Raw(text string) *BoolField {
	return &BoolField{
		UnderField: &rawTextField{
			rawText: text,
		},
	}

}
func (f *BoolField) makeFieldBinary(other interface{}, op string) *fieldBinary {
	return &fieldBinary{
		left:  f,
		right: other,
		op:    op,
	}
}
func (f *BoolField) Eq(value interface{}) *fieldBinary {
	return f.makeFieldBinary(value, "=")
}
func (f *BoolField) And(other interface{}) *BoolField {
	return &BoolField{
		UnderField: f.makeFieldBinary(other, "AND"),
	}

}
func (f *BoolField) Or(other interface{}) *BoolField {
	return &BoolField{
		UnderField: f.makeFieldBinary(other, "OR"),
	}
}
func (f *BoolField) Not() *BoolField {
	return &BoolField{
		UnderField: fieldBinary{
			left:  nil,
			right: f,
			op:    "NOT",
		},
	}
}
