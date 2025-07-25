package vcache

import (
	// Cần cho các thao tác Redis
	"context"
	"crypto/sha256" // Để hash key
	"encoding/hex"  // Để chuyển hash thành chuỗi hex
	"errors"
	"strings"

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
	client interface{}

	prefixKey string // Tiền tố key
	timeOut   time.Duration
	isCluster bool // true nếu là cluster, false nếu là single instance
}

// NewRedisCache tạo một instance mới của RedisCache.
// addr là địa chỉ của Redis server (ví dụ: "localhost:6379").
// password là mật khẩu Redis, db là số database (0-15).
func NewRedisCache(
	//ctx context.Context,

	addr string, password string,
	prefixKey string,
	db int,
	timeOut time.Duration) Cache {

	// tach addr thành các phần tử
	addrs := strings.Split(addr, ",")
	if len(addrs) == 1 {
		rdb := redis.NewClient(&redis.Options{
			Addr:     addrs[0],
			Password: password, // no password set
			DB:       db,       // use default DB
		})
		return &RedisCache{client: rdb, prefixKey: prefixKey}
	} else {
		rdb := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			Password: password, // no password set

		})
		return &RedisCache{client: rdb, prefixKey: prefixKey, isCluster: true}
	}

}

// Set đặt giá trị vào cache với TTL.
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	redisCtx, cancel := context.WithTimeout(ctx, r.timeOut)
	defer cancel()
	typ := reflect.TypeOf(value)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	key = typ.PkgPath() + ":" + typ.Name() + ":" + key
	realKey := fmt.Sprintf("%s:%s", r.prefixKey, key)
	hRealKey := sha256.Sum256([]byte(realKey))
	hashedKey := hex.EncodeToString(hRealKey[:])

	// Serialize giá trị thành JSON []byte
	byteValue, err := bytesEncodeObject(value)
	if err != nil {
		fmt.Printf("Lỗi khi serialize giá trị sang JSON cho key '%s': %v\n", key, err)
		return
	}
	if r.isCluster {
		client, ok := r.client.(*redis.ClusterClient)
		if !ok {
			fmt.Printf("Lỗi khi ép kiểu client sang ClusterClient cho key '%s'\n", key)
			return
		}
		err = client.Set(redisCtx, hashedKey, byteValue, ttl).Err()
	} else {
		client, ok := r.client.(*redis.Client)
		if !ok {
			fmt.Printf("Lỗi khi ép kiểu client sang Client cho key '%s'\n", key)
			return
		}
		err = client.Set(redisCtx, hashedKey, byteValue, ttl).Err()
	}
	// Đặt giá trị vào Redis với TTL

}

// Get lấy giá trị từ cache.
// dest phải là một CON TRỎ đến biến mà bạn muốn nhận dữ liệu.
func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) bool {
	redisCtx, cancel := context.WithTimeout(ctx, r.timeOut)
	defer cancel() // Luôn gọi cancel để giải phóng tài nguyên context
	typ := reflect.TypeOf(dest)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	key = typ.PkgPath() + ":" + typ.Name() + ":" + key
	realKey := fmt.Sprintf("%s:%s", r.prefixKey, key)
	hRealKey := sha256.Sum256([]byte(realKey))
	hashedKey := hex.EncodeToString(hRealKey[:])
	var val []byte
	if r.isCluster {
		client, ok := r.client.(*redis.ClusterClient)
		if !ok {
			return false
		}
		_val, err := client.Get(redisCtx, hashedKey).Bytes()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				// Lỗi timeout: Redis phản hồi quá chậm
				return false
			} else if err == redis.Nil {
				// Key không tồn tại trong Redis, cần lấy từ Database và lưu lại vào cache
				return false
			} else {
				return false
			}
		}
		val = _val
	} else {
		client, ok := r.client.(*redis.Client)
		if !ok {
			return false
		}
		_val, err := client.Get(redisCtx, hashedKey).Bytes()
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				// Lỗi timeout: Redis phản hồi quá chậm
				return false
			} else if err == redis.Nil {
				// Key không tồn tại trong Redis, cần lấy từ Database và lưu lại vào cache
				return false
			} else {
				return false
			}
		}
		val = _val
	}

	// Lấy giá trị từ Redis

	// Deserialize JSON []byte vào dest
	err := bytesDecodeObject(val, dest)
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
	if r.isCluster {
		client, ok := r.client.(*redis.ClusterClient)
		if !ok {
			fmt.Printf("Lỗi khi ép kiểu client sang ClusterClient cho key '%s'\n", key)
			return
		}
		err := client.Del(ctx, hashedKey).Err()
		if err != nil {
			fmt.Printf("Lỗi khi xóa dữ liệu từ Redis cho key '%s': %v\n", key, err)
		}
	} else {
		client, ok := r.client.(*redis.Client)
		if !ok {
			fmt.Printf("Lỗi khi ép kiểu client sang Client cho key '%s'\n", key)
			return
		}
		err := client.Del(ctx, hashedKey).Err()
		if err != nil {
			fmt.Printf("Lỗi khi xóa dữ liệu từ Redis cho key '%s': %v\n", key, err)
		}
	}

}

// Close đóng kết nối/giải phóng tài nguyên của cache.
func (r *RedisCache) Close() error {
	if r.isCluster {
		client, ok := r.client.(*redis.ClusterClient)
		if !ok {
			return errors.New("Lỗi khi ép kiểu client sang ClusterClient")
		}
		return client.Close()
	} else {
		client, ok := r.client.(*redis.Client)
		if !ok {
			return errors.New("Lỗi khi ép kiểu client sang Client")
		}
		return client.Close()
	}

}
