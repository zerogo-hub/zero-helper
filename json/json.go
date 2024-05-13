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

// UnmarshalNumber 结构体中包含大整型，如 int64、uint64，使用默认的 Unmarshal 可能会丢失精度
// Deprecated: Use Unmarshal directly.
func UnmarshalNumber(data []byte, v interface{}) error {
	return Unmarshal(data, v)
}
