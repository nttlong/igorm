package orm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	EXPR "unvs-orm/expr"
)

type sqlSelectFields struct {
	fields     []interface{}
	aliasMap   *map[string]string
	selectText string
	args       []interface{}
}

type SqlCmdSelect struct {
	source       interface{}
	fields       []interface{}
	where        interface{}
	groups       []interface{}
	buildContext *map[string]string
	having       *BoolField
	cmp          *CompilerUtils
	Err          error
	tables       *[]string
	exprCmp      *EXPR.EXPR
}

func (sql *SqlCmdSelect) getSelectField(field interface{}) string {
	v := reflect.ValueOf(field)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	fieldNameField := v.FieldByName("Name")
	if !fieldNameField.IsValid() {
		return ""
	}
	fieldName := fieldNameField.String()
	return fieldName

}
func (sql *SqlCmdSelect) buildSelect() *sqlCmdSelectResult {

	fields := []string{}
	args := []interface{}{}
	for _, field := range sql.fields {
		if _, ok := field.(*aliasField); ok {
			fieldCompiler, err := sql.cmp.Resolve(sql.tables, sql.buildContext, field, false, true)
			if err != nil {
				return &sqlCmdSelectResult{
					Err: err,
				}
			}
			fields = append(fields, fieldCompiler.Syntax)
			args = append(args, fieldCompiler.Args...)
		} else if exprField, ok := field.(*exprField); ok {
			if sql.tables == nil {
				sql.tables = &[]string{}
			}
			if sql.buildContext == nil {
				sql.buildContext = &map[string]string{}
			}

			exprCmp := EXPR.NewExpressionCompiler(sql.cmp.dialect.driverName())
			compileResult, err := exprCmp.Compile("", sql.tables, sql.buildContext, exprField.Stmt, false, true)
			if err != nil {
				return &sqlCmdSelectResult{
					Err: err,
				}
			}
			fields = append(fields, compileResult.Syntax)

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

type SqlStmt struct {
	From    string
	Where   string
	GroupBy string
	Having  string
	OrderBy string
	Limit   string
	Offset  string
	Select  string
	Args    []interface{}
	err     error
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
		source: expr,
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
	whereResult, err := sql.cmp.Resolve(sql.tables, sql.buildContext, sql.where, false, true)
	if err != nil {
		return &sqlCmdSelectResult{
			Err: err,
		}
	}
	if whereResult != nil {
		return &sqlCmdSelectResult{
			SqlText: whereResult.Syntax,
			Args:    whereResult.Args,
		}
	} else {
		return &sqlCmdSelectResult{
			SqlText: "",
			Args:    []interface{}{},
		}
	}

}
func (sql *SqlCmdSelect) compileSourceByJoinExpr(d DialectCompiler, je *JoinExpr, sourceCache string) *SqlStmt {
	sql.tables = &[]string{}
	sql.buildContext = &map[string]string{}
	join, err := JoinCompiler.Ctx(d).Resolve(je, sql.tables, sql.buildContext, sourceCache)
	if err != nil {
		return &SqlStmt{
			err: err,
		}
	}

	ret := &SqlStmt{
		From: join.Syntax,
		Args: join.Args,
	}

	ret.Select = sql.buildSelect2(join.Syntax)
	return ret
}
func (sql *SqlCmdSelect) Compile(d DialectCompiler) *SqlStmt {
	if sql.exprCmp == nil {
		sql.exprCmp = EXPR.NewExpressionCompiler(d.driverName())
	}
	if txtSource, ok := sql.source.(string); ok {
		return sql.simpleSource(txtSource, d)

	}
	if exprSource, ok := sql.source.(*BoolField); ok {
		if je, ok := exprSource.underField.(*JoinExpr); ok {

			return sql.compileSourceByJoinExpr(d, je, "")

		}
		retCmp, err := JoinCompiler.Ctx(d).ResolveBoolFieldAsJoin(sql.tables, sql.buildContext, exprSource)
		// retCmp, err := sql.cmp.Resolve(sql.tables, sql.buildContext, exprSource, true, true)
		if err != nil {
			return &SqlStmt{
				err: err,
			}
		}
		selectStmt := sql.buildSelect()
		ret := &SqlStmt{
			From: retCmp.Syntax,

			Select: selectStmt.SqlText,
			Args:   append(selectStmt.Args, retCmp.Args...),
		}
		return ret

	}

	panic(fmt.Errorf("not implemented yet"))

}
func (sql *SqlCmdSelect) GroupBy(fields ...interface{}) *SqlCmdSelect {
	sql.groups = fields
	return sql
}
func (sql *SqlCmdSelect) Having(expr *BoolField) *SqlCmdSelect {
	sql.having = expr
	return sql
}
