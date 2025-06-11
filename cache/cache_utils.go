package caching

import (
	"bytes"
	"encoding/gob"
)

func bytesEncodeObject(obj interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer) // Tạo Encoder
	err := encoder.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func bytesDecodeObject(objBytes []byte, obj interface{}) error {
	decoder := gob.NewDecoder(bytes.NewBuffer(objBytes)) // Tạo Decoder từ []byte
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}
