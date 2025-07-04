package orm

func (c *CompilerUtils) resolveNumberField(aliasSource *map[string]string, ptr interface{}) (*resolverResult, error) {

	if nf, ok := ptr.(*NumberField[int]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(*NumberField[int8]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(*NumberField[int16]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(*NumberField[int32]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(*NumberField[int64]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(*NumberField[uint]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(*NumberField[uint8]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(*NumberField[uint16]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(*NumberField[uint32]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(*NumberField[uint64]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(*NumberField[float32]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(nil, nf.callMethod)
		}
	}
	if nf, ok := ptr.(*NumberField[float64]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	//-----------------------------------------
	if nf, ok := ptr.(NumberField[int]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(NumberField[int8]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(NumberField[int16]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(NumberField[int32]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(NumberField[int64]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(NumberField[uint]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(NumberField[uint8]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(NumberField[uint16]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(NumberField[uint32]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(NumberField[uint64]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(NumberField[float32]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	if nf, ok := ptr.(NumberField[float64]); ok {
		if nf.dbField != nil {
			return c.Resolve(aliasSource, nf.dbField)
		}
		if nf.callMethod != nil {
			return c.Resolve(aliasSource, nf.callMethod)
		}
	}
	return nil, nil
}
