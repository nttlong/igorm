package orm

import (
	"reflect"
	"strconv"
)

type joinExprText struct {
	Expr string
	Args []interface{}
}
type JoinExpr struct {
	*joinExprText

	baseTable string
	previous  *JoinExpr

	on       *BoolField
	aliasMap map[string]string
	joinType string
	index    int
}

func (m *Model[T]) Joins(joins ...*JoinExpr) *JoinExpr {
	root := joins[0]
	for i := 1; i < len(joins); i++ {
		root = root.doJoin(joins[i].joinType, joins[i].baseTable, joins[i].on)
	}
	return root
}
func (m *Model[T]) doJoin(joinType string, other interface{}, on *BoolField) *JoinExpr {
	root := &JoinExpr{
		joinType:  joinType,
		index:     1,
		baseTable: m.TableName,
	}
	ret := &JoinExpr{
		previous: root,
		joinType: joinType,

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
func (m *Model[T]) Join(other interface{}, on *BoolField) *JoinExpr {
	return m.doJoin("INNER", other, on)
}
func (m *Model[T]) LeftJoin(other interface{}, on *BoolField) *JoinExpr {
	return m.doJoin("LEFT", other, on)
}
func (m *Model[T]) RightJoinJoin(other interface{}, on *BoolField) *JoinExpr {
	return m.doJoin("RIGHT", other, on)
}
func (j *JoinExpr) doJoin(joinType string, other interface{}, on *BoolField) *JoinExpr {
	j.joinType = joinType
	ret := &JoinExpr{
		previous: j,
		joinType: joinType,

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
func (j *JoinExpr) Join(other interface{}, on *BoolField) *JoinExpr {
	return j.doJoin("INNER", other, on)
}
func (j *JoinExpr) LeftJoin(other interface{}, on *BoolField) *JoinExpr {
	return j.doJoin("LEFT", other, on)
}
func (j *JoinExpr) RightJoin(other interface{}, on *BoolField) *JoinExpr {
	return j.doJoin("RIGHT", other, on)
}
func (m *Model[T]) JoinBy(on *BoolField) *JoinExpr {
	return joinUnpackUtils.unpack(on)
}

type joinField struct {
	left     interface{}
	right    interface{}
	joinType string
	alias    map[string]string
	tables   []string
}

func (f *BoolField) Join(other interface{}) *BoolField {
	return &BoolField{
		underField: &joinField{
			left:     f,
			right:    other,
			joinType: "INNER",
		},
	}

}
func (f *BoolField) RightJoin(other interface{}) *BoolField {
	return &BoolField{
		underField: &joinField{
			left:     f,
			right:    other,
			joinType: "RIGHT",
		},
	}

}
func (f *BoolField) LeftJoin(other interface{}) *BoolField {
	return &BoolField{
		underField: &joinField{
			left:     f,
			right:    other,
			joinType: "LEFT",
		},
	}
}

func (f *BoolField) FullJoin(other interface{}) *BoolField {
	return &BoolField{
		underField: &joinField{
			left:     f,
			right:    other,
			joinType: "FULL",
		},
	}

}

func (f *BoolField) doJoin() *BoolField {
	if fx, ok := f.underField.(fieldBinary); ok {
		joinInfo := joinUnpackUtils.ExtractJoinInfo(f)
		return &BoolField{
			underField: &joinField{
				left:     fx.left,
				right:    fx.right,
				joinType: "INNER",
				alias:    joinInfo.alias,
				tables:   joinInfo.tables,
			},
		}

	}
	panic("Can not convert to join")

}
