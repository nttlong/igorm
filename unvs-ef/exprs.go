package unvsef

func AddExpr(left Expr, op string, right Expr) *BinaryExpr {
	return &BinaryExpr{Left: left, Op: op, Right: right}
}
