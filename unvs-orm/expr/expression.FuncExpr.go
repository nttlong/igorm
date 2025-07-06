package expr

import (
	"github.com/xwb1989/sqlparser"
)

type MethodCall struct {
	Method string
	Args   []interface{}
}

func (e *expression) funcExpr(expr *sqlparser.FuncExpr) ([]string, error) {

	retArs, err := e.compile(expr.Exprs, true)

	args := make([]interface{}, len(retArs))
	for i, arg := range retArs {
		args[i] = arg
	}
	r, err := e.resolve(nil, &MethodCall{
		Method: expr.Name.CompliantName(),
		Args:   args,
	})
	if err != nil {
		return nil, err
	}
	return []string{r.Syntax}, nil

}
