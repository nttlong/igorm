package unvsef

func (f *Field[TField]) Mod(other interface{}) *Field[TField] {
	return &Field[TField]{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "%",
			Right: other,
		},
	}
}

func (f *Field[TField]) In(other interface{}) *Field[TField] {
	return &Field[TField]{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "IN",
			Right: other,
		},
	}
}
func (f *Field[TField]) NotIn(other interface{}) *Field[TField] {
	return &Field[TField]{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "NOT IN",
			Right: other,
		},
	}
}
func (f *Field[TField]) Between(left interface{}, right interface{}) *Field[TField] {
	return &Field[TField]{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{
			Left:  f,
			Op:    "BETWEEN",
			Right: []interface{}{left, right},
		},
	}
}
func (f *Field[TField]) As(alias string) *AliasField {
	return &AliasField{
		Field: f,
		Alias: alias,
	}
}
