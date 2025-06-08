package cache

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
)

// === Triển khai BadgerCache sử dụng github.com/dgraph-io/badger/v4 ===
// Badger là một embedded key-value store, thích hợp cho cache bền vững cục bộ.
//
// BadgerCache là triển khai của Cache interface sử dụng BadgerDB
type BadgerCache struct {
	db *badger.DB
}

// NewBadgerCache tạo một instance mới của BadgerCache.
// dbPath là đường dẫn tới thư mục lưu trữ dữ liệu của Badger.
func NewBadgerCache(dbPath string) (*BadgerCache, error) {
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
	return &BadgerCache{db: db}, nil
}

// Get implements Cache.Get for BadgerCache
func (c *BadgerCache) Get(key string) (interface{}, bool) {
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
		if err != badger.ErrKeyNotFound {
			log.Printf("Lỗi khi đọc từ BadgerCache cho key '%s': %v\n", key, err)
		}
		return nil, false
	}
	// Do Badger lưu trữ bytes, bạn cần unmarshal trở lại đối tượng gốc.
	// Điều này phức tạp hơn go-cache vì go-cache lưu trữ trực tiếp interface{}.
	// Đối với bài toán login, chúng ta sẽ lưu User struct đã JSON hóa.
	// Để đơn giản, hàm Get này sẽ trả về []byte. Service cần handle việc unmarshal.
	// Hoặc bạn có thể thêm một Type specific Get (ví dụ: GetUser) vào interface nếu các loại đối tượng cache là cố định.
	// Tạm thời trả về []byte và để service xử lý.
	return valBytes, true
}

// Set implements Cache.Set for BadgerCache
func (c *BadgerCache) Set(key string, value interface{}, ttl time.Duration) {
	err := c.db.Update(func(txn *badger.Txn) error {
		// Chuyển đổi value sang []byte.
		// Đối với User struct, bạn sẽ cần JSON marshal nó.
		// Ví dụ: userBytes, _ := json.Marshal(user)
		// Ở đây, chúng ta giả định value đã là []byte hoặc có thể chuyển đổi.
		var valBytes []byte
		switch v := value.(type) {
		case []byte:
			valBytes = v
		case string:
			valBytes = []byte(v)
		case fmt.Stringer: // Nếu value có thể chuyển thành chuỗi
			valBytes = []byte(v.String())
		default:
			// Fallback: có thể marshal thành JSON nếu là struct/map
			// hoặc trả về lỗi nếu không thể chuyển đổi
			log.Printf("Cảnh báo: Không thể chuyển đổi giá trị %T thành []byte cho BadgerCache. Key: %s\n", value, key)
			return fmt.Errorf("kiểu giá trị không được hỗ trợ để lưu vào BadgerCache")
		}

		entry := badger.NewEntry([]byte(key), valBytes)
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
func (c *BadgerCache) Delete(key string) {
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
