package vcache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// MemcachedCache là một triển khai của Cache sử dụng Memcached.
type MemcachedCache struct {
	client    *memcache.Client
	prefixKey string // prefix cho tên key để đảm bảo tính toàn vẹn của key
}

// NewMemcachedCache tạo một instance mới của MemcachedCache.
// servers là danh sách các địa chỉ server Memcached (ví dụ: "127.0.0.1:11211").
var (
	mcClient *memcache.Client
	once     sync.Once
)

func NewMemcachedCache(ownerType reflect.Type, servers []string) Cache {
	once.Do(func() {
		mcClient = memcache.New(servers...)
	})
	mc := memcache.New(servers...)
	prefixKey := fmt.Sprintf("%s:%s:", ownerType.PkgPath(), ownerType.Name())
	return &MemcachedCache{
		client:    mc,
		prefixKey: prefixKey,
	}
}

// getHashedKey tạo một hash SHA256 từ key đầu vào.
// Memcached có giới hạn về độ dài key (250 ký tự), việc hash giúp xử lý các key dài.
func getHashedKey(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Set đặt giá trị vào cache với TTL.
func (m *MemcachedCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	typ := reflect.TypeOf(value)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	key = m.prefixKey + key + ":" + typ.PkgPath() + "." + typ.Name()
	hashedKey := getHashedKey(key)

	var byteValue []byte
	var errMarshal error

	// Cố gắng serialize value thành JSON
	byteValue, errMarshal = bytesEncodeObject(value)
	if errMarshal != nil {
		fmt.Printf("Lỗi khi serialize giá trị sang JSON cho key '%s': %v\n", key, errMarshal)
		return // Không thể lưu nếu serialize thất bại
	}

	expiration := int32(ttl.Seconds())

	err := m.client.Set(&memcache.Item{
		Key:        hashedKey,
		Value:      byteValue,
		Expiration: expiration,
	})
	if err != nil {
		fmt.Printf("Lỗi khi đặt dữ liệu vào Memcached cho key '%s': %v\n", key, err)
	}
}

// Delete xóa một key khỏi cache.
func (m *MemcachedCache) Delete(ctx context.Context, key string) {

	hashedKey := getHashedKey(key)
	err := m.client.Delete(hashedKey)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			// Key không tồn tại, không cần làm gì
			return
		}
		fmt.Printf("Lỗi khi xóa dữ liệu từ Memcached: %v\n", err)
	}
}

// Close đóng kết nối/giải phóng tài nguyên của cache.
// Đối với gomemcache, không có phương thức Close() cụ thể cho client,
// nhưng chúng ta giữ phương thức này để tuân thủ interface và cho các triển khai cache khác.
func (m *MemcachedCache) Close() error {
	// gomemcache client tự động quản lý kết nối.
	// Không có tài nguyên cụ thể nào cần giải phóng ở đây.
	return nil
}

// Get lấy giá trị từ cache.
// Hàm này sẽ TRẢ VỀ []byte ĐƯỢC LẤY TRỰC TIẾP từ Memcached.
// Người gọi sẽ phải tự deserialize []byte này thành kiểu dữ liệu mong muốn.
func (m *MemcachedCache) Get(ctx context.Context, key string, dest interface{}) bool {
	typ := reflect.TypeOf(dest)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	key = m.prefixKey + key + ":" + typ.PkgPath() + "." + typ.Name()
	hashedKey := getHashedKey(key)
	item, err := m.client.Get(hashedKey)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return false // Key không tồn tại
		}
		// Xử lý các lỗi khác nếu cần (ví dụ: log lỗi)
		fmt.Printf("Lỗi khi lấy dữ liệu từ Memcached: %v\n", err)
		return false
	}
	// Trả về []byte trực tiếp. Người gọi phải tự giải mã.
	// desrialize thành kiểu dữ liệu mong muốn.

	bff := item.Value
	err = bytesDecodeObject(bff, dest)
	if err != nil {
		fmt.Printf("Lỗi khi giải mã dữ liệu từ Memcached: %v\n", err)
		return false
	}
	return true
}
