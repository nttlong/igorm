package orm

func (c *CompilerUtils) resolveNumberField(expr interface{}) (*resolverResult, error) {

	switch f := expr.(type) {
	case *NumberField[int], NumberField[int],
		*NumberField[int8], NumberField[int8],
		*NumberField[int16], NumberField[int16],
		*NumberField[int32], NumberField[int32],
		*NumberField[int64], NumberField[int64],
		*NumberField[uint], NumberField[uint],
		*NumberField[uint8], NumberField[uint8],
		*NumberField[uint16], NumberField[uint16],
		*NumberField[uint32], NumberField[uint32],
		*NumberField[uint64], NumberField[uint64],
		*NumberField[float32], NumberField[float32],
		*NumberField[float64], NumberField[float64]:

		// Lấy địa chỉ field nếu là value type
		ptr := any(f)
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
	}
	return nil, nil
}
