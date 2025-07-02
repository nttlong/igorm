package orm

type fieldBinary struct {
	*dbField
	left  interface{}
	right interface{}
	op    string
}

func (f *fieldBinary) As() *aliasField {
	return &aliasField{
		Expr:  &f.dbField,
		Alias: f.Name,
	}
}

func (f *fieldBinary) Eq(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "=",
	}
}
func (f *fieldBinary) Ne(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "!=",
	}
}
func (f *fieldBinary) Gt(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      ">",
	}
}
func (f *fieldBinary) Lt(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "<",
	}
}
func (f *fieldBinary) Ge(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      ">=",
	}
}
func (f *fieldBinary) Le(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "<=",
	}
}
func (f *fieldBinary) In(others interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   others,
		op:      "IN",
	}
}
func (f *fieldBinary) NotIn(others interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   others,
		op:      "NOT IN",
	}
}

// ==============================

func (f *fieldBinary) And(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "AND",
	}
}
func (f *fieldBinary) Or(other interface{}) *BoolField {
	return &BoolField{
		dbField: f.dbField.clone(),
		left:    f,
		right:   other,
		op:      "OR",
	}
}
