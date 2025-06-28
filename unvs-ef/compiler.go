package unvsef

import (
	"fmt"
	"reflect"
)

type sqlCompiler struct {
}
type compilerInfo struct {
	BinaryField *BinaryField
	FuncField   *FuncField
	AliasField  *AliasField
	DbField     *DbField
	Op          string
	Left        interface{}
	Right       interface{}
	Args        []interface{}
}

func (c *sqlCompiler) extract(expr interface{}) *compilerInfo {
	fieldNames := []string{
		"BinaryField",
		"FuncField",
		"AliasField",
		"DbField",
		//"Op",
	}
	ret := compilerInfo{}
	if expr == nil {
		return nil
	}
	getter := reflect.ValueOf(expr)
	getterType := getter.Type()
	if getterType.Kind() == reflect.Ptr {
		getterType = getterType.Elem()
	}
	if getter.Kind() == reflect.Ptr {
		getter = getter.Elem()
	}
	setter := reflect.ValueOf(&ret).Elem()
	hasSetField := false
	for _, fieldName := range fieldNames {
		if _, ok := getterType.FieldByName(fieldName); ok {
			getterField := getter.FieldByName(fieldName)
			setterField := setter.FieldByName(fieldName)

			if !getterField.IsValid() || !setterField.IsValid() {
				continue
			}

			// Nếu là con trỏ thì deref để gán đúng
			if getterField.Kind() == reflect.Ptr && !getterField.IsNil() {
				setterField.Set(getterField)

				hasSetField = true
			} else if getterField.Kind() != reflect.Ptr {
				setterField.Set(getterField)
				hasSetField = true
			}
		}

	}
	if hasSetField {
		return &ret
	} else {
		return nil
	}

}
func (c *sqlCompiler) exprToSQL(v interface{}, d Dialect) (string, []interface{}) {
	val := reflect.ValueOf(v)

	// typ := reflect.TypeOf(v)
	// fmt.Println(typ.Name())
	if val.Kind() == reflect.Struct {
		fmt.Println(reflect.TypeOf(v).String())

		method := val.MethodByName("ToSqlExpr2")
		if method.IsValid() && method.Type().NumIn() == 1 {
			// fmt.Println(method.Type().NumIn())
			res := method.Call([]reflect.Value{reflect.ValueOf(d)})
			if len(res) == 2 {
				if sqlStr, ok := res[0].Interface().(string); ok {
					if args, ok := res[1].Interface().([]interface{}); ok {
						return sqlStr, args
					}
				}
			}
		}

	}

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		fmt.Println(reflect.TypeOf(v).String())
		method := val.MethodByName("ToSqlExpr")
		if method.IsValid() && method.Type().NumIn() == 1 {
			res := method.Call([]reflect.Value{reflect.ValueOf(d)})
			if len(res) == 2 {
				if sqlStr, ok := res[0].Interface().(string); ok {
					if args, ok := res[1].Interface().([]interface{}); ok {
						return sqlStr, args
					}
				}
			}
		}
	}
	return "?", []interface{}{v}
}
func (c *sqlCompiler) compileBinaryField(bf *BinaryField, d Dialect) (string, []interface{}) {
	if bf.Op == "IS NULL" {
		leftExpr, leftArgs := c.exprToSQL(bf.Left, d)
		sql := fmt.Sprintf("(%s %s )", leftExpr, bf.Op)
		return sql, leftArgs
	}
	if bf.Left == nil {
		rightExpr, rightArgs := c.exprToSQL(bf.Right, d)
		sql := fmt.Sprintf("(%s %s)", bf.Op, rightExpr)
		return sql, rightArgs
	}

	leftExpr, leftArgs := c.exprToSQL(bf.Left, d)
	rightExpr, rightArgs := c.exprToSQL(bf.Right, d)
	if bf.Op == "BETWEEN" {
		typOfRightArgs := reflect.ValueOf(rightArgs)
		if typOfRightArgs.Kind() == reflect.Slice {

			leftArgs = []interface{}{typOfRightArgs.Index(0).Elem().Index(0).Interface()}
			rightArgs = []interface{}{typOfRightArgs.Index(0).Elem().Index(1).Interface()}
			sql := fmt.Sprintf("(%s %s %s)", leftExpr, "BETWEEN ? AND ?", rightExpr)
			args := append(leftArgs, rightArgs...)
			return sql, args
		}

	}
	// Nếu cả trái và phải đều là cột (không có args), không truyền tham số
	sql := fmt.Sprintf("(%s %s %s)", leftExpr, bf.Op, rightExpr)
	args := append(leftArgs, rightArgs...)
	return sql, args
}
func (c *sqlCompiler) CompileFuncField(expr *FuncField, d Dialect) (string, []interface{}) {
	args := make([]string, len(expr.Args))
	params := []interface{}{}
	for i, a := range expr.Args {
		expr, p := c.exprToSQL(a, d)
		args[i] = expr
		params = append(params, p...)
	}
	sql := fmt.Sprintf("%s(%s)", expr.FuncName, utils.join(args, ", "))
	return sql, params
}
func (c *sqlCompiler) Compile(expr interface{}, d Dialect) (string, []interface{}) {

	f := c.extract(expr)
	if f == nil {
		if bf, ok := expr.(*BinaryField); ok {
			return c.compileBinaryField(bf, d)

		}
		if ff, ok := expr.(*FuncField); ok {
			return c.CompileFuncField(ff, d)

		}
		typ := reflect.TypeOf(expr)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		tableName := utils.TableNameFromStruct(typ)

		return d.QuoteIdent(tableName), nil

	}

	if f.BinaryField != nil {

		return c.compileBinaryField(f.BinaryField, d)
	}
	if f.FuncField != nil {
		return c.CompileFuncField(f.FuncField, d)
	}
	if f.AliasField != nil {
		sql, args := c.exprToSQL(f.AliasField.Field, d)
		return fmt.Sprintf("%s AS %s", sql, d.QuoteIdent(f.AliasField.Alias)), args
	}
	if f.DbField != nil {
		return d.QuoteIdent(f.DbField.TableName, f.DbField.ColName), nil
	}

	// if f.ValueField != nil {
	// 	return "?", []interface{}{f.ValueField.Value}
	// }
	return "NULL", nil
}

var compiler = &sqlCompiler{}
