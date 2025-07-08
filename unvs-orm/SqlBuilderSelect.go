package orm

import (
	"errors"
	"reflect"
	"strings"
	"unvs-orm/internal"
)

type SqlSelectBuilder struct {
	source    interface{} //<-- table name or join expression or subquery
	condition *BoolField  //<-- condition expression
	selects   []interface{}
	group     []interface{}

	having  interface{}
	Err     error
	noAlias bool
	tables  *[]string
}
type SqlCompilerResult struct {
	Sql  string
	Args []interface{}
}

func (s *SqlSelectBuilder) Where(condition *BoolField) *SqlSelectBuilder {
	s.condition = condition
	return s
}
func (s *SqlSelectBuilder) Select(fields ...interface{}) *SqlSelectBuilder {
	s.selects = fields
	return s
}
func (s *SqlSelectBuilder) GroupBy(fields ...interface{}) *SqlSelectBuilder {
	s.group = fields
	return s

}

func (s *SqlSelectBuilder) Having(condition *BoolField) *SqlSelectBuilder {
	s.having = condition
	return s
}

func (s *SqlSelectBuilder) ToSql(dialectCompiler DialectCompiler) (*SqlCompilerResult, error) {
	if s.source == nil {
		return nil, errors.New("source is nil")
	}
	var source *resolverResult
	ctx := Compiler.Ctx(dialectCompiler)
	joinCtx := JoinCompiler.Ctx(dialectCompiler)
	if strTableName, ok := s.source.(string); ok {
		if s.noAlias {
			source = &resolverResult{
				Syntax:       ctx.Quote(strTableName),
				buildContext: nil,
				Args:         []interface{}{},
			}
		} else {
			buildContext := map[string]string{strTableName: "T1"}
			source = &resolverResult{
				Syntax:       ctx.Quote(strTableName) + " AS " + ctx.Quote("T1"),
				buildContext: &buildContext,
				Args:         []interface{}{},
			}
		}
	} else if join, ok := s.source.(*JoinExpr); ok {
		joinResult, err := joinCtx.Resolve(join)
		if err != nil {
			return nil, err
		}

		source = &resolverResult{
			Syntax:       joinResult.Syntax,
			buildContext: &map[string]string{},
			Args:         joinResult.Args,
		}
	} else {
		typ := reflect.TypeOf(s.source)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		tableName := internal.Utils.TableNameFromStruct(typ)
		buildContext := map[string]string{tableName: "T1"}

		source = &resolverResult{
			Syntax:       tableName,
			buildContext: &buildContext,
			Args:         []interface{}{},
		}
	}

	var condition *resolverResult
	if s.condition != nil {
		_condition, err := ctx.Resolve(s.tables, source.buildContext, s.condition, true)
		if err != nil {
			return nil, err
		}
		condition = _condition
	}

	selects := []resolverResult{}
	for _, selectField := range s.selects {
		field, err := ctx.Resolve(s.tables, source.buildContext, selectField, !s.noAlias)
		if err != nil {
			return nil, err
		}
		selects = append(selects, *field)
	}

	group := []resolverResult{}
	for _, groupField := range s.group {
		field, err := ctx.Resolve(s.tables, source.buildContext, groupField, true)
		if err != nil {
			return nil, err
		}
		group = append(group, *field)
	}

	var having *resolverResult
	if s.having != nil {
		_having, err := ctx.Resolve(s.tables, source.buildContext, s.having, true)
		if err != nil {
			return nil, err
		}
		having = _having
	}

	// === Assemble SQL ===
	sql := "SELECT "

	if len(selects) == 0 {
		sql += "*"
	} else {
		selectParts := []string{}
		for _, sel := range selects {
			selectParts = append(selectParts, sel.Syntax)
		}
		sql += joinComma(selectParts)
	}

	sql += " FROM " + source.Syntax
	args := append([]interface{}{}, source.Args...)

	if condition != nil {
		sql += " WHERE " + condition.Syntax
		args = append(args, condition.Args...)
	}

	if len(group) > 0 {
		groupParts := []string{}
		for _, grp := range group {
			groupParts = append(groupParts, grp.Syntax)
		}
		sql += " GROUP BY " + joinComma(groupParts)
		for _, g := range group {
			args = append(args, g.Args...)
		}
	}

	if having != nil {
		sql += " HAVING " + having.Syntax
		args = append(args, having.Args...)
	}

	return &SqlCompilerResult{
		Sql:  sql,
		Args: args,
	}, nil
}
func joinComma(parts []string) string {
	return strings.Join(parts, ", ")
}
