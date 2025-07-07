package orm

import "fmt"

func (c *CompilerUtils) resolveBinaryField(context *map[string]string, f *fieldBinary) (*resolverResult, error) {

	left, err := c.Resolve(context, f.left)
	if err != nil {
		return nil, err
	}
	right, err := c.Resolve(context, f.right)
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
	if len(args) == 0 {
		return &resolverResult{
			Syntax: fmt.Sprintf("%s %s", left.Syntax, f.op),
			Args:   args,
		}, nil
	} else if len(args) == 1 {
		return &resolverResult{
			Syntax: fmt.Sprintf("%s %s %s", left.Syntax, f.op, right.Syntax),
			Args:   args,
		}, nil
	} else {
		return &resolverResult{
			Syntax: fmt.Sprintf("%s %s %s", left.Syntax, f.op, right.Syntax),
			Args:   args,
		}, nil
	}
}
