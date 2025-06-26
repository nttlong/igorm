package unvsef

type RawExpr struct {
	sql  string
	args []interface{}
}

func (r RawExpr) ToSQL(d Dialect) (string, []interface{}) {
	return r.sql, r.args
}

// Hàm tiện ích để tạo RawExpr
func Raw(sql string, args ...interface{}) Expr {
	return RawExpr{sql: sql, args: args}
}
func AddExpr(left Expr, op string, right Expr) *BinaryExpr {
	return &BinaryExpr{Left: left, Op: op, Right: right}
}
