package orm

type TextField struct {
	*dbField
	callMethod *methodCall
	val        *string
}

func (f *TextField) As(name string) *aliasField {
	return &aliasField{
		Expr:  f.dbField,
		Alias: name,
	}
}
func (f *TextField) Eq(value interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   value,
		op:      "=",
	}
}
func (f *TextField) Ne(value interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   value,
		op:      "!=",
	}
}
func (f *TextField) Like(value interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   value,
		op:      "LIKE",
	}
}
func (f *TextField) NotLike(value interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   value,
		op:      "NOT LIKE",
	}
}
func (f *TextField) In(values interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   values,
		op:      "IN",
	}
}
func (f *TextField) NotIn(values interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   values,
		op:      "NOT IN",
	}
}
func (f *TextField) IsNull() *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		op:      "IS NULL",
	}
}
func (f *TextField) IsNotNull() *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		op:      "IS NOT NULL",
	}
}
func (f *TextField) Between(start, end interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   []interface{}{start, end},
		op:      "BETWEEN",
	}
}
func (f *TextField) NotBetween(start, end interface{}) *fieldBinary {
	return &fieldBinary{
		dbField: f.dbField.clone(),
		left:    f,
		right:   []interface{}{start, end},
		op:      "NOT BETWEEN",
	}
}
func (f *TextField) Set(val *string) {
	f.val = val
}
func (f *TextField) Get() *string {
	return f.val
}

//---------------------------

func (f *TextField) Len() *methodCall {
	return &methodCall{
		method: "LEN",
		args:   []interface{}{f},
	}
}
func (f *TextField) Upper() *methodCall {
	return &methodCall{
		method: "UPPER",
		args:   []interface{}{f},
	}
}
func (f *TextField) Lower() *methodCall {
	return &methodCall{
		method: "LOWER",
		args:   []interface{}{f},
	}
}
func (f *TextField) Trim() *methodCall {
	return &methodCall{
		method: "TRIM",
		args:   []interface{}{f},
	}
}
func (f *TextField) LTrim() *methodCall {
	return &methodCall{
		method: "LTRIM",
		args:   []interface{}{f},
	}
}
func (f *TextField) RTrim() *methodCall {
	return &methodCall{
		method: "RTRIM",
		args:   []interface{}{f},
	}
}
func (f *TextField) Concat(args ...interface{}) *methodCall {
	return &methodCall{
		method: "CONCAT",
		args:   append([]interface{}{f}, args...),
	}
}
func (f *TextField) Replace(old, new interface{}) *methodCall {
	return &methodCall{
		method: "REPLACE",
		args:   []interface{}{f, old, new},
	}
}
func (f *TextField) Substr(start, length interface{}) *methodCall {
	return &methodCall{
		method: "SUBSTR",
		args:   []interface{}{f, start, length},
	}
}
func (f *TextField) Left(length interface{}) *methodCall {
	return &methodCall{
		method: "LEFT",
		args:   []interface{}{f, length},
	}
}
func (f *TextField) Right(length interface{}) *methodCall {
	return &methodCall{
		method: "RIGHT",
		args:   []interface{}{f, length},
	}
}
