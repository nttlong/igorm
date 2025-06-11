package caching

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	_ "encoding/gob"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/dgraph-io/badger/v4"
)

// === Triển khai BadgerCache sử dụng github.com/dgraph-io/badger/v4 ===
// Badger là một embedded key-value store, thích hợp cho cache bền vững cục bộ.
//
// BadgerCache là triển khai của Cache interface sử dụng BadgerDB
type BadgerCache struct {
	db        *badger.DB
	prefixKey string
}

// NewBadgerCache tạo một instance mới của BadgerCache.
// dbPath là đường dẫn tới thư mục lưu trữ dữ liệu của Badger.
func NewBadgerCache(ownerType reflect.Type, dbPath string) (Cache, error) {
	prefixType := ownerType.PkgPath() + "." + ownerType.Name()
	h := sha256.Sum256([]byte(prefixType))
	prefixType = string(h[:])
	// Đảm bảo thư mục tồn tại

	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return nil, fmt.Errorf("không thể tạo thư mục cho Badger DB tại %s: %w", dbPath, err)
	}

	opts := badger.DefaultOptions(dbPath)
	// Tùy chỉnh logger cho Badger nếu cần (để kiểm soát log output)
	// opts.Logger = nil // Tắt log của Badger nếu bạn muốn

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("không thể mở Badger DB tại %s: %w", dbPath, err)
	}

	// Chạy Goroutine để dọn dẹp các mục cũ (lý tưởng là ở một Goroutine riêng)
	// Thường thì bạn sẽ gọi RunValueLogGC trong một vòng lặp định kỳ.
	// For simplicity in this example, we omit a continuous GC loop here,
	// but in a real application, you'd manage this.
	// db.RunValueLogGC(0.7) // Cần quản lý việc này liên tục

	log.Printf("BadgerCache đã mở tại: %s\n", dbPath)
	return &BadgerCache{db: db, prefixKey: prefixType}, nil
}

// Get implements Cache.Get for BadgerCache
func (c *BadgerCache) Get(ctx context.Context, key string, dest interface{}) bool {
	realKey := c.prefixKey + key
	sha256Key := sha256.Sum256([]byte(realKey))
	key = string(sha256Key[:])
	var valBytes []byte
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err // badger.ErrKeyNotFound hoặc lỗi khác
		}
		valBytes, err = item.ValueCopy(nil)
		return err
	})

	if err != nil {
		return false
	}
	// Do Badger lưu trữ bytes, bạn cần unmarshal trở lại đối tượng gốc.
	// Điều này phức tạp hơn go-cache vì go-cache lưu trữ trực tiếp interface{}.
	// Đối với bài toán login, chúng ta sẽ lưu User struct đã JSON hóa.
	// Để đơn giản, hàm Get này sẽ trả về []byte. Service cần handle việc unmarshal.
	// Hoặc bạn có thể thêm một Type specific Get (ví dụ: GetUser) vào interface nếu các loại đối tượng cache là cố định.
	// Tạm thời trả về []byte và để service xử lý.
	decoder := gob.NewDecoder(bytes.NewBuffer(valBytes)) // Tạo Decoder từ []byte
	err = decoder.Decode(dest)                           // Decode vào biến đích
	if err != nil {
		fmt.Printf("Lỗi khi Gob Decode userBytes vào user2: %v\n", err)
		return false
	}

	if err != nil {
		return false
	}

	return true
}

// Set implements Cache.Set for BadgerCache
func (c *BadgerCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
	realKey := c.prefixKey + key
	sha256Key := sha256.Sum256([]byte(realKey))
	key = string(sha256Key[:])
	// Lấy []byte từ buffer
	err := c.db.Update(func(txn *badger.Txn) error {

		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer) // Tạo Encoder
		err := encoder.Encode(value)       // Encode struct
		if err != nil {
			fmt.Printf("Lỗi khi Gob Encode user1: %v\n", err)
			return err
		}
		userBytes := buffer.Bytes()

		entry := badger.NewEntry([]byte(key), userBytes)
		if ttl > 0 {
			entry = entry.WithTTL(ttl)
		}
		return txn.SetEntry(entry)
	})
	if err != nil {
		log.Printf("Lỗi khi ghi vào BadgerCache cho key '%s': %v\n", key, err)
	}
}

// Delete implements Cache.Delete for BadgerCache
func (c *BadgerCache) Delete(ctx context.Context, key string) {
	realKey := c.prefixKey + key
	sha256Key := sha256.Sum256([]byte(realKey))
	key = string(sha256Key[:])
	err := c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	if err != nil {
		log.Printf("Lỗi khi xóa từ BadgerCache cho key '%s': %v\n", key, err)
	}
}

// Close implements Cache.Close for BadgerCache
func (c *BadgerCache) Close() error {
	log.Println("BadgerCache đang đóng...")
	return c.db.Close()
}
