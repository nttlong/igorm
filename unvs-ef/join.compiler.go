package unvsef

import (
	"encoding/json"
	"reflect"
)

type JoinExpr struct {
	LeftTable       string
	LeftTableAlias  string
	RightTable      string
	RightTableAlias string
	On              interface{}
	JoinType        string
	Previous        *JoinExpr
	Index           byte
}

/*
Make inner join expression from binary field
*/
func (expr *BinaryField) InnerJoin() *JoinExpr {

	lefValue := reflect.ValueOf(expr.Left)
	rightValue := reflect.ValueOf(expr.Right)
	if lefValue.Kind() == reflect.Ptr {
		lefValue = lefValue.Elem()
	}
	if rightValue.Kind() == reflect.Ptr {
		rightValue = rightValue.Elem()
	}
	leftTableNameField := lefValue.FieldByName("TableName")

	leftTable := ""
	if leftTableNameField.IsValid() {
		leftTable = leftTableNameField.String()

	}
	rightTableNameField := rightValue.FieldByName("TableName")
	rightTable := ""
	if rightTableNameField.IsValid() {
		rightTable = rightTableNameField.String()

	}
	return &JoinExpr{
		LeftTable:       leftTable,
		RightTable:      rightTable,
		LeftTableAlias:  "L1",
		RightTableAlias: "R1",
		JoinType:        "INNER JOIN",
		On:              expr,
		Index:           byte(1),
	}
}
func (J *JoinExpr) String() string {
	json, _ := json.MarshalIndent(J, "", "  ")
	return string(json)
}

/*
Continue inner join expression from previous join expression
*/
func (J *JoinExpr) InnerJoin(exprs ...*BinaryField) *JoinExpr {
	ret := &JoinExpr{

		JoinType: "INNER JOIN",

		Previous: J,
		Index:    J.Index + 1,
	}
	for i, expr := range exprs {
		leftType := reflect.TypeOf(expr.Left)
		rightType := reflect.TypeOf(expr.Right)
		lefValue := reflect.ValueOf(expr.Left)
		rightValue := reflect.ValueOf(expr.Right)
		if lefValue.Kind() == reflect.Ptr {
			lefValue = lefValue.Elem()
			leftType = lefValue.Type()
		}
		if rightValue.Kind() == reflect.Ptr {
			rightValue = rightValue.Elem()
			rightType = rightValue.Type()
		}

		leftTable := ""
		if leftType.Kind() == reflect.Struct {
			if _, ok := leftType.FieldByName("TableName"); ok {
				leftTableNameField := lefValue.FieldByName("TableName")
				if leftTableNameField.IsValid() {
					leftTable = leftTableNameField.String()
				}
			}
		}
		rightTable := ""
		if rightType.Kind() == reflect.Struct {
			if _, ok := rightType.FieldByName("TableName"); ok {
				rightTableNameField := rightValue.FieldByName("TableName")

				if rightTableNameField.IsValid() {
					rightTable = rightTableNameField.String()
				}
			}
		}
		if i == 0 {
			ret.LeftTable = leftTable
			ret.RightTable = rightTable
			ret.LeftTableAlias = "L" + string(ret.Index)
			ret.RightTableAlias = "R" + string(ret.Index)
			ret.On = expr
			ret.Previous = J

		} else {
			ret = &JoinExpr{
				LeftTable:       leftTable,
				RightTable:      rightTable,
				LeftTableAlias:  "L" + string(ret.Index),
				RightTableAlias: "R" + string(ret.Index),
				JoinType:        "INNER JOIN",
				On:              expr,
				Previous:        ret,
			}
		}

	}
	return ret
}

/*
this function get sql string of on expression and parameters
Examle: onExpr := repo.Orders.OrderId.Eq(repo.OrderItems.OrderId)
return "L1.order_id = R1.order_id", []interface{}{}
*/
func (J *JoinExpr) GetSqlJoin(d Dialect) (string, []interface{}) {
	sql := d.QuoteIdent(J.LeftTable) + " AS " + d.QuoteIdent(J.LeftTableAlias)
	args := []interface{}{}

	for current := J; current != nil; current = current.Previous {
		join := current.JoinType + " " + d.QuoteIdent(current.RightTable) + " AS " + d.QuoteIdent(current.RightTableAlias)
		if current.On != nil {
			onSQL, onArgs := compiler.ToSqlJoinClause(current.On, d)
			join += " ON " + onSQL
			args = append(args, onArgs...)
		}
		sql = join + " " + sql
	}

	return sql, args
}
