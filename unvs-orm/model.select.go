package orm

import (
	EXPR "unvs-orm/expr"
)

type exprField struct {
	UnderField interface{}
	Stmt       string
	Args       []interface{}
}

func (f *Base) Expr(expr string, args ...interface{}) *exprField {
	return &exprField{
		Stmt: expr,
		Args: args,
	}
}
func (m *Model[T]) Select(fields ...interface{}) *SqlSelectBuilder {
	return &SqlSelectBuilder{
		tableName: m.TableName,
		selects:   fields,
		noAlias:   true,
		tables:    &[]string{m.TableName},
		context:   &map[string]string{m.TableName: m.TableName},
	}
}
func (sql *SqlSelectBuilder) Compile(d DialectCompiler) *SqlStmt {
	// convert to SqlCmdSelect then compile
	cmd := &SqlCmdSelect{
		source:       sql.tableName,
		fields:       sql.selects,
		where:        sql.condition,
		buildContext: sql.context,
		cmp:          d.getCompiler(),
		tables:       sql.tables,
		exprCmp:      EXPR.NewExpressionCompiler(d.driverName()),
	}
	return cmd.Compile(d)
}
func (m *Model[T]) Filter(expr *BoolField) *SqlSelectBuilder {
	return &SqlSelectBuilder{
		source:    m.TableName,
		condition: expr,
		noAlias:   true,
	}
}
