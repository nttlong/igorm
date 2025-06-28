package unvsef

// type RawExpr struct {
// 	sql  string
// 	args []interface{}
// 	ft   DbFieldType
// }

// func (r RawExpr) ToSQL(d Dialect) (string, []interface{}) {
// 	return r.sql, r.args
// }
// func (r RawExpr) SetType(ft DbFieldType) {
// 	r.ft = ft
// }

// // Hàm tiện ích để tạo RawExpr
// func Raw(sql string, args ...interface{}) Expr {
// 	return RawExpr{sql: sql, args: args, ft: DbFieldTypeConst}
// }
// func AddExpr(left Expr, op string, right Expr) *BinaryExpr {
// 	return &BinaryExpr{Left: left, Op: op, Right: right, ft: DbFieldTypeExpr}
// }
