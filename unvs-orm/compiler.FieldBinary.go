package orm

import (
	"fmt"
)

func (c *CompilerUtils) resolveBinaryField(tables *[]string, context *map[string]string, f *fieldBinary, extractAlias, applyContext bool) (*resolverResult, error) {

	left, err := c.Resolve(tables, context, f.left, extractAlias, applyContext)
	if err != nil {
		return nil, err
	}
	right, err := c.Resolve(tables, context, f.right, extractAlias, applyContext)
	if err != nil {
		return nil, err
	}
	args := append(left.Args, right.Args...)
	if f.op == "BETWEEN" || f.op == "NOT BETWEEN" {
		return &resolverResult{
			Syntax: fmt.Sprintf("%s %s ? AND ?", left.Syntax, f.op),
			Args:   args,
		}, nil
	}
	if f.op == "IN" || f.op == "NOT IN" {
		return &resolverResult{
			Syntax: fmt.Sprintf("%s %s (%s)", left.Syntax, f.op, right.Syntax),
			Args:   args,
		}, nil
	}

	if left.Syntax != "" && right.Syntax != "" {

		return &resolverResult{
			Syntax: fmt.Sprintf("%s %s %s", left.Syntax, f.op, right.Syntax),
			Args:   args,
		}, nil
	} else if left.Syntax == "" && right.Syntax != "" {
		return &resolverResult{
			Syntax: fmt.Sprintf("%s %s", f.op, right.Syntax),
			Args:   args,
		}, nil

	} else if left.Syntax != "" && right.Syntax == "" {
		return &resolverResult{
			Syntax: fmt.Sprintf("%s %s", left.Syntax, f.op),
			Args:   args,
		}, nil
	} else {
		return nil, fmt.Errorf("invalid binary expression")
	}
}
