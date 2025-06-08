package cache

import (
	// Cần cho các thao tác Redis
	"context"
	"crypto/sha256" // Để hash key
	"encoding/hex"  // Để chuyển hash thành chuỗi hex

	// Để serialize/deserialize object
	"fmt"
	"reflect" // Để lấy thông tin về package path và type name
	"time"

	"github.com/go-redis/redis/v8" // Redis client
)

// Context cho Redis operations. Trong ứng dụng thực tế, nên dùng context được truyền vào từ request.

// Cache interface
// RedisCache là một triển khai của Cache sử dụng Redis.
type RedisCache struct {
	client    *redis.Client
	prefixKey string // Tiền tố key
}

// NewRedisCache tạo một instance mới của RedisCache.
// addr là địa chỉ của Redis server (ví dụ: "localhost:6379").
// password là mật khẩu Redis, db là số database (0-15).
func NewRedisCache(
	ctx context.Context,
	ownerType reflect.Type,
	addr, password string, db int) Cache {
	prefixKey := fmt.Sprintf("%s:%s", ownerType.PkgPath(), ownerType.Name())
	hKey := sha256.Sum256([]byte(prefixKey))
	prefixKey = hex.EncodeToString(hKey[:])
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	// Ping để kiểm tra kết nối
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("Lỗi khi kết nối đến Redis: %v\n", err)
		// Trong ứng dụng thực tế, bạn có thể muốn panic hoặc trả về error ở đây.
	}
	return &RedisCache{client: rdb, prefixKey: prefixKey}
}

// Set đặt giá trị vào cache với TTL.
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	realKey := fmt.Sprintf("%s:%s", r.prefixKey, key)
	hRealKey := sha256.Sum256([]byte(realKey))
	hashedKey := hex.EncodeToString(hRealKey[:])

	// Serialize giá trị thành JSON []byte
	byteValue, err := bytesEncodeObject(value)
	if err != nil {
		fmt.Printf("Lỗi khi serialize giá trị sang JSON cho key '%s': %v\n", key, err)
		return
	}

	// Đặt giá trị vào Redis với TTL
	err = r.client.Set(ctx, hashedKey, byteValue, ttl).Err()
	if err != nil {
		fmt.Printf("Lỗi khi đặt dữ liệu vào Redis cho key '%s': %v\n", key, err)
	}
}

// Get lấy giá trị từ cache.
// dest phải là một CON TRỎ đến biến mà bạn muốn nhận dữ liệu.
func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) bool {
	realKey := fmt.Sprintf("%s:%s", r.prefixKey, key)
	hRealKey := sha256.Sum256([]byte(realKey))
	hashedKey := hex.EncodeToString(hRealKey[:])

	// Lấy giá trị từ Redis
	val, err := r.client.Get(ctx, hashedKey).Bytes()
	if err != nil {
		if err == redis.Nil { // Key không tồn tại
			return false
		}
		fmt.Printf("Lỗi khi lấy dữ liệu từ Redis cho key '%s': %v\n", key, err)
		return false
	}

	// Deserialize JSON []byte vào dest
	err = bytesDecodeObject(val, dest)
	if err != nil {
		fmt.Printf("Lỗi khi deserialize dữ liệu từ JSON cho key '%s': %v\n", key, err)
		return false
	}
	return true
}

// Delete xóa một key khỏi cache.
func (r *RedisCache) Delete(ctx context.Context, key string) {
	realKey := fmt.Sprintf("%s:%s", r.prefixKey, key)
	hRealKey := sha256.Sum256([]byte(realKey))
	hashedKey := hex.EncodeToString(hRealKey[:])

	// Xóa key khỏi Redis
	err := r.client.Del(ctx, hashedKey).Err()
	if err != nil {
		fmt.Printf("Lỗi khi xóa dữ liệu từ Redis cho key '%s': %v\n", key, err)
	}

}

// Close đóng kết nối/giải phóng tài nguyên của cache.
func (r *RedisCache) Close() error {
	return r.client.Close()
}
