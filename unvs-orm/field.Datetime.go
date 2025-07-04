package orm

import "time"

type DateTimeField struct {
	*dbField
	callMethod *methodCall
	Left       interface{}
	Right      interface{}
	Op         string
	Val        *time.Time
}

func (f *DateTimeField) compare(other interface{}, op string) *BoolField {

	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      op,
	}
}
func (f *DateTimeField) Eq(other interface{}) *BoolField {
	return f.compare(other, "=")
}
func (f *DateTimeField) Ne(other interface{}) *BoolField {
	return f.compare(other, "!=")
}
func (f *DateTimeField) Gt(other interface{}) *BoolField {
	return f.compare(other, ">")
}
func (f *DateTimeField) Ge(other interface{}) *BoolField {
	return f.compare(other, ">=")
}
func (f *DateTimeField) Lt(other interface{}) *BoolField {
	return f.compare(other, "<")
}
func (f *DateTimeField) Le(other interface{}) *BoolField {
	return f.compare(other, "<=")
}
func (f *DateTimeField) Between(min, max interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   []interface{}{min, max},
		op:      "BETWEEN",
	}
}
func (f *DateTimeField) In(others ...interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   others,
		op:      "IN",
	}
}
func (f *DateTimeField) NotIn(others ...interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   others,
		op:      "NOT IN",
	}
}
func (f *DateTimeField) IsNull() *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		op:      "IS NULL",
	}
}
func (f *DateTimeField) IsNotNull() *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		op:      "IS NOT NULL",
	}
}
func (f *DateTimeField) NotBetween(min, max interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   []interface{}{min, max},
		op:      "NOT BETWEEN",
	}
}
func (f *DateTimeField) Day() *NumberField[int] {
	return &NumberField[int]{

		callMethod: &methodCall{
			method:  "DAY",
			dbField: f.dbField.clone(),
			args:    []interface{}{},
		},
	}
}
func (f *DateTimeField) Month() *NumberField[int] {
	return &NumberField[int]{

		callMethod: &methodCall{
			method:  "MONTH",
			dbField: f.dbField.clone(),
			args:    []interface{}{},
		},
	}
}
func (f *DateTimeField) Year() *NumberField[int] {
	return &NumberField[int]{

		callMethod: &methodCall{
			method:  "YEAR",
			dbField: f.dbField.clone(),
			args:    []interface{}{},
		},
	}
}
func (f *DateTimeField) Hour() *NumberField[int] {
	return &NumberField[int]{

		callMethod: &methodCall{
			method:  "HOUR",
			dbField: f.dbField.clone(),
			args:    []interface{}{},
		},
	}
}
func (f *DateTimeField) Minute() *NumberField[int] {
	return &NumberField[int]{

		callMethod: &methodCall{
			method:  "MINUTE",
			dbField: f.dbField.clone(),
			args:    []interface{}{},
		},
	}
}
func (f *DateTimeField) Second() *NumberField[int] {
	return &NumberField[int]{

		callMethod: &methodCall{
			method:  "SECOND",
			dbField: f.dbField.clone(),
			args:    []interface{}{},
		},
	}
}
func (f *DateTimeField) Format(layout string) *TextField {
	return &TextField{
		callMethod: &methodCall{
			method:  "FORMAT",
			dbField: f.dbField.clone(),
			args:    []interface{}{layout},
		},
	}
}

/*
DateTimeField_test.go:46: method Min not found in type *orm.DateTimeField
    DateTimeField_test.go:46: method Max not found in type *orm.DateTimeField
    DateTimeField_test.go:46: method Count not found in type *orm.DateTimeField
*/
func (f *DateTimeField) Min() *DateTimeField {
	return &DateTimeField{
		callMethod: &methodCall{
			method:  "MIN",
			dbField: f.dbField.clone(),
			args:    []interface{}{},
		},
	}
}
func (f *DateTimeField) Max() *DateTimeField {
	return &DateTimeField{
		callMethod: &methodCall{
			method:  "MAX",
			dbField: f.dbField.clone(),
			args:    []interface{}{},
		},
	}
}
func (f *DateTimeField) Count() *NumberField[int] {
	return &NumberField[int]{
		callMethod: &methodCall{
			method:  "COUNT",
			dbField: f.dbField.clone(),
			args:    []interface{}{},
		},
	}
}
