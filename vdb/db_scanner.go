package vdb

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unsafe"
)

type fieldOffset struct {
	offset uintptr
	typ    reflect.Type
}

type scanPlan struct {
	elemSize uintptr
	columns  []string
	offsets  []fieldOffset
	newSlice func() reflect.Value       // return reflect.Value of []T
	setters  func(unsafe.Pointer) []any // return []interface{} to pass to rows.Scan
}

var scanPlanCache sync.Map // key = reflect.Type + cols

func getScanPlan(t reflect.Type, cols []string) *scanPlan {
	key := t.PkgPath() + "." + t.Name() + ":" + strings.Join(cols, ",")
	if cached, ok := scanPlanCache.Load(key); ok {
		return cached.(*scanPlan)
	}

	var offsets []fieldOffset
	for _, col := range cols {
		field, ok := t.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, col)
		})
		if ok && field.IsExported() {
			offsets = append(offsets, fieldOffset{
				offset: field.Offset,
				typ:    field.Type,
			})
		} else {
			offsets = append(offsets, fieldOffset{}) // dummy
		}
	}

	plan := &scanPlan{
		elemSize: t.Size(),
		columns:  cols,
		offsets:  offsets,
		newSlice: func() reflect.Value {
			slice := reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
			return slice
		},
		setters: func(ptr unsafe.Pointer) []any {
			args := make([]any, len(cols))
			for i, off := range offsets {
				if off.typ != nil {
					fieldPtr := unsafe.Pointer(uintptr(ptr) + off.offset)
					args[i] = reflect.NewAt(off.typ, fieldPtr).Interface()
				} else {
					var dummy any
					args[i] = &dummy
				}
			}
			return args
		},
	}

	scanPlanCache.Store(key, plan)
	return plan
}
func ScanToStructUnsafeCached[T any](rows *sql.Rows, out *[]T) error {
	defer rows.Close()
	t := reflect.TypeOf((*T)(nil)).Elem()
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	plan := getScanPlan(t, cols)
	sliceVal := plan.newSlice()

	for rows.Next() {
		elem := reflect.New(t).Elem()
		ptr := unsafe.Pointer(elem.UnsafeAddr())
		args := plan.setters(ptr)

		if err := rows.Scan(args...); err != nil {
			return err
		}
		sliceVal = reflect.Append(sliceVal, elem)
	}

	*out = sliceVal.Interface().([]T)
	return rows.Err()
}

type structFieldInfo struct {
	Offset uintptr
	Type   reflect.Type
}

var structMapCache sync.Map // map[string][]structFieldInfo

/*
La ham nhanh nhat
*/
func ScanToStructValueCached[T any](rows *sql.Rows) ([]T, error) {
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var t T
	tType := reflect.TypeOf(t)

	if tType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("T must be struct, got: %s", tType.String())
	}

	cacheKey := tType.PkgPath() + "." + tType.Name() + "|" + strings.Join(cols, ",")

	var fieldInfos []structFieldInfo
	if cached, ok := structMapCache.Load(cacheKey); ok {
		fieldInfos = cached.([]structFieldInfo)
	} else {
		fieldMap := make(map[string]structFieldInfo)
		for i := 0; i < tType.NumField(); i++ {
			f := tType.Field(i)
			dbName := f.Tag.Get("db")
			if dbName == "" {
				dbName = strings.ToLower(f.Name)
			}
			fieldMap[dbName] = structFieldInfo{
				Offset: f.Offset,
				Type:   f.Type,
			}
		}

		fieldInfos = make([]structFieldInfo, len(cols))
		for i, col := range cols {
			if fi, ok := fieldMap[col]; ok {
				fieldInfos[i] = fi
			} else {
				fieldInfos[i] = structFieldInfo{
					Offset: 0,
					Type:   reflect.TypeOf(new(interface{})).Elem(),
				}
			}
		}
		structMapCache.Store(cacheKey, fieldInfos)
	}

	results := make([]T, 0, 512)
	args := make([]interface{}, len(cols))

	for rows.Next() {
		item := reflect.New(tType).Elem() // Struct value (not pointer)
		ptr := unsafe.Pointer(item.UnsafeAddr())

		for i, fi := range fieldInfos {
			fieldPtr := unsafe.Add(ptr, fi.Offset)
			args[i] = reflect.NewAt(fi.Type, fieldPtr).Interface()
		}

		if err := rows.Scan(args...); err != nil {
			return nil, err
		}

		results = append(results, item.Interface().(T))
	}

	return results, nil
}

// structPlan là cache ánh xạ cột -> offset trong struct
type structPlan struct {
	columnToOffset map[string]uintptr
	structType     reflect.Type
}

var (
	planCache    sync.Map // map[string]*structPlan
	ptrSlicePool = sync.Pool{
		New: func() any {
			return make([]interface{}, 0, 20)
		},
	}
)

// getPlan: lấy ánh xạ field offset cho struct
func getPlan[T any](columns []string) (*structPlan, error) {
	var t T
	tType := reflect.TypeOf(t)
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
	}
	cacheKey := tType.PkgPath() + "." + tType.Name() + "|" + fmt.Sprint(columns)

	if cached, ok := planCache.Load(cacheKey); ok {
		return cached.(*structPlan), nil
	}

	fieldOffsets := map[string]uintptr{}
	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		name := field.Tag.Get("db")
		if name == "" {
			name = field.Name
		}
		fieldOffsets[name] = field.Offset
	}

	plan := &structPlan{
		columnToOffset: fieldOffsets,
		structType:     tType,
	}
	planCache.Store(cacheKey, plan)
	return plan, nil
}
func ScanToStructValueCachedFix[T any](rows *sql.Rows) ([]T, error) {
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var t T
	tType := reflect.TypeOf(t)
	if tType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("T must be struct, got: %s", tType.String())
	}

	cacheKey := tType.PkgPath() + "." + tType.Name() + "|" + strings.Join(cols, ",")

	var fieldInfos []structFieldInfo
	if cached, ok := structMapCache.Load(cacheKey); ok {
		fieldInfos = cached.([]structFieldInfo)
	} else {
		fieldMap := map[string]structFieldInfo{}
		for i := 0; i < tType.NumField(); i++ {
			f := tType.Field(i)
			dbName := f.Tag.Get("db")
			if dbName == "" {
				dbName = strings.ToLower(f.Name)
			}
			fieldMap[dbName] = structFieldInfo{
				Offset: f.Offset,
				Type:   f.Type,
			}
		}

		for _, col := range cols {
			if fi, ok := fieldMap[col]; ok {
				fieldInfos = append(fieldInfos, fi)
			} else {
				fieldInfos = append(fieldInfos, structFieldInfo{
					Offset: 0,
					Type:   reflect.TypeOf(new(interface{})).Elem(),
				})
			}
		}
		structMapCache.Store(cacheKey, fieldInfos)
	}

	results := make([]T, 0, 512)
	args := make([]interface{}, len(cols))
	fieldHolders := make([]reflect.Value, len(cols))

	for rows.Next() {
		item := reflect.New(tType).Elem()
		basePtr := unsafe.Pointer(item.UnsafeAddr())

		for i, fi := range fieldInfos {
			// Tạo holder riêng biệt cho từng field (pointer đến đúng kiểu)
			holder := reflect.New(fi.Type)
			fieldHolders[i] = holder
			args[i] = holder.Interface()
		}

		if err := rows.Scan(args...); err != nil {
			return nil, err
		}

		// Sau khi scan, gán lại vào field thật trong item
		for i, fi := range fieldInfos {
			fieldPtr := unsafe.Add(basePtr, fi.Offset)
			dest := reflect.NewAt(fi.Type, fieldPtr).Elem()
			dest.Set(fieldHolders[i].Elem())
		}

		results = append(results, item.Interface().(T))
	}

	return results, nil
}

// unsafeSetField: đặt giá trị bằng unsafe
func unsafeSetField(ptr unsafe.Pointer, offset uintptr, val interface{}) {
	fieldPtr := unsafe.Pointer(uintptr(ptr) + offset)

	switch v := val.(type) {
	case *int64:
		*(*int64)(fieldPtr) = *v
	case *int:
		*(*int)(fieldPtr) = *v
	case *string:
		*(*string)(fieldPtr) = *v
	case *float64:
		*(*float64)(fieldPtr) = *v
	case *[]byte:
		*(*[]byte)(fieldPtr) = *v
	default:
		// fallback qua reflect cho các kiểu chưa xử lý
		reflect.NewAt(reflect.TypeOf(val).Elem(), fieldPtr).Elem().Set(reflect.ValueOf(val).Elem())
	}
}
func ScanToStructUnsafeCachedImproveV2[T any](rows *sql.Rows) ([]T, error) {
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	plan, err := getPlan[T](cols)
	if err != nil {
		return nil, err
	}

	var results []T
	var tZero T
	tType := reflect.TypeOf(tZero)
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
	}

	for rows.Next() {
		var t T
		tPtr := reflect.ValueOf(&t).Elem() // reflect.Value của struct

		ptrs := make([]interface{}, len(cols)) // không dùng sync.Pool vội

		for i, col := range cols {
			if offset, ok := plan.columnToOffset[col]; ok {
				field := reflect.NewAt(tType.FieldByIndex([]int{0}).Type, unsafe.Pointer(uintptr(tPtr.UnsafeAddr())+offset)).Interface()
				ptrs[i] = field
			} else {
				var dummy interface{}
				ptrs[i] = &dummy
			}
		}

		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		results = append(results, t)
	}

	return results, nil
}
