package unvsef

func (f *Field[TField]) Mod(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "%",
			Right: other,
		},
	}
}
func (f *Field[TField]) Like(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "LIKE",
			Right: other,
		},
	}
}
func (f *Field[TField]) NotLike(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "LIKE",
			Right: other,
		},
	}
}
func (f *Field[TField]) In(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "IN",
			Right: other,
		},
	}
}
func (f *Field[TField]) NotIn(other interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "NOT IN",
			Right: other,
		},
	}
}
func (f *Field[TField]) Between(left interface{}, right interface{}) *Field[TField] {
	return &Field[TField]{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "BETWEEN",
			Right: []interface{}{left, right},
		},
	}
}
func (f *FieldDateTime) Between(left interface{}, right interface{}) *FieldBool {
	return &FieldBool{
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "BETWEEN",
			Right: []interface{}{left, right},
		},
	}
}
