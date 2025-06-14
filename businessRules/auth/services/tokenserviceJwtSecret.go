package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	authModels "dbmodels/auth"
	"dbx"
	"encoding/base64"
	"fmt"
	"io"
	"reflect"
	"sync"
)

// GenerateRandomSecret tạo một chuỗi ngẫu nhiên an toàn về mật mã với độ dài mong muốn.
// Chuỗi được mã hóa base64 để đảm bảo tất cả các ký tự đều in được.
func generateRandomSecret(length int) (string, error) {
	// Tính toán số byte cần thiết.
	// Mỗi ký tự base64 mã hóa 6 bit, nên 4 ký tự base64 mã hóa 3 byte.
	// Để có length ký tự base64, chúng ta cần (length * 6) / 8 = (length * 3) / 4 bytes.
	// Làm tròn lên để đảm bảo đủ ký tự.
	numBytes := (length * 3) / 4
	if (length*3)%4 != 0 {
		numBytes++ // Đảm bảo đủ byte cho mã hóa base64
	}

	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Mã hóa byte thành chuỗi Base64.
	// Sử dụng RawURLEncoding để tránh padding ('=') và các ký tự không an toàn cho URL,
	// mặc dù cho secret thì Base64StdEncoding cũng ổn.
	secret := base64.RawURLEncoding.EncodeToString(randomBytes)

	// Cắt chuỗi nếu nó dài hơn độ dài mong muốn (do làm tròn lên numBytes)
	if len(secret) > length {
		secret = secret[:length]
	}

	return secret, nil
}
func (p *TokenService) RemoveCacheJWTSecret() {
	path := p.getPath()
	key := path + "://GetJwtSecret/v2" + p.TenantDb.TenantDbName
	cacheGetJwtSecret.Delete(key)
}

var cacheGetJwtSecret sync.Map

func (p *TokenService) GetJwtSecret() ([]byte, error) {
	path := p.getPath()
	key := path + "://GetJwtSecret/v2" + p.TenantDb.TenantDbName
	if v, ok := cacheGetJwtSecret.Load(key); ok {
		return v.([]byte), nil
	}

	jwtSecret, err := p.getJwtSecret()
	if err != nil {
		return nil, err
	}
	cacheGetJwtSecret.Store(key, jwtSecret)
	return jwtSecret, nil
}

// encryptBytes mã hóa một slice bytes (data) sử dụng AES-GCM.
// keyString là khóa mã hóa (phải có độ dài 16, 24 hoặc 32 byte khi được chuyển đổi sang []byte).
// Hàm này trả về dữ liệu đã mã hóa (ciphertext kèm nonce).
func (p *TokenService) encryptBytes(keyString string, data []byte) (*string, error) {
	// Chuyển đổi khóa từ string sang []byte
	// Đảm bảo keyString có độ dài phù hợp (ví dụ: 32 ký tự cho AES-256)
	// Trong thực tế, keyString nên là một chuỗi Base64 đã được Decode,
	// hoặc một chuỗi bytes thô có độ dài cố định.
	// Nếu keyString không được tạo từ base64, hãy đảm bảo nó có đủ entropy và độ dài.
	key := []byte(keyString)

	// Kiểm tra độ dài khóa
	// Khóa phải có độ dài 16 (AES-128), 24 (AES-192) hoặc 32 (AES-256) byte
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key length for AES: %d bytes (must be 16, 24, or 32)", len(key))
	}

	// Tạo một block cipher từ khóa
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	// Tạo một Galois Counter Mode (GCM) instance từ block cipher.
	// GCM là chế độ mã hóa được khuyến nghị vì nó cung cấp cả tính bảo mật và tính toàn vẹn (authentication).
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	// Tạo một Nonce (Number Used Once - Vector Khởi tạo).
	// Nonce phải là duy nhất cho mỗi lần mã hóa với cùng một khóa.
	// Nonce không cần giữ bí mật nhưng phải được lưu trữ cùng với ciphertext để giải mã.
	// GCM.NonceSize() sẽ trả về kích thước nonce tiêu chuẩn cho GCM (thường là 12 bytes).
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Mã hóa dữ liệu.
	// gcm.Seal() sẽ nối nonce vào đầu của ciphertext trả về.
	// Tham số 'additionalData' (tham số thứ tư) được đặt là nil ở đây,
	// nhưng nó có thể được sử dụng để cung cấp dữ liệu bổ sung không được mã hóa
	// nhưng được xác thực (ví dụ: ID của người dùng hoặc phiên).
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	//convert ciphertext to base64 string
	ciphertextStr := base64.StdEncoding.EncodeToString(ciphertext)
	fx, fy := p.decryptBytes(p.EncryptionKey, ciphertextStr)
	fmt.Println(fx, fy)
	n := len(ciphertext)
	fmt.Println(n)
	n = len(ciphertextStr)
	fmt.Println(n)
	return &ciphertextStr, nil

}

// decryptBytes giải mã một slice bytes (data) đã được mã hóa bằng AES-GCM.
// keyString là khóa mã hóa (phải có độ dài 16, 24 hoặc 32 byte khi được chuyển đổi sang []byte).
// data là dữ liệu đã mã hóa (ciphertext, bao gồm nonce ở đầu).
// Hàm này trả về dữ liệu gốc (plaintext).
func (p *TokenService) decryptBytes(keyString string, dataBase64 string) ([]byte, error) {
	// convert base64 string to bytes
	data, err := base64.StdEncoding.DecodeString(dataBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %w", err)
	}
	// Chuyển đổi khóa từ string sang []byte

	key := []byte(keyString)

	// Kiểm tra độ dài khóa
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key length for AES: %d bytes (must be 16, 24, or 32)", len(key))
	}

	// Tạo một block cipher từ khóa
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	// Tạo một Galois Counter Mode (GCM) instance từ block cipher.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}

	// Tách Nonce (IV) từ dữ liệu đã mã hóa.
	// Nhớ rằng gcm.Seal đã nối nonce vào đầu ciphertext.
	nonceSize := gcm.NonceSize() // Lấy kích thước nonce tiêu chuẩn (thường là 12 bytes)
	if len(data) < nonceSize {
		return nil, fmt.Errorf("encrypted data is too short to contain nonce")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Giải mã dữ liệu.
	// gcm.Open sẽ xác minh authentication tag và trả về lỗi nếu dữ liệu bị giả mạo
	// hoặc khóa không đúng.
	// Tham số 'additionalData' (tham số thứ tư) phải khớp với tham số đã dùng khi mã hóa.
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// Lỗi ở đây có thể do:
		// - Khóa mã hóa không khớp.
		// - Dữ liệu đã bị giả mạo.
		// - Nonce bị thay đổi (nếu nó được thay đổi độc lập với ciphertext).
		return nil, fmt.Errorf("failed to decrypt data (possibly invalid key or corrupted data): %w", err)
	}

	return plaintext, nil
}
func (p *TokenService) getJwtSecret() ([]byte, error) {
	if p.EncryptionKey == "" {
		pkgPath := reflect.TypeOf(*p).PkgPath() + "/getJwtSecret"
		panic(fmt.Sprintf("encryption key is missing in %s", pkgPath))
	}

	jwtSecret, err := generateRandomSecret(255)
	if err != nil {
		return nil, err
	}
	encryptedJwtSecret, err := p.encryptBytes(p.EncryptionKey, []byte(jwtSecret))

	if err != nil {
		return nil, err
	}
	err = p.TenantDb.Insert(&authModels.AppConfig{
		Name:      p.TenantDb.TenantDbName,
		Tenant:    p.TenantDb.TenantDbName,
		AppId:     dbx.NewUUID(),
		JwtSecret: *encryptedJwtSecret,
	})
	if err != nil {
		if dbxErr, ok := err.(*dbx.DBXError); ok {
			if dbxErr.Code == dbx.DBXErrorCodeDuplicate {
				if dbxErr.Fields[0] == "Tenant" || dbxErr.Fields[0] == "Name" {
					qr := dbx.Query[authModels.AppConfig](p.TenantDb, p.Context).Where("Tenant =?", p.TenantDb.TenantDbName)
					qr.Select("JwtSecret")
					appConfig, err := qr.First()
					if err != nil {
						return nil, fmt.Errorf("failed to get jwt secret: %w", err)
					}
					if appConfig.JwtSecret == "" {
						return nil, fmt.Errorf("failed to get jwt secret: %w", err)
					}
					// Decrypt jwt secret

					return p.decryptBytes(p.EncryptionKey, appConfig.JwtSecret)

				}
			}
		}
		return nil, err
	}

	return []byte(jwtSecret), nil
}
