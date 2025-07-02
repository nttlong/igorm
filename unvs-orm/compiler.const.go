package orm

func (c *CompilerUtils) resolveConstant(expr interface{}) (*resolverResult, error) {
	switch expr.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
		return &resolverResult{
			Syntax: "?",
			Args:   []interface{}{expr},
		}, nil
	case bool:
		return &resolverResult{
			Syntax: "?",
			Args:   []interface{}{expr},
		}, nil
	}
	return nil, nil
}
