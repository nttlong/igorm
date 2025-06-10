package dbx

import (
	"fmt"
	"strings"

	"github.com/xwb1989/sqlparser"
)

func (w Compiler) walkOnFuncExpr(expr *sqlparser.FuncExpr, ctx *ParseContext) (string, error) {

	// params := []string{}
	args := []Node{}
	for _, p := range expr.Exprs {
		s, err := w.walkSQLNode(p, ctx)
		if err != nil {
			return "", err
		}
		args = append(args, Node{Nt: FunctionArg, V: s})
	}
	funcName := expr.Name.String()
	if strings.ToLower(funcName) == "row_number" {
		if selectStm, ok := ctx.Original.(*sqlparser.Select); ok {
			if selectStm.OrderBy != nil {
				strOrderBy, err := w.walkOnOrderBy(&selectStm.OrderBy, ctx)
				selectStm.OrderBy = nil
				if err != nil {
					return "", err
				} else {
					return "ROW_NUMBER() OVER (ORDER BY " + strOrderBy + ")", nil
				}
			}
		}
		return "", fmt.Errorf("row_number require order by")

	}

	n, err := w.OnParse(Node{Nt: Function, V: funcName, C: args, IsResolved: false, ctx: ctx})
	if err != nil {
		return "", err
	}
	if n.IsResolved {
		return n.V, nil

	}
	Params := []string{}
	for _, p := range n.C {
		Params = append(Params, p.V)
	}
	return n.V + "(" + strings.Join(Params, ", ") + ")", nil
}
