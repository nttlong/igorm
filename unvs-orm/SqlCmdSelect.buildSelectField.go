package orm

import (
	"fmt"
)

func (sql *SqlCmdSelect) buildSelectField(field interface{}) (string, []interface{}, error) {
	cmp := sql.cmp
	txtSql := ""
	args := []interface{}{}
	err := error(nil)
	aliasField := ""
	switch f := field.(type) {
	case *BoolField:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *DateTimeField:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *TextField:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[int64]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[int32]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[int16]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[int8]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[uint64]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[uint32]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[uint16]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[uint8]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[float64]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[float32]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case *NumberField[int]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
		//------------------------------------------
	case BoolField:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case DateTimeField:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case TextField:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[int64]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[int32]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[int16]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[int8]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[uint64]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[uint32]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[uint16]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[uint8]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[float64]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[float32]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)
	case NumberField[int]:
		aliasField = f.field.Name
		txtSql, args, err = sql.buildSelectFieldNoAlias(f)

	default:
		panic(fmt.Errorf("unsupported field type %T, file SqlCmdSelect.buildSelectField.go, line 14", f))

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
	ret, err := cmp.Resolve(sql.aliasMap, field)
	if err != nil {
		return "", nil, err
	}
	return ret.Syntax, ret.Args, nil
}
