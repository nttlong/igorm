package orm

import (
	"fmt"
	"strings"
	"sync"
)

type JoinCompilerUtils struct {
	dialect DialectCompiler
}

var cacheJoinCompilerCtx sync.Map

func (c *JoinCompilerUtils) Ctx(dialect DialectCompiler) *JoinCompilerUtils {
	if dialect == nil {
		panic("dialect is nil")
	}

	key := dialect.driverName()

	if v, ok := cacheJoinCompilerCtx.Load(key); ok {
		return v.(*JoinCompilerUtils)
	}
	ret := &JoinCompilerUtils{dialect: dialect}
	dialect.setJoinCompiler(ret)
	cacheJoinCompilerCtx.Store(key, ret)
	return ret
}
func (c *JoinCompilerUtils) Resolve(expr *JoinExpr) (*resolverResult, error) {
	cmd := Compiler.Ctx(c.dialect)

	stack := []*JoinExpr{}
	for node := expr; node != nil; node = node.previous {
		stack = append(stack, node)
	}
	for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		stack[i], stack[j] = stack[j], stack[i]
	}

	sqlParts := []string{}
	args := []interface{}{}

	// FROM part
	baseAlias := expr.aliasMap[expr.baseTable]
	sqlParts = append(sqlParts, cmd.Quote(expr.baseTable)+" AS "+cmd.Quote(baseAlias))

	for _, join := range stack {
		right := join.rightTable
		rightAlias := join.aliasMap[right]

		onRes, err := cmd.Resolve(&join.aliasMap, join.on)
		if err != nil {
			return nil, err
		}

		sql := fmt.Sprintf("%s JOIN %s AS %s ON %s",
			join.joinType,
			cmd.Quote(right),
			cmd.Quote(rightAlias),
			onRes.Syntax,
		)
		sqlParts = append(sqlParts, sql)
		args = append(args, onRes.Args...)
	}

	return &resolverResult{
		Syntax: strings.Join(sqlParts, " "),
		Args:   args,
	}, nil
}

// func (c *JoinCompilerUtils) Resolve(expr *JoinExpr) (*resolverResult, error) {
// 	cmd := Compiler.Ctx(c.dialect)

// 	// Dịch điều kiện ON
// 	onRes, err := cmd.Resolve(&expr.aliasSource, expr.on)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Bảng được JOIN (phải là bảng bên phải)
// 	// Giả sử expr.tables luôn là [leftTable, rightTable]
// 	if len(expr.tables) != 2 {
// 		return nil, errors.New("JoinExpr expects exactly 2 tables")
// 	}

// 	// Bảng bên phải là bảng được JOIN vào FROM gốc
// 	rightTable := expr.tables[1]
// 	rightAlias := expr.aliasSource[rightTable]

// 	// Sinh JOIN clause
// 	tblJoin := cmd.Quote(rightTable) + " AS " + cmd.Quote(rightAlias)

// 	ret := &resolverResult{
// 		Syntax:      expr.joinType + " JOIN " + tblJoin + " ON " + onRes.Syntax,
// 		Args:        onRes.Args,
// 		AliasSource: expr.aliasSource,
// 	}
// 	return ret, nil
// }

var JoinCompiler = JoinCompilerUtils{}
