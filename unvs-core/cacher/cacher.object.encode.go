package cacher

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"sync"
)

type utils struct {
	cacheGetHash sync.Map
}

func (u *utils) bytesEncodeObject(obj interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer) // Tạo Encoder
	err := encoder.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func (u *utils) bytesDecodeObject(objBytes []byte, obj interface{}) error {
	decoder := gob.NewDecoder(bytes.NewBuffer(objBytes)) // Tạo Decoder từ []byte
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}

// getHashedKey tạo một hash SHA256 từ key đầu vào.
// Memcached có giới hạn về độ dài key (250 ký tự), việc hash giúp xử lý các key dài.
func (u *utils) getHashedKey(key string) string {
	if cachedHash, ok := u.cacheGetHash.Load(key); ok {
		return cachedHash.(string)
	}
	hasher := sha256.New()
	hasher.Write([]byte(key))
	ret := hex.EncodeToString(hasher.Sum(nil))
	u.cacheGetHash.Store(key, ret)
	return ret
}

var Utils = &utils{}
