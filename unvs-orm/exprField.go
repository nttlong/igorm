package orm

func (f *exprField) makeBinaryExpr(left interface{}, right interface{}, op string) *exprField {
	return &exprField{
		UnderField: fieldBinary{
			left:  left,
			right: right,
			op:    op,
		},
	}
}
func (f *exprField) Add(order interface{}) *exprField {
	return f.makeBinaryExpr(f, order, "+")
}

// --------------------
func (f *exprField) makeMethodCallExpr(method string, args []interface{}) *exprField {
	return &exprField{
		UnderField: methodCall{

			method: method,
			args:   args,
		},
	}
}
func (f *exprField) Year() *exprField {
	return f.makeMethodCallExpr("Year", []interface{}{f})
}
