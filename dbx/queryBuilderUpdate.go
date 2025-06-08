package dbx

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// Update cập nhật các trường của entity T vào cơ sở dữ liệu.
// Chỉ những trường có giá trị "được set" (bao gồm cả con trỏ không nil) sẽ được cập nhật.
// Các điều kiện WHERE đã được thêm vào QrBuilder trước đó sẽ được áp dụng.
func (q QrBuilder[T]) Update(entity T) error {
	// Lấy tên bảng từ kiểu T
	var zero T
	tableName := reflect.TypeOf(zero).Name() // Giả định tên bảng là tên struct

	// Lấy giá trị reflect của entity
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Nếu entity là con trỏ, lấy giá trị mà nó trỏ tới
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("entity must be a struct or a pointer to a struct, got %s", val.Kind())
	}

	// Xây dựng danh sách các trường cần cập nhật và giá trị của chúng
	updates := make(map[string]interface{})
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Lấy tên cột từ tag `db` hoặc tên trường
		// Giả định bạn có tag `db:"column_name"` hoặc dùng tên trường trực tiếp
		columnName := field.Name
		if dbTag := field.Tag.Get("db"); dbTag != "" {
			columnName = dbTag
		}

		// Kiểm tra xem trường có "được set" hay không
		if fieldVal.Kind() == reflect.Ptr {
			// Nếu là con trỏ, chỉ cập nhật nếu con trỏ KHÔNG phải là nil
			if !fieldVal.IsNil() {
				// Lấy giá trị mà con trỏ trỏ tới
				updates[columnName] = fieldVal.Elem().Interface()
			}
		} else {
			// Nếu không phải con trỏ, luôn cập nhật giá trị này
			// Giả định rằng nếu nó có mặt trong struct, nó đã được "set"
			updates[columnName] = fieldVal.Interface()
		}
	}

	if len(updates) == 0 {
		return errors.New("no fields to update in the provided entity")
	}

	// Gọi hàm UpdateWithContext nội bộ của dbx để thực thi
	// Cần context. Nếu QrBuilder không có context, phải thêm vào chữ ký Update
	// Hiện tại, tôi giả định QrBuilder đã có ctx từ lúc khởi tạo.
	if q.ctx == nil {
		q.ctx = context.Background() // Fallback nếu QrBuilder không có context
	}

	// dbx.UpdateWithContext sẽ lấy tableName, updates, conditions, args để tạo SQL
	return updateWithContext[T](q.ctx, q.dbx, tableName, updates, []string{q.where}, q.args)
}
