package orm_test

import (
	"reflect"
	"testing"
	orm "unvs-orm"
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
	orm.Model[User]
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
	retVal := orm.EntityUtils.QueryableFromType(typ, tblName)
	ret := retVal.Interface()
	qr := ret.(*User)

	t.Log(qr)

}
func TestQueryableNullField(t *testing.T) {
	typ := reflect.TypeOf(&UserNullable{}).Elem()
	tblName := orm.Utils.TableNameFromStruct(typ)
	retVal2 := orm.EntityUtils.QueryableFromType(typ, tblName)
	ret2 := retVal2.Interface()

	qr2 := ret2.(*UserNullable)

	t.Log(qr2)
}
