package orm

import (
	"fmt"
	"strings"
)

func (c *CompilerUtils) resolveBoolFieldRightJoin(tables *[]string, context *map[string]string, f *joinField, requireAlias bool) (*resolverResult, error) {
	left, err := c.Resolve(tables, context, f.left, requireAlias)
	if err != nil {
		return nil, err
	}
	right, err := c.Resolve(tables, context, f.right, requireAlias)
	if err != nil {
		return nil, err
	}
	args := append(left.Args, right.Args...)
	last2Table := (*tables)[len((*tables))-2:]
	leftTable := last2Table[1]
	rightTable := last2Table[0]
	if strings.Contains(leftTable, "*") {
		leftTable = strings.Split(leftTable, "*")[0]
	}
	if strings.Contains(rightTable, "*") {
		rightTable = strings.Split(rightTable, "*")[0]
	}
	leftSource := c.Quote(leftTable) + " AS " + c.Quote((*context)[last2Table[1]])
	rightSource := c.Quote(rightTable) + " AS " + c.Quote((*context)[last2Table[0]])
	onSyntax := left.Syntax + " = " + right.Syntax
	syntax := leftSource + " " + f.joinType + " JOIN " + rightSource + " ON " + onSyntax

	return &resolverResult{
		Syntax:     syntax, //fmt.Sprintf("%s %s %s ON %s", left.Syntax, f.op, right.Syntax, f.rawText),
		Args:       args,
		IsJoinExpr: true,
		NextJoin:   left.Syntax,
	}, nil
}
func (c *CompilerUtils) resolveBoolFieldJoin(tables *[]string, context *map[string]string, bf *BoolField, requireAlias bool) (*resolverResult, error) {
	if f, ok := bf.underField.(*joinField); ok {
		if f.joinType == "RIGHT" {
			return c.resolveBoolFieldRightJoin(tables, context, f, requireAlias)

		}
		left, err := c.Resolve(tables, context, f.left, requireAlias)
		if err != nil {
			return nil, err
		}
		right, err := c.Resolve(tables, context, f.right, requireAlias)
		if err != nil {
			return nil, err
		}
		args := append(left.Args, right.Args...)
		last2Table := (*tables)[len((*tables))-2:]
		leftTable := last2Table[0]
		rightTable := last2Table[1]
		if strings.Contains(leftTable, "*") {
			leftTable = strings.Split(leftTable, "*")[0]
		}
		if strings.Contains(rightTable, "*") {
			rightTable = strings.Split(rightTable, "*")[0]
		}
		leftSource := c.Quote(leftTable) + " AS " + c.Quote((*context)[last2Table[0]])
		rightSource := c.Quote(rightTable) + " AS " + c.Quote((*context)[last2Table[1]])
		onSyntax := left.Syntax + " = " + right.Syntax
		syntax := ""
		if f.joinType == "RIGHT" {
			syntax = rightSource + " " + f.joinType + " JOIN " + leftSource + " ON " + onSyntax
		} else {
			syntax = leftSource + " " + f.joinType + " JOIN " + rightSource + " ON " + onSyntax
		}

		return &resolverResult{
			Syntax:     syntax, //fmt.Sprintf("%s %s %s ON %s", left.Syntax, f.op, right.Syntax, f.rawText),
			Args:       args,
			IsJoinExpr: true,
			NextJoin:   right.Syntax,
		}, nil
	}
	return nil, fmt.Errorf("unsupported expression type: %T, file %s, line %d", bf.underField, "unvs-orm/compiler.go", 23)
}
