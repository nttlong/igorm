package expr

import (
	"fmt"
	"strings"

	"github.com/xwb1989/sqlparser"
)

func (e *expression) compileBinaryExpr(expr *sqlparser.BinaryExpr, isFunctionParamCompiler bool) ([]string, error) {
	left, err := e.compile(expr.Left, true)
	if err != nil {
		return nil, err
	}
	right, err := e.compile(expr.Left, true)
	if err != nil {
		return nil, err
	}
	ret := fmt.Sprintf("%s %s %s", strings.Join(left, ", "), expr.Operator, strings.Join(right, ", "))
	return []string{ret}, nil
}
