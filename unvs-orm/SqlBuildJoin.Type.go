package orm

import (
	"reflect"
	"strconv"
)

type JoinExpr struct {
	baseTable string
	previous  *JoinExpr

	on       *BoolField
	aliasMap map[string]string
	joinType string
	index    int
}

// type JoinExpr struct {
// 	joinType    string
// 	aliasSource map[string]string
// 	on          interface{} // can be nil for CROSS JOIN or delayed assignment

// 	tables []string
// }

//	func InnerJoin(source ...interface{}) *JoinExpr {
//		AliasSource := map[string]string{}
//		for i, s := range source {
//			alias := "T" + strconv.Itoa(i+1)
//			if tableName, ok := s.(string); ok {
//				AliasSource[tableName] = alias
//			} else {
//				typ := reflect.TypeOf(s)
//				if typ.Kind() == reflect.Ptr {
//					typ = typ.Elem()
//				}
//				tableName = Utils.TableNameFromStruct(typ)
//				AliasSource[tableName] = alias
//			}
//		}
//		return &JoinExpr{
//			joinType: "INNER JOIN",
//			aliasMap: AliasSource,
//			on:       nil,
//		}
//	}
func (j *JoinExpr) On(on *BoolField) *JoinExpr {
	j.on = on
	return j
}
func (m *Model[T]) Join(other interface{}, on *BoolField) *JoinExpr {
	root := &JoinExpr{
		joinType:  "INNER",
		index:     1,
		baseTable: m.TableName,
	}
	ret := &JoinExpr{
		previous: root,

		aliasMap: map[string]string{},
		on:       on,
		index:    2,
	}
	otherTyp := reflect.TypeOf(other)
	if otherTyp.Kind() == reflect.Ptr {
		otherTyp = otherTyp.Elem()
	}
	otherTableName := Utils.TableNameFromStruct(otherTyp)
	// ret.tables = append(ret.tables, m.TableName, otherTableName)
	ret.aliasMap[m.TableName] = "T1"
	ret.aliasMap[otherTableName] = "T2"

	ret.baseTable = otherTableName

	root.aliasMap = ret.aliasMap

	return ret
}

//	func (j *JoinExpr) getBaseTable() string {
//		// Tìm node gốc (previous = nil)
//		node := j
//		for node.previous != nil {
//			node = node.previous
//		}
//		return node.baseTable
//	}
func (j *JoinExpr) Join(other interface{}, on *BoolField) *JoinExpr {
	j.joinType = "INNER"
	ret := &JoinExpr{
		previous: j,
		joinType: "INNER",

		aliasMap: j.aliasMap,
		on:       on,
		index:    j.index + 1,
	}
	otherTyp := reflect.TypeOf(other)
	if otherTyp.Kind() == reflect.Ptr {
		otherTyp = otherTyp.Elem()
	}
	otherTableName := Utils.TableNameFromStruct(otherTyp)
	// ret.tables = append(ret.tables, m.TableName, otherTableName)
	t1 := "T" + strconv.Itoa(j.index+1)

	ret.aliasMap[Utils.TableNameFromStruct(otherTyp)] = t1
	ret.baseTable = otherTableName

	return ret
}
