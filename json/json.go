package json

import (
	libJSON "github.com/bytedance/sonic"
)

// Marshal ..
func Marshal(v interface{}) ([]byte, error) {
	return libJSON.Marshal(v)
}

// Unmarshal ...
func Unmarshal(data []byte, v interface{}) error {
	return libJSON.Unmarshal(data, v)
}
