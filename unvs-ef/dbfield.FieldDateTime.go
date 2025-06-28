package unvsef

import "time"

type FieldDateTime Field[time.Time]

func (f *FieldDateTime) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
func (f *FieldDateTime) Year() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "YEAR",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Month() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "MONTH",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Day() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "DAY",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Hour() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "HOUR",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Minute() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "MINUTE",
			Args:     []interface{}{f},
		},
	}
}
func (f *FieldDateTime) Second() *Field[int] {
	return &Field[int]{
		FuncField: &FuncField{
			FuncName: "SECOND",
			Args:     []interface{}{f},
		},
	}
}
func (f *BinaryField) ToSqlExpr(d Dialect) (string, []interface{}) {
	return compiler.Compile(f, d)
}
