package orm

func (c *CompilerUtils) resolveDBField(aliasSource *map[string]string, f *dbField) (*resolverResult, error) {
	if aliasSource == nil {
		return &resolverResult{
			Syntax: c.Quote(f.Table, f.Name),
			Args:   nil,
		}, nil
	}
	if alias, ok := (*aliasSource)[f.Table]; ok {
		return &resolverResult{
			Syntax:      c.Quote(alias, f.Name),
			Args:        nil,
			AliasSource: *aliasSource,
		}, nil
	}
	return &resolverResult{
		Syntax: c.Quote(f.Table, f.Name),
		Args:   nil,
	}, nil
}
