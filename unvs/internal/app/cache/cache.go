// internal/app/cache/cache.go
package cache

import (
	// Mặc dù không sử dụng trực tiếp trong go-cache/badger đơn giản, nhưng tốt cho tương lai (ví dụ: Redis client)

	"time"
	// Để dùng BadgerCache
	// Đổi tên import để tránh xung đột với package "cache" này
)

// Cache interface định nghĩa các phương thức mà bất kỳ triển khai cache nào cũng phải có.
type Cache interface {
	Get(key string) (interface{}, bool)                   // Lấy giá trị từ cache
	Set(key string, value interface{}, ttl time.Duration) // Đặt giá trị vào cache với TTL
	Delete(key string)                                    // Xóa một key khỏi cache
	Close() error                                         // Đóng kết nối/giải phóng tài nguyên của cache
}

// === Triển khai InMemoryCache sử dụng github.com/patrickmn/go-cache ===

// InMemoryCache là triển khai của Cache interface sử dụng go-cache
