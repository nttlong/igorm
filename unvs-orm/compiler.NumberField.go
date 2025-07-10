package orm

func (c *CompilerUtils) resolveNumberField(
	tables *[]string,
	context *map[string]string,
	ptr interface{},
	extractAlias, applyContext bool,
) (*resolverResult, error) {

	switch nf := ptr.(type) {
	case *NumberField[int],
		*NumberField[int8],
		*NumberField[int16],
		*NumberField[int32],
		*NumberField[int64],
		*NumberField[uint],
		*NumberField[uint8],
		*NumberField[uint16],
		*NumberField[uint32],
		*NumberField[uint64],
		*NumberField[float32],
		*NumberField[float64]:

		return c.Resolve(tables, context, nf.(interface{ GetUnderField() any }).GetUnderField(), extractAlias, applyContext)

	case NumberField[int],
		NumberField[int8],
		NumberField[int16],
		NumberField[int32],
		NumberField[int64],
		NumberField[uint],
		NumberField[uint8],
		NumberField[uint16],
		NumberField[uint32],
		NumberField[uint64],
		NumberField[float32],
		NumberField[float64]:

		return c.Resolve(tables, context, nf.(interface{ GetUnderField() any }).GetUnderField(), extractAlias, applyContext)

	default:
		return nil, nil
	}
}
