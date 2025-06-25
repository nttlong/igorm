// Package unvsef provides a type-safe SQL query builder using Go generics.
// It supports multiple SQL dialects, composable expressions, aggregates,
// binary operations, CASE WHEN expressions, JOINs, GROUP BY, HAVING, and more.
package unvsef

import (
	"fmt"
	"strings"
)

// Expr is the core interface representing any SQL expression.
type Expr interface {
	ToSQL(d Dialect) (string, []interface{})
}

// --------------------- Literal Expression ---------------------
// Literal represents a constant value in SQL.
type Literal[T any] struct {
	Value T
}

// Lit is a helper to construct a typed literal.
func Lit[T any](v T) Literal[T] {
	return Literal[T]{Value: v}
}

// --------------------- Aggregate Expression ---------------------
// AggregateFunc represents functions like SUM, COUNT, etc.
type AggregateFunc struct {
	Func  string
	Arg   Expr
	Alias *string
}

func (a AggregateFunc) ToSQL(d Dialect) (string, []interface{}) {
	sql, args := a.Arg.ToSQL(d)
	if a.Alias != nil {
		return fmt.Sprintf("%s(%s) AS %s", a.Func, sql, *a.Alias), args
	}
	return fmt.Sprintf("%s(%s)", a.Func, sql), args
}

// Helper functions for aggregates.
func Sum(arg Expr) AggregateFunc   { return AggregateFunc{Func: "SUM", Arg: arg} }
func Count(arg Expr) AggregateFunc { return AggregateFunc{Func: "COUNT", Arg: arg} }

// --------------------- Binary Expression ---------------------
// BinaryExpr represents arithmetic or logical operations between two expressions.
type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
	Alias *string
}

func (b BinaryExpr) ToSQL(d Dialect) (string, []interface{}) {
	lsql, largs := b.Left.ToSQL(d)
	rsql, rargs := b.Right.ToSQL(d)
	sql := fmt.Sprintf("(%s %s %s)", lsql, b.Op, rsql)
	if b.Alias != nil {
		sql += " AS " + *b.Alias
	}
	return sql, append(largs, rargs...)
}

// Arithmetic helpers.
func Add(left, right Expr) *BinaryExpr { return &BinaryExpr{Left: left, Op: "+", Right: right} }
func Sub(left, right Expr) *BinaryExpr { return &BinaryExpr{Left: left, Op: "-", Right: right} }
func Mul(left, right Expr) *BinaryExpr { return &BinaryExpr{Left: left, Op: "*", Right: right} }
func Div(left, right Expr) *BinaryExpr { return &BinaryExpr{Left: left, Op: "/", Right: right} }

// As sets an alias for the expression.
func (b *BinaryExpr) As(alias string) *BinaryExpr {
	b.Alias = &alias
	return b
}

// --------------------- Case Expression ---------------------
// CaseExpr represents SQL CASE WHEN THEN ELSE END expression.
type CaseExpr struct {
	whens []struct {
		Cond Expr
		Then Expr
	}
	elseExpr Expr
	alias    *string
}

// Case starts a new CASE expression.
func Case() *CaseExpr {
	return &CaseExpr{}
}

// When adds a WHEN condition.
func (c *CaseExpr) When(cond Expr, then Expr) *CaseExpr {
	c.whens = append(c.whens, struct {
		Cond Expr
		Then Expr
	}{cond, then})
	return c
}

// Else adds the ELSE fallback.
func (c *CaseExpr) Else(e Expr) *CaseExpr {
	c.elseExpr = e
	return c
}

// As sets an alias for the CASE expression.
func (c *CaseExpr) As(alias string) *CaseExpr {
	c.alias = &alias
	return c
}

// ToSQL generates the SQL string and arguments for the CASE expression.
func (c *CaseExpr) ToSQL(d Dialect) (string, []interface{}) {
	var sb strings.Builder
	args := []interface{}{}
	sb.WriteString("CASE")
	for _, w := range c.whens {
		condSQL, condArgs := w.Cond.ToSQL(d)
		thenSQL, thenArgs := w.Then.ToSQL(d)
		sb.WriteString(fmt.Sprintf(" WHEN %s THEN %s", condSQL, thenSQL))
		args = append(args, condArgs...)
		args = append(args, thenArgs...)
	}
	if c.elseExpr != nil {
		elseSQL, elseArgs := c.elseExpr.ToSQL(d)
		sb.WriteString(" ELSE " + elseSQL)
		args = append(args, elseArgs...)
	}
	sb.WriteString(" END")
	if c.alias != nil {
		sb.WriteString(" AS " + *c.alias)
	}
	return sb.String(), args
}

// --------------------- Join Clause ---------------------
// JoinClause represents a JOIN statement.
type JoinClause struct {
	Table    string
	On       Expr
	JoinType string
}

func (j JoinClause) ToSQL(d Dialect) (string, []interface{}) {
	onSQL, args := j.On.ToSQL(d)
	return fmt.Sprintf("%s JOIN %s ON %s", j.JoinType, j.Table, onSQL), args
}

// --------------------- Query Object ---------------------
// Query represents a full SQL SELECT query.
type Query struct {
	selects []Expr
	from    string
	joins   []JoinClause
	where   Expr
	groupBy []Expr
	having  Expr
}

// NewQuery creates an empty query builder.
func NewQuery() *Query {
	return &Query{}
}

func (q *Query) Select(exprs ...Expr) *Query {
	q.selects = exprs
	return q
}

func (q *Query) From(table string) *Query {
	q.from = table
	return q
}

func (q *Query) Join(table string, on Expr) *Query {
	q.joins = append(q.joins, JoinClause{Table: table, On: on, JoinType: "INNER"})
	return q
}

func (q *Query) LeftJoin(table string, on Expr) *Query {
	q.joins = append(q.joins, JoinClause{Table: table, On: on, JoinType: "LEFT"})
	return q
}

func (q *Query) Where(expr Expr) *Query {
	q.where = expr
	return q
}

func (q *Query) GroupBy(exprs ...Expr) *Query {
	q.groupBy = exprs
	return q
}

func (q *Query) Having(expr Expr) *Query {
	q.having = expr
	return q
}

// ToSQL generates the full SQL SELECT statement.
func (q *Query) ToSQL(d Dialect) (string, []interface{}) {
	var sb strings.Builder
	args := []interface{}{}

	sb.WriteString("SELECT ")
	parts := []string{}
	for _, e := range q.selects {
		sql, a := e.ToSQL(d)
		parts = append(parts, sql)
		args = append(args, a...)
	}
	sb.WriteString(strings.Join(parts, ", "))
	sb.WriteString(" FROM " + q.from)

	for _, j := range q.joins {
		jsql, jargs := j.ToSQL(d)
		sb.WriteString(" " + jsql)
		args = append(args, jargs...)
	}

	if q.where != nil {
		wsql, wargs := q.where.ToSQL(d)
		sb.WriteString(" WHERE " + wsql)
		args = append(args, wargs...)
	}

	if len(q.groupBy) > 0 {
		gbParts := []string{}
		for _, g := range q.groupBy {
			sql, _ := g.ToSQL(d)
			gbParts = append(gbParts, sql)
		}
		sb.WriteString(" GROUP BY " + strings.Join(gbParts, ", "))
	}

	if q.having != nil {
		hsql, hargs := q.having.ToSQL(d)
		sb.WriteString(" HAVING " + hsql)
		args = append(args, hargs...)
	}

	return sb.String(), args
}

// --------------------- Utility ---------------------
// toExpr ensures that a raw value is wrapped as a Literal if not already an Expr.
func toExpr(val interface{}) Expr {
	switch v := val.(type) {
	case Expr:
		return v
	default:
		return Literal[interface{}]{Value: v}
	}
}
