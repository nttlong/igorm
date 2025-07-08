package orm

func (c *CompilerUtils) resolveNumberField(tables *[]string, context *map[string]string, ptr interface{}, requireAlias bool) (*resolverResult, error) {

	if nf, ok := ptr.(*NumberField[int]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(*NumberField[int8]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(*NumberField[int16]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(*NumberField[int32]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(*NumberField[int64]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(*NumberField[uint]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(*NumberField[uint8]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(*NumberField[uint16]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(*NumberField[uint32]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(*NumberField[uint64]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(*NumberField[float32]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(*NumberField[float64]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	//-----------------------------------------
	if nf, ok := ptr.(NumberField[int]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(NumberField[int8]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(NumberField[int16]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(NumberField[int32]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(NumberField[int64]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(NumberField[uint]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(NumberField[uint8]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(NumberField[uint16]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(NumberField[uint32]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(NumberField[uint64]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(NumberField[float32]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	if nf, ok := ptr.(NumberField[float64]); ok {
		return c.Resolve(tables, context, nf.UnderField, requireAlias)
	}
	return nil, nil
}
