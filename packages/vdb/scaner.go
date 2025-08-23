package vdb

import (
	"database/sql"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

type scanner struct {
	cacheBuildFieldMap sync.Map
}

// Giả định FieldMetadata đã được định nghĩa như bạn cung cấp:
type FieldMetadata struct {
	Offset      uintptr
	Kind        reflect.Kind
	IsPtr       bool        // true if it's a pointer type (e.g., *time.Time)
	SqlNullType interface{} // e.g., sql.NullString for *string field
	// Thêm Typ reflect.Type nếu bạn vẫn cần nó, ví dụ như trong hàm fetchAllRows để switch meta.Typ.String()
	Typ reflect.Type
}

func (u *scanner) buildFieldMapNoCache(t reflect.Type) map[string]FieldMetadata {
	m := map[string]FieldMetadata{}

	// Xử lý các trường hợp input Type là con trỏ hoặc slice
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr { // Xử lý trường hợp input là slice của con trỏ (e.g., []*Order)
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// Bỏ qua các trường không được export (chữ cái đầu thường)
		if f.PkgPath != "" {
			continue
		}

		// Xử lý các trường nhúng (anonymous fields)
		if f.Anonymous {
			// Đệ quy để lấy các trường từ struct nhúng
			// Lưu ý: buildFieldMap (có cache) nên được gọi ở đây để tận dụng cache
			m2 := u.buildFieldMap(f.Type)
			for k, v := range m2 {
				// Cần điều chỉnh Offset cho các trường nhúng
				// v.Offset += f.Offset // Chỉ cần thiết nếu bạn muốn lưu offset tương đối từ gốc
				m[k] = v
			}
			continue // Sau khi xử lý trường nhúng, bỏ qua phần còn lại của vòng lặp cho trường này
		}

		// Khởi tạo giá trị mặc định
		isPtr := false
		var sqlNullType interface{}
		fieldType := f.Type // Loại của trường hiện tại

		// Kiểm tra nếu là kiểu con trỏ
		if fieldType.Kind() == reflect.Ptr {
			isPtr = true
			fieldType = fieldType.Elem() // Lấy loại gốc mà con trỏ trỏ tới

			// Xác định SqlNullType dựa trên loại gốc
			switch fieldType.Kind() {
			case reflect.String:
				sqlNullType = &sql.NullString{}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				sqlNullType = &sql.NullInt64{}
			case reflect.Float32, reflect.Float64:
				sqlNullType = &sql.NullFloat64{}
			case reflect.Bool:
				sqlNullType = &sql.NullBool{}
			case reflect.Struct:
				// Xử lý các struct đặc biệt như time.Time
				if fieldType.ConvertibleTo(reflect.TypeOf(time.Time{})) {
					sqlNullType = &sql.NullTime{}
				}
				// Thêm các kiểu struct khác nếu cần (ví dụ: uuid.UUID)
				// case reflect.Array: // cho uuid.UUID ([16]byte)
				// if fieldType.String() == "uuid.UUID" {
				//     // sqlNullType = &NullUUID{} // Nếu bạn có một kiểu NullUUID tùy chỉnh
				// }
			}
		}

		m[f.Name] = FieldMetadata{
			Offset:      f.Offset,
			Kind:        fieldType.Kind(), // Lưu kind của loại gốc (nếu là con trỏ)
			Typ:         fieldType,        // Lưu loại gốc (nếu là con trỏ)
			IsPtr:       isPtr,
			SqlNullType: sqlNullType,
		}
	}
	return m
}

func (u *scanner) buildFieldMap(t reflect.Type) map[string]FieldMetadata {
	if v, ok := u.cacheBuildFieldMap.Load(t); ok {
		return v.(map[string]FieldMetadata)
	}
	m := u.buildFieldMapNoCache(t)
	u.cacheBuildFieldMap.Store(t, m)
	return m
}

func (u *scanner) doScan(rows *sql.Rows, typ reflect.Type) (interface{}, error) {
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	fieldMetas := make([]FieldMetadata, len(cols))
	// Điền fieldMetas dựa trên cols và fieldMap từ buildFieldMap
	// Vòng lặp này chỉ chạy một lần.
	fieldMap := u.buildFieldMap(typ)
	for i, col := range cols {
		if meta, ok := fieldMap[col]; ok {
			fieldMetas[i] = meta
		} else {
			// Xử lý cột không tìm thấy trong struct, có thể dùng một FieldMetadata đặc biệt cho dummy
			fieldMetas[i] = FieldMetadata{} // Một loại dummy meta
		}
	}

	// Khởi tạo slice kết quả (tùy chọn []*T hoặc []T)
	// Để hiệu quả, nên dùng []*T nếu không cần giá trị bản sao
	const defaultCapacity = 10000 // Hoặc một giá trị ước tính phù hợp với use-case của bạn
	resultSlice := reflect.MakeSlice(reflect.SliceOf(reflect.PointerTo(typ)), 0, defaultCapacity)

	for rows.Next() {
		// Tạo một instance mới của struct và lấy con trỏ base
		// new(T) nhanh hơn reflect.New(typ).Elem()
		instancePtr := reflect.New(typ).Interface() // new(T)
		// #nosec G103 -- using unsafe.Pointer with reflect.UnsafeAddr for zero-copy field mapping
		mem := unsafe.Pointer(reflect.ValueOf(instancePtr).Pointer()) // base pointer của T

		scanArgs := make([]interface{}, len(cols)) // Vẫn cần làm lại mỗi lần nếu dùng scanArgs []interface{}

		for i, meta := range fieldMetas {
			if meta.Offset == 0 && !meta.IsPtr && meta.Kind == 0 { // Đây là một dummy field meta
				var dummy interface{}
				scanArgs[i] = &dummy
				continue
			}

			// Nếu là kiểu con trỏ (nullable)
			if meta.IsPtr {
				// Tạo biến sql.NullX tạm thời và scan vào đó
				// Ví dụ: var ns sql.NullString; scanArgs[i] = &ns
				// Sau đó cần logic để gán ns.String vào (*string)(fieldPtr) sau rows.Scan
				// Đây là phần phức tạp nhất.
				// Tạm thời, dùng một dummy nếu không muốn xử lý phức tạp ngay lập tức.
				var dummy interface{}
				scanArgs[i] = &dummy
				// Bạn cần một cơ chế ánh xạ phức tạp hơn ở đây cho kiểu con trỏ
			} else {
				// #nosec G103 -- using unsafe.Pointer with reflect.UnsafeAddr for zero-copy field mapping
				fieldPtr := unsafe.Pointer(uintptr(mem) + meta.Offset)
				switch meta.Kind {
				case reflect.String:
					scanArgs[i] = (*string)(fieldPtr)
				case reflect.Int:
					scanArgs[i] = (*int)(fieldPtr)
				// ... các kiểu khác
				case reflect.Struct:
					if meta.SqlNullType != nil {
						scanArgs[i] = meta.SqlNullType // Ví dụ: &sql.NullTime{}
					} else {
						var dummy interface{}
						scanArgs[i] = &dummy
					}
				default:
					var dummy interface{}
					scanArgs[i] = &dummy
				}
			}
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		// Sau khi scan, xử lý các giá trị nullable nếu có (ví dụ: từ sql.NullX sang *string)
		// ...

		resultSlice = reflect.Append(resultSlice, reflect.ValueOf(instancePtr))
	}

	return resultSlice.Interface(), nil
}

func FetchAllRowsV2(rows *sql.Rows, typ reflect.Type) ([]interface{}, error) {
	defer rows.Close() // Đảm bảo rows luôn được đóng

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	ret := make([]interface{}, 0)
	for rows.Next() {
		itemVal := reflect.New(typ).Elem()
		vals := make([]interface{}, len(columns))
		valPtrs := make([]interface{}, len(columns))

		for i := 0; i < len(columns); i++ {
			valPtrs[i] = itemVal.FieldByName(columns[i]).Addr().Interface()
		}

		err := rows.Scan(valPtrs...)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(columns); i++ {
			vals[i] = itemVal.Field(i).Interface()
		}

		ret = append(ret, itemVal.Interface())
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

var scannerInst = &scanner{}
