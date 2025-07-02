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
	SortField   *SortField
	Op          string
	Left        interface{}
	Right       interface{}
	Args        []interface{}
}

func (c *sqlCompiler) extract(expr interface{}) *compilerInfo {
	if typ := reflect.TypeOf(expr); typ.Kind() != reflect.Struct {
		return nil
	}
	fieldNames := []string{
		"BinaryField",
		"FuncField",
		"AliasField",
		"DbField",
		"SortField",
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
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	// typ := reflect.TypeOf(v)
	// fmt.Println(typ.Name())
	if val.Kind() == reflect.Struct {
		fmt.Println(reflect.TypeOf(v).String())

		method := val.MethodByName("ToSqlExpr")
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
func (c *sqlCompiler) compilerField(expr *DbField, d Dialect) (string, []interface{}) {
	panic("not support")
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

	if bf, ok := expr.(*FieldBool); ok {
		left, argsLeft := c.Compile(bf.Left, d)
		right, argsRight := c.Compile(bf.Right, d)
		sql := fmt.Sprintf("(%s %s %s)", left, bf.Op, right)
		return sql, append(argsLeft, argsRight...)

	}
	if bf, ok := expr.(*BinaryField); ok {
		return c.compileBinaryField(bf, d)

	}
	if bf, ok := expr.(BinaryField); ok {
		return c.compileBinaryField(&bf, d)

	}
	if ff, ok := expr.(*FuncField); ok {
		return c.CompileFuncField(ff, d)

	}
	if df, ok := expr.(*DbField); ok {
		tableName := d.QuoteIdent(df.TableName)
		colName := d.QuoteIdent(df.ColName)

		return tableName + "." + colName, nil
	}
	if df, ok := expr.(*SortField); ok {
		sqlText, args := c.Compile(df.Field, d)

		return sqlText + " " + df.Sort, args
	}
	if df, ok := expr.(*AliasField); ok {
		sqlText, args := c.Compile(df.Field, d)

		return sqlText + " AS " + d.QuoteIdent(df.Alias), args
	}
	if df, ok := expr.(AliasField); ok {
		sqlText, args := c.Compile(df.Field, d)

		return sqlText + " AS " + d.QuoteIdent(df.Alias), args
	}
	typ := reflect.TypeOf(expr)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	tableName := utils.TableNameFromStruct(typ)

	return d.QuoteIdent(tableName), nil
}
func (c *sqlCompiler) ToSqlJoinClause(expr interface{}, d Dialect) (string, []interface{}) {
	if expr == nil {
		return "", nil
	}
	if bf, ok := expr.(*BinaryField); ok {
		return c.compileBinaryField(bf, d)

	}
	panic("Not support join clause")

}

var compiler = &sqlCompiler{}
