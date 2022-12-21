package bytes

import (
	"reflect"
	"strings"
	"unsafe"
)

// StringToByte 字符串转 []byte，转出的 []byte 可读不可写
//
// 来源: kitex/pkg/utils/byte2str.go
//
// 来源: https://www.toutiao.com/a6918883127146349067/
func StringToBytes(s string) []byte {
	p := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data)
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Data = uintptr(p)
	hdr.Cap = len(s)
	hdr.Len = len(s)
	return b
}

// SliceByteToString ..
// 来源: kitex/pkg/utils/byte2str.go
func SliceByteToString(b []byte) string {
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
