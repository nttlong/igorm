package vdb

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"
	"vdb/migrate"
	"vdb/tenantDB"
)

var decoderSelectorCache sync.Map // map[reflect.Type]func([]interface{}) (any, error)

func getSelectorDecoder[T any](typ reflect.Type, cols []migrate.ColumnDef) func([]interface{}) (any, error) {
	if fn, ok := decoderSelectorCache.Load(typ); ok {
		return fn.(func([]interface{}) (any, error))
	}
	fn := func(scanData []interface{}) (any, error) {
		ptr := reflect.New(typ).Elem() // instance T

		for i, col := range cols {
			field := ptr.FieldByName(col.Field.Name)
			if !field.IsValid() || !field.CanSet() {
				continue
			}

			raw := scanData[i]
			val := reflect.ValueOf(raw)

			if !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil()) {
				continue // nil or invalid, skip
			}

			// Giải con trỏ nếu cần
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
				if !val.IsValid() {
					continue
				}
			}

			// Gán nếu type tương thích
			if val.Type().AssignableTo(field.Type()) {
				field.Set(val)
			} else if val.Type().ConvertibleTo(field.Type()) {
				field.Set(val.Convert(field.Type()))
			}
		}

		return ptr.Addr().Interface(), nil
	}

	decoderSelectorCache.Store(typ, fn)
	return fn
}

func SelectAllOriginalVersion[T any](db *tenantDB.TenantDB) ([]*T, error) {
	// 1. Khởi tạo thông tin entity
	dialect := dialectFactory.Create(db.GetDriverName())
	repoType := inserterObj.getEntityInfo(reflect.TypeFor[T]())
	tableName := repoType.tableName
	columns := repoType.entity.GetColumns()

	// 2. Tạo SELECT query
	fieldsSelect := make([]string, len(columns))
	for i, col := range columns {
		fieldsSelect[i] = dialect.Quote(col.Name) + " AS " + dialect.Quote(col.Field.Name)
	}
	sql := fmt.Sprintf("SELECT %s FROM %s", strings.Join(fieldsSelect, ", "), dialect.Quote(tableName))

	// 3. Thực thi query
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 4. Chuẩn bị thông tin reflect
	typ := reflect.TypeFor[T]()
	fieldIndexes := make([][]int, len(columns)) // cache field index paths
	fieldTypes := make([]reflect.Type, len(columns))
	for i, col := range columns {
		fieldIndexes[i] = col.IndexOfField
		fieldTypes[i] = col.Field.Type
	}

	// 5. Buffer để scan dữ liệu từ DB
	vals := make([]interface{}, len(columns))
	ptrs := make([]interface{}, len(columns))
	for i := range ptrs {
		ptrs[i] = &vals[i]
	}

	// 6. Lặp và scan từng dòng
	result := make([]*T, 0, 1024)
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		// Tạo instance T
		ptr := reflect.New(typ).Elem()

		for i := range columns {
			raw := vals[i]
			if raw == nil {
				continue
			}

			val := reflect.ValueOf(raw)
			if !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil()) {
				continue
			}
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
				if !val.IsValid() {
					continue
				}
			}

			field := ptr.FieldByIndex(fieldIndexes[i])
			if !field.CanSet() {
				continue
			}

			if val.Type().AssignableTo(fieldTypes[i]) {
				field.Set(val)
			} else if val.Type().ConvertibleTo(fieldTypes[i]) {
				field.Set(val.Convert(fieldTypes[i]))
			}
		}

		result = append(result, ptr.Addr().Interface().(*T))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

var __timeType = reflect.TypeOf(time.Time{})

func SelectAllUnsafe[T any](db *tenantDB.TenantDB) ([]*T, error) {
	dialect := dialectFactory.Create(db.GetDriverName())
	repoType := inserterObj.getEntityInfo(reflect.TypeFor[T]())
	tableName := repoType.tableName
	columns := repoType.entity.GetColumns()

	fieldsSelect := make([]string, len(columns))
	for i, col := range columns {
		fieldsSelect[i] = dialect.Quote(col.Name) + " AS " + dialect.Quote(col.Field.Name)
	}
	sql := fmt.Sprintf("SELECT %s FROM %s", strings.Join(fieldsSelect, ", "), dialect.Quote(tableName))

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// reflect type
	typ := reflect.TypeFor[T]()

	// Prepare buffer for Scan
	vals := make([]interface{}, len(columns))
	ptrs := make([]interface{}, len(columns))
	for i := range vals {
		ptrs[i] = &vals[i]
	}

	// Create unsafe field writers
	fieldSetters := make([]func(obj unsafe.Pointer, val interface{}), len(columns))
	for i, col := range columns {
		offset := col.Field.Offset
		fieldType := col.Field.Type

		fieldSetters[i] = func(offset uintptr, fieldType reflect.Type) func(unsafe.Pointer, interface{}) {
			return func(obj unsafe.Pointer, val interface{}) {
				if val == nil {
					return
				}
				fieldPtr := unsafe.Pointer(uintptr(obj) + offset)

				switch {
				case fieldType == __timeType:
					// handle []byte -> time.Time
					switch val := val.(type) {
					case time.Time:
						*(*time.Time)(fieldPtr) = val
					case []byte:
						parsed, err := time.Parse("2006-01-02 15:04:05", string(val))
						if err == nil {
							*(*time.Time)(fieldPtr) = parsed
						}
					case string:
						parsed, err := time.Parse("2006-01-02 15:04:05", val)
						if err == nil {
							*(*time.Time)(fieldPtr) = parsed
						}
					}

				case fieldType.Kind() == reflect.Int:
					*(*int)(fieldPtr) = int(reflect.ValueOf(val).Convert(fieldType).Int())
				case fieldType.Kind() == reflect.Int64:
					*(*int64)(fieldPtr) = reflect.ValueOf(val).Convert(fieldType).Int()
				case fieldType.Kind() == reflect.Int32:
					*(*int32)(fieldPtr) = int32(reflect.ValueOf(val).Convert(fieldType).Int())
				case fieldType.Kind() == reflect.Float64:
					*(*float64)(fieldPtr) = reflect.ValueOf(val).Convert(fieldType).Float()
				case fieldType.Kind() == reflect.String:
					*(*string)(fieldPtr) = reflect.ValueOf(val).Convert(fieldType).String()
				case fieldType.Kind() == reflect.Bool:
					*(*bool)(fieldPtr) = reflect.ValueOf(val).Convert(fieldType).Bool()
				default:
					// fallback chậm
					reflect.NewAt(fieldType, fieldPtr).Elem().Set(reflect.ValueOf(val).Convert(fieldType))
				}
			}
		}(offset, fieldType)
	}

	result := make([]*T, 0, 1024)
	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		obj := reflect.New(typ).Interface().(*T)
		objPtr := unsafe.Pointer(obj)

		for i, set := range fieldSetters {
			set(objPtr, vals[i]) //<--panic: reflect.Value.Convert: value of type []uint8 cannot be converted to type time.Time

		}
		result = append(result, obj)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
func SelectAllZeroAlloc[T any](db *tenantDB.TenantDB) ([]T, error) {
	dialect := dialectFactory.Create(db.GetDriverName())
	repoType := inserterObj.getEntityInfo(reflect.TypeFor[T]())
	tableName := repoType.tableName
	columns := repoType.entity.GetColumns()

	fieldsSelect := make([]string, len(columns))
	for i, col := range columns {
		fieldsSelect[i] = dialect.Quote(col.Name) + " AS " + dialect.Quote(col.Field.Name)
	}
	sql := fmt.Sprintf("SELECT %s FROM %s", strings.Join(fieldsSelect, ", "), dialect.Quote(tableName))

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//typ := reflect.TypeFor[T]()
	vals := make([]interface{}, len(columns))
	ptrs := make([]interface{}, len(columns))
	for i := range ptrs {
		ptrs[i] = &vals[i]
	}

	// Chuẩn bị fieldSetters dùng unsafe
	var timeType = reflect.TypeOf(time.Time{})
	fieldSetters := make([]func(obj unsafe.Pointer, val interface{}), len(columns))
	for i, col := range columns {
		offset := col.Field.Offset
		fieldType := col.Field.Type

		fieldSetters[i] = func(offset uintptr, fieldType reflect.Type) func(unsafe.Pointer, interface{}) {
			return func(obj unsafe.Pointer, val interface{}) {
				if val == nil {
					return
				}
				fieldPtr := unsafe.Pointer(uintptr(obj) + offset)

				switch {
				case fieldType == timeType:
					switch v := val.(type) {
					case time.Time:
						*(*time.Time)(fieldPtr) = v
					case []byte:
						t, err := time.Parse("2006-01-02 15:04:05", string(v))
						if err == nil {
							*(*time.Time)(fieldPtr) = t
						}
					case string:
						t, err := time.Parse("2006-01-02 15:04:05", v)
						if err == nil {
							*(*time.Time)(fieldPtr) = t
						}
					}
				case fieldType.Kind() == reflect.Int:
					*(*int)(fieldPtr) = int(reflect.ValueOf(val).Convert(fieldType).Int())
				case fieldType.Kind() == reflect.Int64:
					*(*int64)(fieldPtr) = reflect.ValueOf(val).Convert(fieldType).Int()
				case fieldType.Kind() == reflect.Int32:
					*(*int32)(fieldPtr) = int32(reflect.ValueOf(val).Convert(fieldType).Int())
				case fieldType.Kind() == reflect.Float64:
					*(*float64)(fieldPtr) = reflect.ValueOf(val).Convert(fieldType).Float()
				case fieldType.Kind() == reflect.String:
					*(*string)(fieldPtr) = reflect.ValueOf(val).Convert(fieldType).String()
				case fieldType.Kind() == reflect.Bool:
					*(*bool)(fieldPtr) = reflect.ValueOf(val).Convert(fieldType).Bool()
				default:
					reflect.NewAt(fieldType, fieldPtr).Elem().Set(reflect.ValueOf(val).Convert(fieldType))
				}
			}
		}(offset, fieldType)
	}

	// Preallocate slice of T
	result := make([]T, 0, 1024)

	for rows.Next() {
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		var row T
		ptr := unsafe.Pointer(&row)

		for i, set := range fieldSetters {
			set(ptr, vals[i])
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
func SelectAllZeroAllocScanDirect[T any](db *tenantDB.TenantDB) ([]T, error) {
	dialect := dialectFactory.Create(db.GetDriverName())
	repoType := inserterObj.getEntityInfo(reflect.TypeFor[T]())
	tableName := repoType.tableName
	columns := repoType.entity.GetColumns()

	fieldsSelect := make([]string, len(columns))
	for i, col := range columns {
		fieldsSelect[i] = dialect.Quote(col.Name) + " AS " + dialect.Quote(col.Field.Name)
	}
	sql := fmt.Sprintf("SELECT %s FROM %s", strings.Join(fieldsSelect, ", "), dialect.Quote(tableName))

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vals := make([]interface{}, len(columns))
	ptrs := make([]interface{}, len(columns))

	var timeType = reflect.TypeOf(time.Time{})

	// build unsafe setters
	fieldSetters := make([]func(unsafe.Pointer, interface{}), len(columns))
	for i, col := range columns {
		offset := col.Field.Offset
		fieldType := col.Field.Type

		fieldSetters[i] = func(offset uintptr, fieldType reflect.Type) func(unsafe.Pointer, interface{}) {
			return func(obj unsafe.Pointer, val interface{}) {
				if val == nil {
					return
				}
				fieldPtr := unsafe.Pointer(uintptr(obj) + offset)

				switch {
				case fieldType == timeType:
					switch v := val.(type) {
					case time.Time:
						*(*time.Time)(fieldPtr) = v
					case []byte:
						t, err := time.Parse("2006-01-02 15:04:05", string(v))
						if err == nil {
							*(*time.Time)(fieldPtr) = t
						}
					case string:
						t, err := time.Parse("2006-01-02 15:04:05", v)
						if err == nil {
							*(*time.Time)(fieldPtr) = t
						}
					}
				case fieldType.Kind() == reflect.Int:
					*(*int)(fieldPtr) = int(reflect.ValueOf(val).Convert(fieldType).Int())
				case fieldType.Kind() == reflect.Int64:
					*(*int64)(fieldPtr) = reflect.ValueOf(val).Convert(fieldType).Int()
				case fieldType.Kind() == reflect.Int32:
					*(*int32)(fieldPtr) = int32(reflect.ValueOf(val).Convert(fieldType).Int())
				case fieldType.Kind() == reflect.Float64:
					*(*float64)(fieldPtr) = reflect.ValueOf(val).Convert(fieldType).Float()
				case fieldType.Kind() == reflect.String:
					*(*string)(fieldPtr) = reflect.ValueOf(val).Convert(fieldType).String()
				case fieldType.Kind() == reflect.Bool:
					*(*bool)(fieldPtr) = reflect.ValueOf(val).Convert(fieldType).Bool()
				default:
					reflect.NewAt(fieldType, fieldPtr).Elem().Set(reflect.ValueOf(val).Convert(fieldType))
				}
			}
		}(offset, fieldType)
	}

	// Scan
	result := make([]T, 0, 1024)
	for rows.Next() {
		// tạm push zero-value để có result[i]
		result = append(result, *new(T))
		i := len(result) - 1

		for k := range ptrs {
			ptrs[k] = &vals[k]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		rowPtr := unsafe.Pointer(&result[i])
		for j, set := range fieldSetters {
			set(rowPtr, vals[j])
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func SelectAll[T any](db *tenantDB.TenantDB) ([]*T, error) {
	return SelectAllOriginalVersion[T](db)

}
