package uorm

func (f *Field) Sub(other Field) *Field {
	return makeBinaryExpr(*f, other, "-")
}

func (f *Field) Mul(other Field) *Field {
	return makeBinaryExpr(*f, other, "*")
}

func (f *Field) Div(other Field) *Field {
	return makeBinaryExpr(*f, other, "/")
}

func (f *Field) Mod(other Field) *Field {
	return makeBinaryExpr(*f, other, "%")
}

func (f *Field) Add(other any) *Field {
	return makeBinaryExpr(*f, other, "+")
}
