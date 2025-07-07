package orm_test

import (
	"reflect"
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

type User struct {
	orm.Model[User] `db:"users"`
	TextField       orm.TextField `db:"len(50)"`
	BoolField       orm.BoolField
	IntField        orm.NumberField[int]
	FloatField      orm.NumberField[float64]
	Uint16Field     orm.NumberField[uint16]
	Uint32Field     orm.NumberField[uint32]
	Uint64Field     orm.NumberField[uint64]
	Int8Field       orm.NumberField[int8]
	Int16Field      orm.NumberField[int16]
	Int32Field      orm.NumberField[int32]
	Int64Field      orm.NumberField[int64]
	Float32Field    orm.NumberField[float32]
	DateTimeField   orm.DateTimeField
	NIntField       *orm.NumberField[int]
	NFloatField     *orm.NumberField[float64]
}
type UserNullable struct {
	*orm.Model[User]
	NTextField     *orm.TextField `db:"len(50)"`
	NBoolField     *orm.BoolField
	NIntField      *orm.NumberField[int]
	NFloatField    *orm.NumberField[float64]
	NUint16Field   *orm.NumberField[uint16]
	NUint32Field   *orm.NumberField[uint32]
	NUint64Field   *orm.NumberField[uint64]
	NInt8Field     *orm.NumberField[int8]
	NInt16Field    *orm.NumberField[int16]
	NInt32Field    *orm.NumberField[int32]
	NInt64Field    *orm.NumberField[int64]
	NFloat32Field  *orm.NumberField[float32]
	NDateTimeField *orm.DateTimeField
}

func TestGetMeta(t *testing.T) {
	typ := reflect.TypeOf(&User{}).Elem()
	meta := orm.Utils.GetMetaInfo(typ)
	t.Log(meta)
}
func TestQueryable(t *testing.T) {
	typ := reflect.TypeOf(&User{}).Elem()
	tblName := orm.Utils.TableNameFromStruct(typ)

	retVal := orm.EntityUtils.QueryableFromType(typ, tblName, nil)
	ret := retVal.Interface()
	qr := ret.(*User)

	t.Log(qr)

}
func TestQueryableNullField(t *testing.T) {
	typ := reflect.TypeOf(&UserNullable{}).Elem()
	tblName := orm.Utils.TableNameFromStruct(typ)

	retVal2 := orm.EntityUtils.QueryableFromType(typ, tblName, nil)
	ret2 := retVal2.Interface()

	qr2 := ret2.(*UserNullable)
	expr := qr2.NTextField.Eq("test")
	cmp := orm.Compiler.Ctx(mssql())
	r1, err := cmp.Resolve(nil, expr)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_text_field] = ?", r1.Syntax)
	assert.Equal(t, []interface{}{"test"}, r1.Args)
	// c := qr2.NIntField.Eq(1)
	//c1 :=
	expr2 := qr2.NBoolField.Eq(true).And(qr2.NIntField.Eq(1))
	cmp = orm.Compiler.Ctx(mssql())
	r2, err := cmp.Resolve(nil, expr2)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_bool_field] = ? AND [user_nullables].[n_int_field] = ?", r2.Syntax)
	assert.Equal(t, []interface{}{true, 1}, r2.Args)
	expr3 := qr2.NIntField.Eq(1).And(qr2.NFloatField.Eq(1.1))
	cmp = orm.Compiler.Ctx(mssql())
	r3, err := cmp.Resolve(nil, expr3)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_int_field] = ? AND [user_nullables].[n_float_field] = ?", r3.Syntax)
	assert.Equal(t, []interface{}{1, 1.1}, r3.Args)
	expr4 := qr2.NIntField.Eq(1).And(qr2.NFloatField.Eq(1.1)).And(qr2.NInt8Field.Eq(1))
	cmp = orm.Compiler.Ctx(mssql())
	r4, err := cmp.Resolve(nil, expr4)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_int_field] = ? AND [user_nullables].[n_float_field] = ? AND [user_nullables].[n_int8_field] = ?", r4.Syntax)
	assert.Equal(t, []interface{}{1, 1.1, 1}, r4.Args)
	expr5 := qr2.NIntField.Eq(1).And(qr2.NFloatField.Eq(1.1)).And(qr2.NInt8Field.Eq(1)).And(qr2.NInt16Field.Eq(1))

	r5, err := cmp.Resolve(nil, expr5)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_int_field] = ? AND [user_nullables].[n_float_field] = ? AND [user_nullables].[n_int8_field] = ? AND [user_nullables].[n_int16_field] = ?", r5.Syntax)
	assert.Equal(t, []interface{}{1, 1.1, 1, 1}, r5.Args)
	expr6 := qr2.NIntField.Eq(1).And(qr2.NFloatField.Eq(1.1)).And(qr2.NInt8Field.Eq(1)).And(qr2.NInt16Field.Eq(1)).And(qr2.NInt32Field.Eq(1))
	r6, err := cmp.Resolve(nil, expr6)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_int_field] = ? AND [user_nullables].[n_float_field] = ? AND [user_nullables].[n_int8_field] = ? AND [user_nullables].[n_int16_field] = ? AND [user_nullables].[n_int32_field] = ?", r6.Syntax)
	assert.Equal(t, []interface{}{1, 1.1, 1, 1, 1}, r6.Args)
	expr7 := qr2.NIntField.Eq(1).And(qr2.NFloatField.Eq(1.1)).And(qr2.NInt8Field.Eq(1)).And(qr2.NInt16Field.Eq(1)).And(qr2.NInt32Field.Eq(1)).And(qr2.NInt64Field.Eq(1))
	r7, err := cmp.Resolve(nil, expr7)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_int_field] = ? AND [user_nullables].[n_float_field] = ? AND [user_nullables].[n_int8_field] = ? AND [user_nullables].[n_int16_field] = ? AND [user_nullables].[n_int32_field] = ? AND [user_nullables].[n_int64_field] = ?", r7.Syntax)
	assert.Equal(t, []interface{}{1, 1.1, 1, 1, 1, 1}, r7.Args)
	expr8 := qr2.NIntField.Eq(1).And(qr2.NFloatField.Eq(1.1)).And(qr2.NInt8Field.Eq(1)).And(qr2.NInt16Field.Eq(1)).And(qr2.NInt32Field.Eq(1)).And(qr2.NInt64Field.Eq(1)).And(qr2.NFloat32Field.Eq(1.1))
	r8, err := cmp.Resolve(nil, expr8)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_int_field] = ? AND [user_nullables].[n_float_field] = ? AND [user_nullables].[n_int8_field] = ? AND [user_nullables].[n_int16_field] = ? AND [user_nullables].[n_int32_field] = ? AND [user_nullables].[n_int64_field] = ? AND [user_nullables].[n_float32_field] = ?", r8.Syntax)
	assert.Equal(t, []interface{}{1, 1.1, 1, 1, 1, 1, 1.1}, r8.Args)
	tn := orm.Now()
	expr9 := qr2.NIntField.Eq(1).And(qr2.NFloatField.Eq(1.1)).And(qr2.NInt8Field.Eq(1)).And(qr2.NInt16Field.Eq(1)).And(qr2.NInt32Field.Eq(1)).And(qr2.NInt64Field.Eq(1)).And(qr2.NFloat32Field.Eq(1.1)).And(qr2.NDateTimeField.Eq(tn))
	r9, err := cmp.Resolve(nil, expr9)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_int_field] = ? AND [user_nullables].[n_float_field] = ? AND [user_nullables].[n_int8_field] = ? AND [user_nullables].[n_int16_field] = ? AND [user_nullables].[n_int32_field] = ? AND [user_nullables].[n_int64_field] = ? AND [user_nullables].[n_float32_field] = ? AND [user_nullables].[n_date_time_field] = ?", r9.Syntax)
	assert.Equal(t, []interface{}{1, 1.1, 1, 1, 1, 1, 1.1, tn}, r9.Args)
	expr10 := qr2.NIntField.Eq(1).And(qr2.NFloatField.Eq(1.1)).And(qr2.NInt8Field.Eq(1)).And(qr2.NInt16Field.Eq(1)).And(qr2.NInt32Field.Eq(1)).And(qr2.NInt64Field.Eq(1)).And(qr2.NFloat32Field.Eq(1.1)).And(qr2.NDateTimeField.Eq(tn)).And(qr2.NTextField.Eq("test"))
	r10, err := cmp.Resolve(nil, expr10)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_int_field] = ? AND [user_nullables].[n_float_field] = ? AND [user_nullables].[n_int8_field] = ? AND [user_nullables].[n_int16_field] = ? AND [user_nullables].[n_int32_field] = ? AND [user_nullables].[n_int64_field] = ? AND [user_nullables].[n_float32_field] = ? AND [user_nullables].[n_date_time_field] = ? AND [user_nullables].[n_text_field] = ?", r10.Syntax)
	assert.Equal(t, []interface{}{1, 1.1, 1, 1, 1, 1, 1.1, tn, "test"}, r10.Args)

}
func TestQueryableNullFieldWithCall(t *testing.T) {
	typ := reflect.TypeOf(&UserNullable{}).Elem()
	tblName := orm.Utils.TableNameFromStruct(typ)
	retVal2 := orm.EntityUtils.QueryableFromType(typ, tblName, nil)
	ret2 := retVal2.Interface()
	cmp := orm.Compiler.Ctx(mssql())
	qr2 := ret2.(*UserNullable)
	expr := qr2.NTextField.Len().Eq(5)
	r1, err := cmp.Resolve(nil, expr)
	assert.NoError(t, err)
	assert.Equal(t, "LEN([user_nullables].[n_text_field]) = ?", r1.Syntax)
	assert.Equal(t, []interface{}{5}, r1.Args)
	expr2 := qr2.NIntField.IsNull().Or(qr2.NIntField.Eq(1))
	r2, err := cmp.Resolve(nil, expr2)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_int_field] IS NULL OR [user_nullables].[n_int_field] = ?", r2.Syntax)
	assert.Equal(t, []interface{}{1}, r2.Args)
	expr3 := qr2.NIntField.IsNull().Or(qr2.NIntField.Eq(1)).And(qr2.NFloatField.Eq(1.1))
	r3, err := cmp.Resolve(nil, expr3)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_int_field] IS NULL OR [user_nullables].[n_int_field] = ? AND [user_nullables].[n_float_field] = ?", r3.Syntax)
	assert.Equal(t, []interface{}{1, 1.1}, r3.Args)
	expr4 := qr2.NIntField.IsNull().Or(qr2.NIntField.Eq(1)).And(qr2.NFloatField.Eq(1.1)).And(qr2.NInt8Field.Eq(1))
	r4, err := cmp.Resolve(nil, expr4)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_int_field] IS NULL OR [user_nullables].[n_int_field] = ? AND [user_nullables].[n_float_field] = ? AND [user_nullables].[n_int8_field] = ?", r4.Syntax)
	assert.Equal(t, []interface{}{1, 1.1, 1}, r4.Args)
	expr5 := qr2.NIntField.IsNull().Or(qr2.NIntField.Eq(1)).And(qr2.NFloatField.Eq(1.1)).And(qr2.NInt8Field.Eq(1)).And(qr2.NInt16Field.Eq(1))
	r5, err := cmp.Resolve(nil, expr5)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_int_field] IS NULL OR [user_nullables].[n_int_field] = ? AND [user_nullables].[n_float_field] = ? AND [user_nullables].[n_int8_field] = ? AND [user_nullables].[n_int16_field] = ?", r5.Syntax)
	assert.Equal(t, []interface{}{1, 1.1, 1, 1}, r5.Args)
	expr6 := qr2.NIntField.IsNull().Or(qr2.NIntField.Eq(1)).And(qr2.NFloatField.Eq(1.1)).And(qr2.NInt8Field.Eq(1)).And(qr2.NInt16Field.Eq(1)).And(qr2.NInt32Field.Eq(1))
	r6, err := cmp.Resolve(nil, expr6)
	assert.NoError(t, err)
	assert.Equal(t, "[user_nullables].[n_int_field] IS NULL OR [user_nullables].[n_int_field] = ? AND [user_nullables].[n_float_field] = ? AND [user_nullables].[n_int8_field] = ? AND [user_nullables].[n_int16_field] = ? AND [user_nullables].[n_int32_field] = ?", r6.Syntax)
	assert.Equal(t, []interface{}{1, 1.1, 1, 1, 1}, r6.Args)

}
func TestSumOfLen(t *testing.T) {
	typ := reflect.TypeOf(&UserNullable{}).Elem()
	tblName := orm.Utils.TableNameFromStruct(typ)
	retVal2 := orm.EntityUtils.QueryableFromType(typ, tblName, nil)
	ret2 := retVal2.Interface()
	cmp := orm.Compiler.Ctx(mssql())
	qr2 := ret2.(*UserNullable)
	expr := qr2.NTextField.Len().Sum()
	r1, err := cmp.Resolve(nil, expr)
	assert.NoError(t, err)
	assert.Equal(t, "SUM(LEN([user_nullables].[n_text_field]))", r1.Syntax)
	assert.Equal(t, []interface{}{}, r1.Args)
}
