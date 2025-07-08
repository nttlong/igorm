package orm

import "strings"

func (c *CompilerUtils) resolveJoinFieldASRightJoin(tables *[]string, context *map[string]string, JoinField joinField, requireAlias bool) (*resolverResult, error) {
	left, err := c.Resolve(tables, context, JoinField.left, requireAlias)
	if err != nil {
		return nil, err
	}
	right, err := c.Resolve(tables, context, JoinField.right, requireAlias)
	if err != nil {
		return nil, err
	}
	last2Tables := (*tables)[0:2]
	leftTable := last2Tables[0]
	if strings.Contains(leftTable, "*") {
		leftTable = strings.Split(leftTable, "*")[0]
	}
	rightTable := last2Tables[1]
	if strings.Contains(rightTable, "*") {
		rightTable = strings.Split(rightTable, "*")[0]
	}
	leftAlias := (*context)[last2Tables[1]]
	rightAlias := (*context)[last2Tables[0]]
	leftExpr := c.Quote(leftTable) + " AS " + c.Quote(leftAlias)
	rightExpr := c.Quote(rightTable) + " AS " + c.Quote(rightAlias)
	syntax := ""
	if left.IsJoinExpr {
		syntax = left.Syntax + " " + JoinField.joinType + " JOIN " + rightExpr + " ON " + left.NextJoin + " = " + right.Syntax

	} else {
		on := left.Syntax + " = " + right.NextJoin
		syntax = right.Syntax + " " + JoinField.joinType + " JOIN " + leftExpr + " ON " + on
		//syntax = leftExpr + " " + JoinField.joinType + " JOIN " + rightExpr + " ON " + left.Syntax + " = " + right.Syntax
	}
	return &resolverResult{
		Syntax:       syntax,
		Args:         append(left.Args, right.Args...),
		buildContext: context,
		Tables:       tables,
		hasNewTable:  true,
		NewTableName: rightTable,
		IsJoinExpr:   true,
		NextJoin:     left.Syntax,
	}, nil

}
func (c *CompilerUtils) resolveJoinField(tables *[]string, context *map[string]string, JoinField joinField, requireAlias bool) (*resolverResult, error) {
	if JoinField.joinType == "RIGHT" {
		return c.resolveJoinFieldASRightJoin(tables, context, JoinField, requireAlias)

	}
	left, err := c.Resolve(tables, context, JoinField.left, requireAlias)
	if err != nil {
		return nil, err
	}
	if jf, ok := JoinField.right.([]interface{}); ok {
		leftTable := (*tables)[len(*tables)-1]
		leftAlias := (*context)[leftTable]
		leftSource := c.Quote(leftTable) + " AS " + c.Quote(leftAlias)
		syntax := leftSource
		args := left.Args
		for _, right := range jf {
			right, err := c.Resolve(tables, context, right, requireAlias)
			args = append(args, right.Args...)
			if err != nil {
				return nil, err
			}
			nextTable := (*tables)[len(*tables)-1]
			nextAlias := (*context)[nextTable]
			nextTable = Utils.GetTableNameFromVirtualName(nextTable)
			nextSourceTable := c.Quote(nextTable) + " AS " + c.Quote(nextAlias)
			if JoinField.joinType == "RIGHT" {
				syntax = syntax + " " + JoinField.joinType + " JOIN " + nextSourceTable + " ON " + left.Syntax + " = " + right.Syntax
			} else {
				syntax = syntax + " " + JoinField.joinType + " JOIN " + nextSourceTable + " ON " + left.Syntax + " = " + right.Syntax
			}

		}
		return &resolverResult{
			Syntax:       syntax,
			Args:         args,
			buildContext: context,
			Tables:       tables,
			IsJoinExpr:   true,
		}, nil

	}

	right, err := c.Resolve(tables, context, JoinField.right, requireAlias)
	if err != nil {
		return nil, err
	}
	last2Tables := (*tables)[len(*tables)-2:]

	leftTable := last2Tables[0]
	if strings.Contains(leftTable, "*") {
		leftTable = strings.Split(leftTable, "*")[0]
	}
	rightTable := last2Tables[1]
	if strings.Contains(rightTable, "*") {
		rightTable = strings.Split(rightTable, "*")[0]
	}

	leftAlias := (*context)[last2Tables[0]]
	rightAlias := (*context)[last2Tables[1]]
	leftExpr := c.Quote(leftTable) + " AS " + c.Quote(leftAlias)
	rightExpr := c.Quote(rightTable) + " AS " + c.Quote(rightAlias)
	syntax := ""
	if left.IsJoinExpr {
		syntax = left.Syntax + " " + JoinField.joinType + " JOIN " + rightExpr + " ON " + left.NextJoin + " = " + right.Syntax

	} else {
		if JoinField.joinType == "RIGHT" {
			syntax = rightExpr + " " + JoinField.joinType + " JOIN " + leftExpr + " ON " + left.Syntax + " = " + right.Syntax
		} else {
			syntax = leftExpr + " " + JoinField.joinType + " JOIN " + rightExpr + " ON " + left.Syntax + " = " + right.Syntax
		}
	}

	return &resolverResult{
		Syntax:       syntax,
		Args:         append(left.Args, right.Args...),
		buildContext: context,
		Tables:       tables,
	}, nil

}
