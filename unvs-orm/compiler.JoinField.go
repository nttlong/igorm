package orm

func (c *CompilerUtils) resolveJoinField(tables *[]string, context *map[string]string, joinField joinField, requireAlias bool) (*resolverResult, error) {
	c.Resolve(tables, context, joinField.left, requireAlias)
	c.Resolve(tables, context, joinField.right, requireAlias)
	return nil, nil

}
