package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"
	internal "unvs-orm/internal"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/microsoft/go-mssqldb"
)

var EntityUtils = internal.EntityUtils
var mapType = map[reflect.Type]string{
	reflect.TypeOf(NumberField[int64]{}):   "int64",
	reflect.TypeOf(NumberField[int]{}):     "int",
	reflect.TypeOf(NumberField[float32]{}): "float32",
	reflect.TypeOf(NumberField[float64]{}): "float64",
	reflect.TypeOf(TextField{}):            "string",
	reflect.TypeOf(BoolField{}):            "bool",

	// reflect.TypeOf(FieldUUID{}):     "uuid.UUID",
	reflect.TypeOf(NumberField[int16]{}):  "int16",
	reflect.TypeOf(NumberField[int32]{}):  "int32",
	reflect.TypeOf(NumberField[int64]{}):  "int64",
	reflect.TypeOf(NumberField[uint]{}):   "uint",
	reflect.TypeOf(NumberField[uint16]{}): "uint16",
	reflect.TypeOf(NumberField[uint32]{}): "uint32",
	reflect.TypeOf(NumberField[uint64]{}): "uint64",
	reflect.TypeOf(DateTimeField{}):       "time.Time",
}
var entityTypeName = strings.Split(reflect.TypeOf(Model[any]{}).String(), "[")[0] + "["
var currentPackagePath = reflect.TypeOf(DateTimeField{}).PkgPath()

func init() {
	internal.InitUtils(currentPackagePath, entityTypeName, mapType)
	internal.EntityUtils.FieldResolver = resolveFieldType
}

func resolveFieldType(tableName string, colName string, field reflect.StructField) reflect.Value {
	fieldType := field.Type

	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	if fieldType == reflect.TypeOf(TextField{}) {
		ret := TextField{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal

	}
	if fieldType == reflect.TypeOf(BoolField{}) {
		ret := BoolField{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}
	if fieldType == reflect.TypeOf(DateTimeField{}) {
		ret := DateTimeField{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}
	if fieldType == reflect.TypeOf(NumberField[int64]{}) {
		ret := NumberField[int64]{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}
	if fieldType == reflect.TypeOf(NumberField[int]{}) {
		ret := NumberField[int]{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}
	if fieldType == reflect.TypeOf(NumberField[float32]{}) {
		ret := NumberField[float32]{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}
	if fieldType == reflect.TypeOf(NumberField[float64]{}) {
		ret := NumberField[float64]{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}
	if fieldType == reflect.TypeOf(NumberField[int8]{}) {
		ret := NumberField[int8]{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}
	if fieldType == reflect.TypeOf(NumberField[int16]{}) {
		ret := NumberField[int16]{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}
	if fieldType == reflect.TypeOf(NumberField[int32]{}) {
		ret := NumberField[int32]{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}
	if fieldType == reflect.TypeOf(NumberField[uint]{}) {
		ret := NumberField[uint]{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}
	if fieldType == reflect.TypeOf(NumberField[uint16]{}) {
		ret := NumberField[uint16]{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}
	if fieldType == reflect.TypeOf(NumberField[uint32]{}) {
		ret := NumberField[uint32]{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}
	if fieldType == reflect.TypeOf(NumberField[uint64]{}) {
		ret := NumberField[uint64]{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}
	if fieldType == reflect.TypeOf(NumberField[uint8]{}) {
		ret := NumberField[uint8]{
			dbField: &dbField{
				field: field,
				Name:  colName,
				Table: tableName,
			},
		}
		retVal := reflect.ValueOf(&ret).Elem()
		return retVal
	}

	panic(fmt.Errorf("not implemented resolveFieldType for %s, source utils.go", fieldType.String()))
}

var Utils = internal.Utils

func Now() time.Time {
	return time.Now()
}

type TenantDb = internal.TenantDb

func Open(driverName, dataSourceName string) (*sql.DB, error) {
	return sql.Open(driverName, dataSourceName)
}

func NewTenantDb(db *sql.DB) (*TenantDb, error) {
	return internal.Utils.NewTenantDb(db)

}

type utilsObject struct {
	useSafeAssign          bool
	enableFallbackDebugLog bool
	didFallbackOnce        sync.Once
}

var utilsObjectIns = &utilsObject{
	useSafeAssign:          true, // bật để dùng switch-case
	enableFallbackDebugLog: true, // bật log khi fallback
}

// Switch-case fallback
func (u *utilsObject) Assign(fieldPtr any, dbf *dbField) {
	switch v := fieldPtr.(type) {
	case *TextField:
		v.dbField = dbf
	case *NumberField[int]:
		v.dbField = dbf
	case *NumberField[int8]:
		v.dbField = dbf
	case *NumberField[int16]:
		v.dbField = dbf
	case *NumberField[int32]:
		v.dbField = dbf
	case *NumberField[int64]:
		v.dbField = dbf
	case *NumberField[uint8]:
		v.dbField = dbf
	case *NumberField[uint16]:
		v.dbField = dbf
	case *NumberField[uint32]:
		v.dbField = dbf
	case *NumberField[uint64]:
		v.dbField = dbf
	case *NumberField[float32]:
		v.dbField = dbf
	case *NumberField[float64]:
		v.dbField = dbf
	case *BoolField:
		v.dbField = dbf
	case *DateTimeField:
		v.dbField = dbf
	default:
		panic(fmt.Sprintf("Unsupported field type: %T", v))
	}
}

// Gán dbField bằng unsafe + fallback an toàn
func (u *utilsObject) AssignDbFieldSmart(f reflect.Value, dbf *dbField) {
	if !f.IsValid() || !f.CanAddr() {
		panic("AssignDbFieldSmart: field is not addressable or invalid")
	}

	defer func() {
		if r := recover(); r != nil {
			u.didFallbackOnce.Do(func() {
				if u.enableFallbackDebugLog {
					fmt.Printf("[⚠️ Fallback] AssignDbFieldSmart panic: %v. Switching to Assign()\n", r)
				}
			})
			// fallback
			fieldPtr := f.Addr().Interface()
			u.Assign(fieldPtr, dbf)
		}
	}()

	if u.useSafeAssign {
		fieldPtr := f.Addr().Interface()
		u.Assign(fieldPtr, dbf)
		return
	}

	fieldPtr := unsafe.Pointer(f.UnsafeAddr())
	*(*uintptr)(fieldPtr) = uintptr(unsafe.Pointer(dbf))
}

func Ptr[T any](v T) *T {
	return &v
}
