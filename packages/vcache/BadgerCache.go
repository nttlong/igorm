package vcache

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	_ "encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
func NewBadgerCache(dbPath string, prefixKey string) (Cache, error) {

	// Đảm bảo thư mục tồn tại

	if err := os.MkdirAll(dbPath, 0755); err != nil {
		absPath, err1 := filepath.Abs(dbPath)
		if err1 != nil {
			return nil, err1
		}

		return nil, fmt.Errorf("can not open Badger DB at %s: %w", absPath, err)
	}

	opts := badger.DefaultOptions(dbPath)
	// Tùy chỉnh logger cho Badger nếu cần (để kiểm soát log output)
	// opts.Logger = nil // Tắt log của Badger nếu bạn muốn

	db, err := badger.Open(opts)
	if err != nil {
		absPath, err1 := filepath.Abs(dbPath)
		if err1 != nil {
			return nil, err1
		}
		return nil, fmt.Errorf("can not open Badger DB at %s: %w", absPath, err)
	}

	// Chạy Goroutine để dọn dẹp các mục cũ (lý tưởng là ở một Goroutine riêng)
	// Thường thì bạn sẽ gọi RunValueLogGC trong một vòng lặp định kỳ.
	// For simplicity in this example, we omit a continuous GC loop here,
	// but in a real application, you'd manage this.
	// db.RunValueLogGC(0.7) // Cần quản lý việc này liên tục

	log.Printf("BadgerCache đã mở tại: %s\n", dbPath)
	return &BadgerCache{db: db, prefixKey: prefixKey}, nil
}
func (c *BadgerCache) GetBool(ctx context.Context, key string) (bool, bool) {
	var val bool
	found := c.Get(ctx, key, &val)
	return val, found
}

// Get implements Cache.Get for BadgerCache
func (c *BadgerCache) Get(ctx context.Context, key string, dest interface{}) bool {
	val := reflect.ValueOf(dest)
	typ := reflect.TypeOf(dest)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}
	realKey := c.prefixKey + ":" + key + ":" + typ.PkgPath() + "." + typ.Name()
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
		typ := reflect.TypeOf(dest)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		fmt.Printf("gob Decode type %s error: %v\n", err, typ.String())
		return false
	}

	if err != nil {
		return false
	}

	return true
}

// Set implements Cache.Set for BadgerCache
func (c *BadgerCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {

	val := reflect.ValueOf(value)
	typ := reflect.TypeOf(value)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}
	realKey := c.prefixKey + ":" + key + ":" + typ.PkgPath() + "." + typ.Name()
	sha256Key := sha256.Sum256([]byte(realKey))
	key = string(sha256Key[:])
	// Lấy []byte từ buffer
	err := c.db.Update(func(txn *badger.Txn) error {

		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)     // Tạo Encoder
		err := encoder.Encode(val.Interface()) // Encode struct
		if err != nil {
			typ := reflect.TypeOf(value)
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem()
			}
			fmt.Printf("Gob Encode type %s error: %v\n", err, typ.String())
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
		log.Printf("BadgerCache can not write with key '%s', type for writing is '%s', error: %v\n", key, typ.String(), err)
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
	log.Println("BadgerCache is closing...")
	return c.db.Close()
}
