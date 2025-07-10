package orm

// type TextField struct {
// 	*dbField
// 	callMethod *methodCall
// 	Val        *string
// }
type TextField struct {
	underField underFieldObject
	Val        *string
}

func (f *TextField) As(name string) *aliasField {
	return &aliasField{
		underField: f,
		Alias:      name,
	}
}
func (f *TextField) makFieldBinary(other interface{}, op string) *fieldBinary {
	return &fieldBinary{
		left:  f,
		right: other,
		op:    op,
	}
}
func (f *TextField) Eq(value interface{}) *BoolField {
	return &BoolField{
		underField: f.makFieldBinary(value, "="),
	}
}
func (f *TextField) Ne(value interface{}) *BoolField {
	return &BoolField{
		underField: f.makFieldBinary(value, "!="),
	}
}
func (f *TextField) Like(value interface{}) *BoolField {
	return &BoolField{
		underField: f.makFieldBinary(value, "LIKE"),
	}
}
func (f *TextField) NotLike(value interface{}) *BoolField {
	return &BoolField{
		underField: f.makFieldBinary(value, "NOT LIKE"),
	}
}
func (f *TextField) In(values interface{}) *BoolField {
	return &BoolField{
		underField: f.makFieldBinary(values, "IN"),
	}
}
func (f *TextField) NotIn(values interface{}) *BoolField {
	return &BoolField{
		underField: f.makFieldBinary(values, "NOT IN"),
	}
}
func (f *TextField) IsNull() *BoolField {
	return &BoolField{
		underField: f.makFieldBinary(nil, "IS NULL"),
	}
}
func (f *TextField) IsNotNull() *BoolField {
	return &BoolField{
		underField: f.makFieldBinary(nil, "IS NOT NULL"),
	}
}
func (f *TextField) Between(start, end interface{}) *BoolField {
	return &BoolField{
		underField: f.makFieldBinary([]interface{}{start, end}, "BETWEEN"),
	}
}

func (f *TextField) NotBetween(start, end interface{}) *BoolField {
	return &BoolField{
		underField: f.makFieldBinary([]interface{}{start, end}, "NOT BETWEEN"),
	}
}

//---------------------------

func (f *TextField) Len() *NumberField[int] {
	return &NumberField[int]{
		underField: &methodCall{
			method: "LEN",
			args:   []interface{}{f},
		},
	}
}

func (f *TextField) Upper() *TextField {

	return &TextField{
		underField: &methodCall{
			method: "UPPER",
			args:   []interface{}{f},
		},
	}
}
func (f *TextField) Lower() *TextField {
	return &TextField{
		underField: &methodCall{
			method: "LOWER",
			args:   []interface{}{f},
		},
	}
}
func (f *TextField) Trim() *TextField {
	return &TextField{
		underField: &methodCall{
			method: "TRIM",
			args:   []interface{}{f},
		},
	}
}
func (f *TextField) LTrim() *TextField {
	return &TextField{
		underField: &methodCall{
			method: "LTRIM",
			args:   []interface{}{f},
		},
	}
}
func (f *TextField) RTrim() *TextField {
	return &TextField{
		underField: &methodCall{
			method: "RTRIM",
			args:   []interface{}{f},
		},
	}
}
func (f *TextField) Concat(args ...interface{}) *TextField {
	return &TextField{
		underField: &methodCall{
			method: "CONCAT",
			args:   append([]interface{}{f}, args...),
		},
	}
}
func (f *TextField) Replace(old, new interface{}) *TextField {
	return &TextField{
		underField: &methodCall{
			method: "REPLACE",
			args:   []interface{}{f, old, new},
		},
	}
}
func (f *TextField) Substr(start, length interface{}) *TextField {
	return &TextField{
		underField: &methodCall{
			method: "SUBSTR",
			args:   []interface{}{f, start, length},
		},
	}

}
func (f *TextField) Left(length interface{}) *TextField {
	return &TextField{
		underField: &methodCall{
			method: "LEFT",
			args:   []interface{}{f, length},
		},
	}

}
func (f *TextField) Right(length interface{}) *TextField {
	return &TextField{
		underField: &methodCall{
			method: "RIGHT",
			args:   []interface{}{f, length},
		},
	}
}
