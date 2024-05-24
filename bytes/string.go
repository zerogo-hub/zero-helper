package bytes

import (
	"strings"
	"unsafe"
)

// StringToByte 字符串转 []byte，避免 []byte(str) 带来的数据复制，转出的数据不可写
func StringToBytes(s string) []byte {
	// return unsafe.Slice(unsafe.StringData(s), len(s))
	return *(*[]byte)(unsafe.Pointer(&s))
}

// SliceByteToString []byte 转字符串，避免 string([]byte) 带来的数据复制
func SliceByteToString(b []byte) string {
	// return unsafe.String(unsafe.SliceData(b), len(b))
	return *(*string)(unsafe.Pointer(&b))
}

// CharUpper 转大写 a -> A
func CharUpper(char byte) byte {
	if char >= 'a' && char <= 'z' {
		return char - 32
	}

	return char
}

// CharLower 转小写 A -> a
func CharLower(char byte) byte {
	if char >= 'A' && char <= 'Z' {
		return char + 32
	}

	return char
}

// FirstUpper 首字母大写
func FirstUpper(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 首字母小写
func FirstLower(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToLower(s[:1]) + s[1:]
}
