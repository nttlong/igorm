package orm

import (
	"fmt"
)

func (c *CompilerUtils) resolveJoinField(tables *[]string, context *map[string]string, joinField joinField, requireAlias bool) (*resolverResult, error) {
	left, err := c.Resolve(tables, context, joinField.left, requireAlias)
	if err != nil {
		return nil, err
	}

	right, err := c.Resolve(tables, context, joinField.right, requireAlias)
	if err != nil {
		return nil, err
	}
	fmt.Println(left.Syntax, right.Syntax)
	leftTable := (*left.Tables)[0]
	rightTable := (*right.Tables)[0]
	leftAlias := (*left.buildContext)[leftTable]
	rightAlias := (*right.buildContext)[rightTable]
	leftExpr := c.Quote(leftTable) + " AS " + c.Quote(leftAlias)
	rightExpr := c.Quote(rightTable) + " AS " + c.Quote(rightAlias)
	syntax := leftExpr + " " + joinField.joinType + " JOIN " + rightExpr + " ON " + left.Syntax + " = " + right.Syntax

	return &resolverResult{
		Syntax:       syntax,
		Args:         append(left.Args, right.Args...),
		buildContext: context,
		Tables:       tables,
	}, nil

}
