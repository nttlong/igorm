package unvsef

import (
	"time"
)

type FieldDateTime Field[time.Time]

func (f *FieldDateTime) Set(val *time.Time) {
	f.val = val
}

func (f *FieldDateTime) Year() *FieldNumber[int] {
	return &FieldNumber[int]{
		DbField: f.DbField.clone(),
		FuncField: &FuncField{

			FuncName: "YEAR",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Month() *Field[int] {
	return &Field[int]{
		DbField: f.DbField.clone(),
		FuncField: &FuncField{

			FuncName: "MONTH",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Day() *Field[int] {
	return &Field[int]{
		DbField: f.DbField.clone(),
		FuncField: &FuncField{
			FuncName: "DAY",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Hour() *Field[int] {
	return &Field[int]{
		DbField: f.DbField.clone(),
		FuncField: &FuncField{
			FuncName: "HOUR",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Minute() *Field[int] {
	return &Field[int]{
		DbField: f.DbField.clone(),
		FuncField: &FuncField{
			FuncName: "MINUTE",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Second() *Field[int] {
	return &Field[int]{
		DbField: f.DbField.clone(),
		FuncField: &FuncField{
			FuncName: "SECOND",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Between(left interface{}, right interface{}) *FieldBool {
	return &FieldBool{
		DbField: f.DbField.clone(),
		BinaryField: &BinaryField{

			Left:  f,
			Op:    "BETWEEN",
			Right: []interface{}{left, right},
		},
	}
}
func (f *FieldDateTime) As(alias string) *AliasField {
	return &AliasField{
		Field: f,
		Alias: alias,
	}
}
