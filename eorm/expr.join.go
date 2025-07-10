package eorm

import (
	"github.com/xwb1989/sqlparser"
)

type exprCompileContext struct {
	tables      []string
	schema      map[string]bool
	alias       map[string]string
	dialect     Dialect
	IsBuildJoin bool
}
type exprJoin struct {
	context *exprCompileContext
	content string
}

func (e *exprJoin) build(joinText string) error {
	joinText = utils.EXPR.QuoteExpression(joinText)

	sqlTest := "select * from " + joinText
	stm, err := sqlparser.Parse(sqlTest)
	if err != nil {
		return err
	}
	if sqlSelect, ok := stm.(*sqlparser.Select); ok {

		for _, expr := range sqlSelect.From {
			strResult, err := exprs.compile(e.context, expr)
			if err != nil {
				return err
			}
			e.content = strResult
		}
	}

	return nil

}
