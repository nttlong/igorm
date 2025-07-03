package orm

func (c *CompilerUtils) resolveNumberField(ptr interface{}) (*resolverResult, error) {

	if nf, ok := ptr.(*NumberField[int]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(*NumberField[int8]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(*NumberField[int16]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(*NumberField[int32]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(*NumberField[int64]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(*NumberField[uint]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(*NumberField[uint8]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(*NumberField[uint16]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(*NumberField[uint32]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(*NumberField[uint64]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(*NumberField[float32]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(*NumberField[float64]); ok {
		return c.Resolve(nf.dbField)
	}
	//-----------------------------------------
	if nf, ok := ptr.(NumberField[int]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(NumberField[int8]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(NumberField[int16]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(NumberField[int32]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(NumberField[int64]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(NumberField[uint]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(NumberField[uint8]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(NumberField[uint16]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(NumberField[uint32]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(NumberField[uint64]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(NumberField[float32]); ok {
		return c.Resolve(nf.dbField)
	}
	if nf, ok := ptr.(NumberField[float64]); ok {
		return c.Resolve(nf.dbField)
	}
	return nil, nil
}
