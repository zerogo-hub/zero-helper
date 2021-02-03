package bytes

import (
	"reflect"
	"unsafe"
)

// StringToByte 字符串转 []byte，转出的 []byte 可读不可写
//
// 来源: https://www.toutiao.com/a6918883127146349067/
func StringToByte(s string) []byte {
	l := len(s)
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (*(*reflect.StringHeader)(unsafe.Pointer(&s))).Data,
		Len:  l,
		Cap:  l,
	}))
}
