package bytes

import (
	"encoding/binary"
)

// ToInt8 转为 int8
func ToInt8(p []byte) int8 {
	return int8(p[0])
}

// PutInt8 int8 转 []byte
func PutInt8(n int8) []byte {
	bytes := make([]byte, 1)
	bytes[0] = byte(n)
	return bytes
}

// ToUint8 转为 uint8
func ToUint8(p []byte) uint8 {
	return uint8(p[0])
}

// PutUint8 uint8 转 []byte
func PutUint8(n uint8) []byte {
	bytes := make([]byte, 1)
	bytes[0] = byte(n)
	return bytes
}

// ToInt16 转 int16
func ToInt16(p []byte) int16 {
	return int16(binary.BigEndian.Uint16(p))
}

// PutInt16 int16 转 []byte
func PutInt16(n int16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, uint16(n))
	return bytes
}

// ToUint16 转 uint16
func ToUint16(p []byte) uint16 {
	return binary.BigEndian.Uint16(p)
}

// PutUint16 uint16 转 []byte
func PutUint16(n uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, n)
	return bytes
}

// ToInt32 转 int32
func ToInt32(p []byte) int32 {
	return int32(binary.BigEndian.Uint32(p))
}

// PutInt32 int32 转 []byte
func PutInt32(n int32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(n))
	return bytes
}

// ToUint32 转 uint32
func ToUint32(p []byte) uint32 {
	return binary.BigEndian.Uint32(p)
}

// PutUint32 uint32 转 []byte
func PutUint32(n uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, n)
	return bytes
}

// ToInt64 转 int64
func ToInt64(p []byte) int64 {
	return int64(binary.BigEndian.Uint64(p))
}

// PutInt64 int64 转 []byte
func PutInt64(n int64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(n))
	return bytes
}

// ToUint64 转 uint64
func ToUint64(p []byte) uint64 {
	return binary.BigEndian.Uint64(p)
}

// PutUint64 uint64 转 []byte
func PutUint64(n uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, n)
	return bytes
}
