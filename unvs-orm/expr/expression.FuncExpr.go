package expr

import (
	"github.com/xwb1989/sqlparser"
)

type MethodCall struct {
	Method string
	Args   []interface{}
}

func (e *expression) funcExpr(funcExpr *sqlparser.FuncExpr, context *ResolveContext) (*expressionCompileResult, error) {

	args := make([]interface{}, len(funcExpr.Exprs))
	for i, expr := range funcExpr.Exprs {
		retArs, err := e.compile(expr, context, true)

		if err != nil {
			return nil, err
		}
		args[i] = retArs.Syntax

	}

	r, err := e.resolve(nil, &MethodCall{
		Method: funcExpr.Name.CompliantName(),
		Args:   args,
	})
	if err != nil {
		return nil, err
	}
	return &expressionCompileResult{Syntax: r.Syntax}, nil

}
