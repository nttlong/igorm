package orm

import (
	"errors"
	"fmt"
	"strings"
)

type sqlSelectSource struct {
	expr       interface{}
	sourceText string
	args       []interface{}
}
type sqlSelectFields struct {
	fields     []interface{}
	aliasMap   *map[string]string
	selectText string
	args       []interface{}
}

func (s *sqlSelectSource) build(tables *[]string, context *map[string]string, d DialectCompiler) (*resolverResult, error) {
	ctx := JoinCompiler.Ctx(d)
	if bf, ok := s.expr.(*BoolField); ok {
		return ctx.ResolveBoolFieldAsJoin(tables, context, bf)
	}
	if expr, ok := s.expr.(*JoinExpr); ok {
		return ctx.Resolve(expr)
		// return nil, errors.New("unsupported expression type: " + fmt.Sprintf("%T", expr))
	}
	return nil, errors.New("unsupported expression type: " + fmt.Sprintf("%T", s.expr))

}

type SqlCmdSelect struct {
	source       *sqlSelectSource
	fields       []interface{}
	where        interface{}
	buildContext *map[string]string
	cmp          *CompilerUtils
	Err          error
	tables       *[]string
}

func (sql *SqlCmdSelect) buildSelect() *sqlCmdSelectResult {

	fields := []string{}
	args := []interface{}{}
	for _, field := range sql.fields {
		if _, ok := field.(*aliasField); ok {
			fieldCompiler, err := sql.cmp.Resolve(sql.tables, sql.buildContext, field, true)
			if err != nil {
				return &sqlCmdSelectResult{
					Err: err,
				}
			}
			fields = append(fields, fieldCompiler.Syntax)
			args = append(args, fieldCompiler.Args...)
		} else {
			txtField, argsField, err := sql.buildSelectField(field)
			if err != nil {
				return &sqlCmdSelectResult{
					Err: err,
				}
			}
			fields = append(fields, txtField)
			args = append(args, argsField...)

		}
	}
	return &sqlCmdSelectResult{
		SqlText: strings.Join(fields, ", "),
		Args:    args,
	}
}

type sqlCmdSelectResult struct {
	Err     error
	SqlText string
	Args    []interface{}
}

func (expr *BoolField) Select(fields ...interface{}) *SqlCmdSelect {
	if expr == nil {
		return &SqlCmdSelect{
			Err: errors.New("source was not found"),
		}
	}
	return &SqlCmdSelect{
		source: &sqlSelectSource{
			expr: expr,
		},
		fields: fields,
	}
}
func (sql *SqlCmdSelect) Where(expr *BoolField) *SqlCmdSelect {
	sql.where = expr
	return sql
}
func (sql *SqlCmdSelect) buildWhere() *sqlCmdSelectResult {
	if sql.where == nil {
		return &sqlCmdSelectResult{
			SqlText: "",
			Args:    []interface{}{},
		}
	}
	whereResult, err := sql.cmp.Resolve(sql.tables, sql.buildContext, sql.where, true)
	if err != nil {
		return &sqlCmdSelectResult{
			Err: err,
		}
	}
	return &sqlCmdSelectResult{
		SqlText: whereResult.Syntax,
		Args:    whereResult.Args,
	}

}
func (sql *SqlCmdSelect) Compile(d DialectCompiler) *sqlCmdSelectResult {
	var args []interface{}
	sql.cmp = Compiler.Ctx(d)

	sourceCompiler, err := sql.source.build(sql.tables, sql.buildContext, d)
	if err != nil {
		return &sqlCmdSelectResult{
			Err: err,
		}
	}
	sql.source.sourceText = sourceCompiler.Syntax
	sql.source.args = sourceCompiler.Args
	args = append(args, sourceCompiler.Args...)

	sql.buildContext = sourceCompiler.buildContext
	resultBuildSelect := sql.buildSelect()
	if resultBuildSelect.Err != nil {
		return resultBuildSelect
	}
	args = append(args, resultBuildSelect.Args...)
	resultBuildWhere := sql.buildWhere()
	if resultBuildWhere.Err != nil {
		return resultBuildWhere
	}

	sqlStr := "SELECT " + resultBuildSelect.SqlText + " FROM " + sql.source.sourceText
	if resultBuildWhere.SqlText != "" {
		sqlStr += " WHERE " + resultBuildWhere.SqlText
	}
	args = append(args, resultBuildWhere.Args...)

	return &sqlCmdSelectResult{
		SqlText: sqlStr,
		Args:    args,
	}

}
