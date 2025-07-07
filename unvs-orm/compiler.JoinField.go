package orm

func (c *CompilerUtils) resolveJoinField(context *map[string]string, joinField joinField) (*resolverResult, error) {
	c.Resolve(context, joinField.left)
	c.Resolve(context, joinField.right)
	return nil, nil

}
