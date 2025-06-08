package dbx

import (
	"context"
	"fmt"
	"strings"
)

// Implement UpdateWithContext (đã thấy trong code của bạn)
// Đây là hàm sẽ được gọi nội bộ bởi QrBuilder.Update
// Hàm này sẽ nhận một map các trường cần cập nhật
func updateWithContext[T any](ctx context.Context, client *DBXTenant, tableName string, updates map[string]interface{}, conditions []string, args []interface{}) error {
	// Mô phỏng việc cập nhật
	setClauses := []string{}
	setArgs := []interface{}{}
	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", field))
		setArgs = append(setArgs, value)
	}

	query := fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(setClauses, ", "))
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	allArgs := append(setArgs, args...)

	fmt.Printf("Mô phỏng UpdateWithContext: SQL=%s ARGS=%v\n", query, allArgs)
	// Trong thực tế, bạn sẽ gọi client.ExecContext(ctx, query, allArgs...)
	return nil
}
