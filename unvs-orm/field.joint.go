package orm

func (f *NumberField[T]) Join(other interface{}) *BoolField {

	return &BoolField{
		dbField:  f.dbField.clone(),
		left:     f,
		right:    other,
		op:       "=",
		joinType: "INNER",
	}
}
func (f *NumberField[T]) LeftJoin(other interface{}) *BoolField {

	return &BoolField{
		dbField:  f.dbField.clone(),
		left:     f,
		right:    other,
		op:       "=",
		joinType: "LEFT",
	}
}
func (f *NumberField[T]) RightJoin(other interface{}) *BoolField {

	return &BoolField{
		dbField:  f.dbField.clone(),
		left:     f,
		right:    other,
		op:       "=",
		joinType: "RIGHT",
	}
}
func (f *TextField) Join(other interface{}) *BoolField {

	return &BoolField{
		dbField:  f.dbField.clone(),
		left:     f,
		right:    other,
		op:       "=",
		joinType: "INNER",
	}
}
func (f *TextField) LeftJoin(other interface{}) *BoolField {

	return &BoolField{
		dbField:  f.dbField.clone(),
		left:     f,
		right:    other,
		op:       "=",
		joinType: "LEFT",
	}
}
func (f *TextField) RightJoin(other interface{}) *BoolField {

	return &BoolField{
		dbField:  f.dbField.clone(),
		left:     f,
		right:    other,
		op:       "=",
		joinType: "RIGHT",
	}
}
func (f *DateTimeField) Join(other interface{}) *BoolField {

	return &BoolField{
		dbField:  f.dbField.clone(),
		left:     f,
		right:    other,
		op:       "=",
		joinType: "INNER",
	}
}
func (f *DateTimeField) LeftJoin(other interface{}) *BoolField {

	return &BoolField{
		dbField:  f.dbField.clone(),
		left:     f,
		right:    other,
		op:       "=",
		joinType: "LEFT",
	}
}
func (f *DateTimeField) RightJoin(other interface{}) *BoolField {

	return &BoolField{
		dbField:  f.dbField.clone(),
		left:     f,
		right:    other,
		op:       "=",
		joinType: "RIGHT",
	}
}

// fieldBinary
func (f *fieldBinary) Join(other interface{}) *BoolField {

	return &BoolField{
		dbField:  f.dbField.clone(),
		left:     f,
		right:    other,
		op:       "=",
		joinType: "INNER",
	}
}
func (f *fieldBinary) LeftJoin(other interface{}) *BoolField {

	return &BoolField{
		dbField:  f.dbField.clone(),
		left:     f,
		right:    other,
		op:       "=",
		joinType: "LEFT",
	}
}
func (f *fieldBinary) RightJoin(other interface{}) *BoolField {

	return &BoolField{
		dbField:  f.dbField.clone(),
		left:     f,
		right:    other,
		op:       "=",
		joinType: "RIGHT",
	}
}
