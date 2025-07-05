package orm

import (
	"errors"
	"strings"
)

type sqlSelectSource struct {
	expr       *BoolField
	sourceText string
	args       []interface{}
}
type sqlSelectFields struct {
	fields     []interface{}
	aliasMap   *map[string]string
	selectText string
	args       []interface{}
}

func (s *sqlSelectSource) build(d DialectCompiler) (*resolverResult, error) {
	ctx := JoinCompiler.Ctx(d)
	return ctx.ResolveBoolFieldAsJoin(s.expr)

}

type SqlCmdSelect struct {
	source   *sqlSelectSource
	fields   []interface{}
	aliasMap *map[string]string
	cmp      *CompilerUtils
	Err      error
}

func (sql *SqlCmdSelect) buildSelect() *sqlCmdSelectResult {

	fields := []string{}
	args := []interface{}{}
	for _, field := range sql.fields {
		if _, ok := field.(*aliasField); ok {
			fieldCompiler, err := sql.cmp.Resolve(sql.aliasMap, field)
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
func (sql *SqlCmdSelect) Compile(d DialectCompiler) *sqlCmdSelectResult {
	var args []interface{}
	sql.cmp = Compiler.Ctx(d)

	sourceCompiler, err := sql.source.build(d)
	if err != nil {
		return &sqlCmdSelectResult{
			Err: err,
		}
	}
	sql.source.sourceText = sourceCompiler.Syntax
	sql.source.args = sourceCompiler.Args
	args = append(args, sourceCompiler.Args...)

	sql.aliasMap = &sourceCompiler.AliasSource
	resultBuildSelect := sql.buildSelect()
	if resultBuildSelect.Err != nil {
		return resultBuildSelect
	}
	args = append(args, resultBuildSelect.Args...)

	sqlStr := "SELECT " + resultBuildSelect.SqlText + " FROM " + sql.source.sourceText
	return &sqlCmdSelectResult{
		SqlText: sqlStr,
		Args:    args,
	}

}
