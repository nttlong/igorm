package orm

import (
	"fmt"
	"reflect"
	"unsafe"
	internal "unvs-orm/internal"
)

type Model[T any] struct {

	// dataPointer unsafe.Pointer
	TableName string
}

func verifyModelFieldFirst[T any]() {
	typ := reflect.TypeOf((*T)(nil)).Elem()
	if typ.NumField() == 0 || typ.Field(0).Name != "Model" {
		panic(fmt.Sprintf("orm.Model must be the first field in struct %s", typ.Name()))
	}
}
func (e *Model[T]) New() *Object[T] {
	verifyModelFieldFirst[T]()
	var t T

	valE := reflect.ValueOf(&t).Elem()
	// testData := reflect.New(reflect.TypeFor[T]()).Elem()

	meta := internal.Utils.GetMetaInfo(reflect.TypeFor[T]())

	for tableName, fieldMeta := range meta {

		for fieldName, field := range fieldMeta {
			f := valE.FieldByName(field.Field.Name)

			if !f.IsValid() || !f.CanAddr() {
				continue // hoặc panic nếu bạn muốn strict hơn
			}

			dbField := &dbField{
				Name:  fieldName,
				Table: tableName,
				field: field.Field,
			}

			utilsObjectIns.AssignDbFieldSmart(f, dbField) // <-- auto fallback
		}
	}

	return &Object[T]{

		Data: t,
	}
}

//go:linkname memmove runtime.memmove
func memmove(dst, src unsafe.Pointer, n uintptr)
func memset(ptr unsafe.Pointer, val byte, size uintptr) {
	b := (*[1 << 30]byte)(ptr)[:size:size]
	for i := range b {
		b[i] = val
	}
}
func (e *Model[T]) Clone(from T) T {
	var newT T
	size := unsafe.Sizeof(newT)

	src := unsafe.Pointer(&from)
	dst := unsafe.Pointer(&newT)

	// copy toàn bộ memory vùng struct
	// shallow copy (nếu có pointer, sẽ trỏ cùng)
	// đảm bảo struct không chứa slice/map cần deep copy
	memmove(dst, src, size)

	// Gán lại Model pointer cho bản clone
	ptr := (*uintptr)(dst)
	*ptr = uintptr(unsafe.Pointer(e))

	return newT
}
func (e *Model[T]) Reset(t *T) {
	size := unsafe.Sizeof(*t)
	ptr := unsafe.Pointer(t)
	memset(ptr, 0, size)

	// gán lại model sau khi reset
	*(*uintptr)(ptr) = uintptr(unsafe.Pointer(e))
}
func Queryable[T any](tenantDb *internal.TenantDb) *T {
	var v T
	typ := reflect.TypeOf(v)
	if typ == nil {
		typ = reflect.TypeOf((*T)(nil)).Elem()
	}

	e := internal.EntityUtils.QueryableFromType(typ, internal.Utils.TableNameFromStruct(typ), nil)
	ret := e.Interface().(*T)
	return ret
}
