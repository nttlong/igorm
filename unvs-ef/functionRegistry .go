package unvsef

import "fmt"

type functionRegistry struct {
	Dialect Dialect
}

var Funcs = &functionRegistry{}

func (f *functionRegistry) Len(arg Expr) Expr {
	return f.Dialect.Func("LEN", arg) // SQL Server
}

func (f *functionRegistry) Lower(arg Expr) Expr {
	return f.Dialect.Func("LOWER", arg)
}

func (f *functionRegistry) Upper(arg Expr) Expr {
	return f.Dialect.Func("UPPER", arg)
}
func (f *functionRegistry) FullTextContains(col Expr, keyword string) Expr {
	switch f.Dialect.(type) {
	case *SqlServerDialect:
		return Raw(fmt.Sprintf("CONTAINS(%s, ?)", toSQLExpr(col, f.Dialect)), keyword)
	case *PostgresDialect:
		return Raw(fmt.Sprintf("%s @@ to_tsquery(?)", toSQLExpr(col, f.Dialect)), keyword)
	case *MySqlDialect:
		return Raw(fmt.Sprintf("MATCH(%s) AGAINST (?)", toSQLExpr(col, f.Dialect)), keyword)
	default:
		panic("Unsupported dialect for full-text search")
	}
}
func toSQLExpr(e Expr, d Dialect) string {
	sql, _ := e.ToSQL(d)
	return sql
}
