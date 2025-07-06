package orm

import (
	"github.com/xwb1989/sqlparser"
)

func (e *expression) funcExpr(expr *sqlparser.FuncExpr) ([]string, error) {

	retArs, err := e.compile(expr.Exprs, true)

	args := make([]interface{}, len(retArs))
	for i, arg := range retArs {
		args[i] = arg
	}
	r, err := e.dialect.resolve(nil, &methodCall{
		method: expr.Name.CompliantName(),
		args:   args,
	})
	if err != nil {
		return nil, err
	}
	return []string{r.Syntax}, nil

}
