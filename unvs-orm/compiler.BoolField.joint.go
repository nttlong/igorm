package orm

import (
	"fmt"
)

func (c *CompilerUtils) resolveBoolFieldJoin(aliasSource *map[string]string, f *BoolField) (*resolverResult, error) {
	if f.op == "INNER JOIN" {
		left, err := c.Resolve(aliasSource, f.left)
		if err != nil {
			return nil, err
		}
		right, err := c.Resolve(aliasSource, f.right)
		if err != nil {
			return nil, err
		}
		args := append(left.Args, right.Args...)
		rightSource := c.Quote(f.joinSource) + " AS " + c.Quote(f.joinSourceAlias)
		syntax := left.Syntax + " " + f.op + " " + rightSource + " ON " + right.Syntax
		return &resolverResult{
			Syntax: syntax, //fmt.Sprintf("%s %s %s ON %s", left.Syntax, f.op, right.Syntax, f.rawText),
			Args:   args,
		}, nil
	}

	if f.op == "LEFT JOIN" {
		panic("not implemented LEFT JOIN at complex.go:239")
	}
	if f.op == "LEFT JOIN" {
		panic("not implemented LEFT JOIN at complex.go:242")
	}
	return nil, fmt.Errorf("unsupported join type: %s", f.op)
}
