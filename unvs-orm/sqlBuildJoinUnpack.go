package orm

import (
	"fmt"
	"reflect"
)

type joinUnpack struct {
	typeCheck []reflect.Type
}
type joinUnpackResult struct {
	aliasMap map[string]string
	tables   []string
	on       *BoolField
}
type joinInfoFromBoolField struct {
	alias  map[string]string
	tables []string
}

func (j *joinUnpack) unpack(on *BoolField) *JoinExpr {
	panic("not implemented") // TODO: Implement)

}
func (j *joinUnpack) check(typ reflect.Type) bool {
	for _, t := range j.typeCheck {
		if t == typ {
			return true
		}
	}
	return false
}
func (j *joinUnpack) extractJoinInfos(refTable *joinRefInfo, expr ...interface{}) *joinRefInfo {

	for _, field := range expr {
		if field == nil {
			continue
		}
		typ := reflect.TypeOf(field)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if j.check(typ) {
			continue
		}
		switch v := field.(type) {
		case *BoolField:
			if bf, ok := v.UnderField.(*fieldBinary); ok {
				refTable = j.extractJoinInfos(refTable, bf.left, bf.right)
				if refTable.hasNewTable {
					bf.op = refTable.joinType + " JOIN"
					// reset all
					refTable.hasNewTable = false
					refTable.newTableName = ""
					refTable.newTableNameAlias = ""
				}
			}

		case NumberField[int]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case *NumberField[int]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)

		case NumberField[uint64]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case NumberField[uint32]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case NumberField[uint16]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case NumberField[uint8]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case NumberField[int64]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case NumberField[int32]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case NumberField[int16]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case NumberField[int8]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)

		case *NumberField[uint64]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case *NumberField[uint32]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case *NumberField[uint16]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case *NumberField[uint8]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case *NumberField[int64]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case *NumberField[int32]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case *NumberField[int16]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case *NumberField[int8]:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case DateTimeField:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case *DateTimeField:
			refTable = j.extractJoinInfos(refTable, v.UnderField)
		case *methodCall:
			if v == nil {
				continue
			}
			refTable = j.extractJoinInfos(refTable, v.args)
		case *fieldBinary:
			refTable = j.extractJoinInfos(refTable, v.left, v.right)
		case fieldBinary:
			refTable = j.extractJoinInfos(refTable, v.left, v.right)
		case *dbField:
			if v == nil {
				continue
			}
			if _, ok := refTable.alias[v.Table]; !ok {
				refTable.alias[v.Table] = fmt.Sprintf("T%d", len(refTable.tables)+1)
				refTable.tables = append(refTable.tables, v.Table)
				refTable.hasNewTable = true
				refTable.newTableName = v.Table
				refTable.newTableNameAlias = refTable.alias[v.Table]
			}

			return refTable

		default:
			panic(fmt.Sprintf("unsupported type %T", v))
		}

	}
	return refTable

}

type joinRefInfo struct {
	alias             map[string]string
	tables            []string
	hasNewTable       bool
	joinType          string
	newTableName      string
	newTableNameAlias string
}

func (j *joinUnpack) ExtractJoinInfo(on *BoolField) *joinInfoFromBoolField {
	if on.UnderField == nil {
		return nil
	}
	if fx, ok := on.UnderField.(*joinField); ok {
		refTables := &joinRefInfo{
			alias:       map[string]string{},
			tables:      []string{},
			hasNewTable: false,
			joinType:    fx.joinType,
		}
		refTables = j.extractJoinInfos(refTables, on)
		ret := &joinInfoFromBoolField{
			alias:  map[string]string{},
			tables: []string{},
		}
		retTables := []string{}
		tblIndex := 1
		for _, table := range refTables.tables {
			if _, ok := ret.alias[table]; ok {
				continue
			}
			ret.alias[table] = fmt.Sprintf("T%d", tblIndex)
			tblIndex++
			retTables = append(retTables, table)
		}
		ret.tables = retTables
		return ret
	}
	panic(fmt.Sprintf("unsupported type %T in orm/sqlBuildJoinUnpack.go, line 100", on.UnderField))

}

var joinUnpackUtils = joinUnpack{
	typeCheck: []reflect.Type{
		reflect.TypeOf(int(0)),
		reflect.TypeOf(uint(0)),
		reflect.TypeOf(int64(0)),
		reflect.TypeOf(int32(0)),
		reflect.TypeOf(int16(0)),
		reflect.TypeOf(int8(0)),
		reflect.TypeOf(uint64(0)),
		reflect.TypeOf(uint32(0)),
		reflect.TypeOf(uint16(0)),
		reflect.TypeOf(uint8(0)),
		reflect.TypeOf(float32(0)),
		reflect.TypeOf(float64(0)),
		reflect.TypeOf(""),
		reflect.TypeOf(bool(false)),
	},
}
