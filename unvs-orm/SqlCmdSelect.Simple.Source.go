package orm

import (
	"errors"
	"fmt"
	"strings"
)

func (sql *SqlCmdSelect) simpleSource(source string, d DialectCompiler) *SqlStmt {
	ret := &SqlStmt{
		From: sql.cmp.Quote(source),
	}
	ret.Select = sql.buildSelect2(source)
	return ret
}
func (sql *SqlCmdSelect) buildSelect2(sourceCache string) string {
	fields := []string{}
	for _, field := range sql.fields {
		if baseExpr, ok := field.(*exprField); ok {
			if sql.tables == nil {
				sql.tables = &[]string{}
			}
			if sql.buildContext == nil {
				sql.buildContext = &map[string]string{}
			}

			compileResult, err := sql.exprCmp.Compile(sourceCache, sql.tables, sql.buildContext, baseExpr.Stmt, false, true)
			if err != nil {
				sql.Err = err
			}
			fields = append(fields, compileResult.Syntax)

		} else {
			panic(errors.New(fmt.Sprintf("unsupported field type %T", field)))
		}
	}
	return strings.Join(fields, ", ")
}
