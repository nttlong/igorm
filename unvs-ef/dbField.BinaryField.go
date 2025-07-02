package unvsef

// func (f *BinaryField) And(other BinaryField) *FieldBool {
// 	return &FieldBool{
// 		DbField: f.DbField.clone(),
// 		BinaryField: &BinaryField{

// 			Left: f, Op: "AND",
// 			Right: other,
// 		},
// 	}
// }
func (f *BinaryField) Or(other BinaryField) *BinaryField {

	return &BinaryField{
		Left: f, Op: "OR",
		Right: other,
	}
}
