package expr

import (
	"github.com/xwb1989/sqlparser"
)

type MethodCall struct {
	Method string
	Args   []interface{}
}

func (e *expression) funcExpr(funcExpr *sqlparser.FuncExpr, tables *[]string, context *map[string]string, requireAlias bool) (*expressionCompileResult, error) {

	args := make([]interface{}, len(funcExpr.Exprs))
	for i, expr := range funcExpr.Exprs {
		retArs, err := e.compile(expr, tables, context, true, requireAlias)

		if err != nil {
			return nil, err
		}
		args[i] = retArs.Syntax

	}

	r, err := e.resolve(tables, context, &MethodCall{
		Method: funcExpr.Name.CompliantName(),
		Args:   args,
	}, requireAlias)
	if err != nil {
		return nil, err
	}
	return &expressionCompileResult{Syntax: r.Syntax}, nil

}
