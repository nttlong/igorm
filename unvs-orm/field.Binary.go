package orm

// type fieldBinary2 struct {
// 	*dbField
// 	left  interface{}
// 	right interface{}
// 	op    string
// }
type fieldBinary struct {
	left  interface{}
	right interface{}
	op    string
}

func (f *fieldBinary) As(Name string) *aliasField {
	return &aliasField{
		UnderField: f,
		Alias:      Name,
	}
}
func (f *fieldBinary) makeFieldBinary(right interface{}, op string) *fieldBinary {
	return &fieldBinary{
		left:  f,
		right: right,
		op:    op,
	}
}

func (f *fieldBinary) Eq(other interface{}) *BoolField {
	return &BoolField{
		UnderField: f.makeFieldBinary(other, "="),
	}
}
func (f *fieldBinary) Ne(other interface{}) *BoolField {
	return &BoolField{
		UnderField: f.makeFieldBinary(other, "!="),
	}
}
func (f *fieldBinary) Gt(other interface{}) *BoolField {
	return &BoolField{
		UnderField: f.makeFieldBinary(other, ">"),
	}
}
func (f *fieldBinary) Lt(other interface{}) *BoolField {
	return &BoolField{
		UnderField: f.makeFieldBinary(other, "<"),
	}
}
func (f *fieldBinary) Ge(other interface{}) *BoolField {
	return &BoolField{
		UnderField: f.makeFieldBinary(other, ">="),
	}
}

func (f *fieldBinary) Le(other interface{}) *BoolField {
	return &BoolField{
		UnderField: f.makeFieldBinary(other, "<="),
	}
}
func (f *fieldBinary) In(others interface{}) *BoolField {
	return &BoolField{
		UnderField: f.makeFieldBinary(others, "IN"),
	}
}
func (f *fieldBinary) NotIn(others interface{}) *BoolField {
	return &BoolField{
		UnderField: f.makeFieldBinary(others, "NOT IN"),
	}
}

// ==============================

func (f *fieldBinary) And(other interface{}) *BoolField {
	return &BoolField{
		UnderField: f.makeFieldBinary(other, "AND"),
	}
}
func (f *fieldBinary) Or(other interface{}) *BoolField {
	return &BoolField{
		UnderField: f.makeFieldBinary(other, "OR"),
	}
}
