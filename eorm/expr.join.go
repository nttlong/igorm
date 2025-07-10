package eorm

import (
	"fmt"
	"strings"

	"github.com/xwb1989/sqlparser"
)

type build_purpose int

const (
	build_purpose_select build_purpose = iota
	build_purpose_join
	build_purpose_where
	build_purpose_group
	build_purpose_having
	build_purpose_order
	build_purpose_limit
	build_purpose_offset
	build_purpose_for_function
)

type exprCompileContext struct {
	tables           []string
	schema           map[string]bool
	alias            map[string]string
	dialect          Dialect
	purpose          build_purpose
	stackAliasFields stack[string]
}
type exprCompiler struct {
	context *exprCompileContext
	content string
}

func (e *exprCompiler) buildSelectField(selector string) error {
	e.context.purpose = build_purpose_select
	selector = utils.EXPR.QuoteExpression(selector)
	sqlTest := "select " + selector
	stm, err := sqlparser.Parse(sqlTest)
	if err != nil {
		return err
	}
	if sqlSelect, ok := stm.(*sqlparser.Select); ok {
		selectors := make([]string, len(sqlSelect.SelectExprs))
		for i, expr := range sqlSelect.SelectExprs {
			if sqlExpr, ok := expr.(*sqlparser.AliasedExpr); ok {
				if !sqlExpr.As.IsEmpty() {
					e.context.stackAliasFields.Push(sqlExpr.As.String())
				}
				if sqlExpr.Expr != nil {
					strResult, err := exprs.compile(e.context, sqlExpr.Expr)

					if err != nil {
						return err
					}
					selectors[i] = strResult

				}
			} else {
				panic(fmt.Errorf("unsupported select type is %T", expr))
			}
		}
		e.content = strings.Join(selectors, ", ")
	}

	return nil
}
func (e *exprCompiler) build(joinText string) error {
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
