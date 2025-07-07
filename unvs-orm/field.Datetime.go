package orm

import "time"

// type DateTimeField struct {
// 	*dbField
// 	callMethod *methodCall
// 	Left       interface{}
// 	Right      interface{}
// 	Op         string
// 	Val        *time.Time
// }
type DateTimeField struct {
	UnderField interface{}
	Val        *time.Time
}

func (f *DateTimeField) As(alias string) *aliasField {
	return &aliasField{
		UnderField: f,
		Alias:      alias,
	}
}
func (f *DateTimeField) makeFieldBinary(other interface{}, op string) *fieldBinary {
	return &fieldBinary{
		left:  f,
		right: other,
		op:    op,
	}
}
func (f *DateTimeField) compare(other interface{}, op string) *BoolField {

	return &BoolField{
		UnderField: f.makeFieldBinary(other, op),
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
	return f.compare(interface{}([]interface{}{min, max}), "BETWEEN")
}
func (f *DateTimeField) In(others ...interface{}) *BoolField {
	return f.compare(others, "IN")
}
func (f *DateTimeField) NotIn(others ...interface{}) *BoolField {
	return f.compare(others, "NOT IN")
}
func (f *DateTimeField) IsNull() *BoolField {
	return f.compare(nil, "IS NULL")
}
func (f *DateTimeField) IsNotNull() *BoolField {
	return f.compare(nil, "IS NOT NULL")

}
func (f *DateTimeField) NotBetween(min, max interface{}) *BoolField {
	return f.compare(interface{}([]interface{}{min, max}), "NOT BETWEEN")

}
func (f *DateTimeField) makeDateTimeMethodCall(method string) *NumberField[int] {
	return &NumberField[int]{
		UnderField: &methodCall{
			method: method,
			args:   []interface{}{f},
		},
	}
}
func (f *DateTimeField) Day() *NumberField[int] {
	return f.makeDateTimeMethodCall("DAY")
}
func (f *DateTimeField) Month() *NumberField[int] {
	return f.makeDateTimeMethodCall("MONTH")
}
func (f *DateTimeField) Year() *NumberField[int] {

	return f.makeDateTimeMethodCall("YEAR")
}
func (f *DateTimeField) Hour() *NumberField[int] {
	return f.makeDateTimeMethodCall("HOUR")
}
func (f *DateTimeField) Minute() *NumberField[int] {
	return f.makeDateTimeMethodCall("MINUTE")
}
func (f *DateTimeField) Second() *NumberField[int] {
	return f.makeDateTimeMethodCall("SECOND")
}
func (f *DateTimeField) Format(layout string) *TextField {
	return &TextField{
		UnderField: &methodCall{
			method: "FORMAT",
			args:   []interface{}{f, layout},
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
		UnderField: &methodCall{
			method: "MIN",
			args:   []interface{}{f},
		},
	}
}
func (f *DateTimeField) Max() *DateTimeField {
	return &DateTimeField{
		UnderField: &methodCall{
			method: "MAX",
			args:   []interface{}{f},
		},
	}
}
func (f *DateTimeField) Count() *NumberField[int] {
	return f.makeDateTimeMethodCall("COUNT")
}
