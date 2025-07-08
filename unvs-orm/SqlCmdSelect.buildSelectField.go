package orm

import (
	"fmt"
)

func (sql *SqlCmdSelect) getFieldName(field interface{}) string {
	if f, ok := field.(*dbField); ok {
		return f.Name
	}
	panic(fmt.Errorf("unsupported field type %T, file orm/SqlCmdSelect.buildSelectField.go, line 11", field))
}
func (sql *SqlCmdSelect) buildSelectField(field interface{}) (string, []interface{}, error) {
	cmp := sql.cmp
	txtSql := ""
	args := []interface{}{}
	err := error(nil)
	aliasField := ""
	switch f := field.(type) {
	case *BoolField:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *DateTimeField:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *TextField:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[int64]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[int32]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[int16]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[int8]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[uint64]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[uint32]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[uint16]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[uint8]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[float64]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[float32]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[int]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
		//------------------------------------------
	case BoolField:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case DateTimeField:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case TextField:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[int64]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[int32]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[int16]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[int8]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[uint64]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[uint32]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[uint16]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[uint8]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[float64]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[float32]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[int]:
		aliasField = sql.getFieldName(f.UnderField)
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *exprField:
		return sql.buildExprField(f)
	case exprField:
		return sql.buildExprField(&f)

	default:
		panic(fmt.Errorf("unsupported field type %T, file orm/SqlCmdSelect.buildSelectField.go, line 101", f))

	}
	if err != nil {
		return "", nil, err
	}
	if aliasField != "" {
		txtSql = txtSql + " AS " + cmp.Quote(aliasField)
	}
	return txtSql, args, nil
}
func (sql *SqlCmdSelect) buildSelectFieldNoAlias(field interface{}) (string, []interface{}, error) {
	cmp := sql.cmp

	ret, err := cmp.Resolve(sql.tables, sql.buildContext, field, true)
	if err != nil {
		return "", nil, err
	}
	return ret.Syntax, ret.Args, nil
}
