package random

import (
	"bytes"
	cr "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"
	"time"
)

var (
	allLetters         = []byte("abcdefghijklmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789")
	lowerLetters       = []byte("abcdefghijklmnopqrstuvwxyz")
	lowerNumberLetters = []byte("abcdefghjkmnpqrstuvwxyz23456789")
	upperLetters       = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	upperNumberLetters = []byte("ABCDEFGHJKMNPQRSTUVWXYZ23456789")
)

// String 获取指定长度的字符串，包括大小写字母和数字
func String(length int) string {
	return rs(allLetters, length)
}

// Lower 获取指定长度的字符串，仅包含小写字母
func Lower(length int) string {
	return rs(lowerLetters, length)
}

// LowerWithNumber 获取指定长度的字符串，包含小写字母和数字
func LowerWithNumber(length int) string {
	return rs(lowerNumberLetters, length)
}

// Upper 获取指定长度的字符串，仅包含大写字母
func Upper(length int) string {
	return rs(upperLetters, length)
}

// UpperWithNumber 获取指定长度的字符串，包含大写字母和数字
func UpperWithNumber(length int) string {
	return rs(upperNumberLetters, length)
}

func rs(letters []byte, length int) string {
	buf := buffer()
	defer releaseBuffer(buf)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for start := 0; start < length; start++ {
		buf.WriteByte(letters[r.Intn(len(letters))])
	}

	return buf.String()
}

// Int 获取指定范围内的整数
func Int(min, max int64) int64 {
	if min >= max || min == max {
		return max
	}
	return rand.Int63n(max-min) + min
}

// Uint32 获取随机数，类型为 uint32
func Uint32() uint32 {
	var v uint32
	if err := binary.Read(cr.Reader, binary.BigEndian, &v); err == nil {
		return v
	}
	panic("Random failed")
}

var bufferPool *sync.Pool

func init() {
	rand.Seed(time.Now().UnixNano())

	bufferPool = &sync.Pool{}
	bufferPool.New = func() interface{} {
		return &bytes.Buffer{}
	}
}

// buffer 从池中获取 buffer
func buffer() *bytes.Buffer {
	buff := bufferPool.Get().(*bytes.Buffer)
	buff.Reset()
	return buff
}

// releaseBuffer 将 buff 放入池中
func releaseBuffer(buff *bytes.Buffer) {
	bufferPool.Put(buff)
}
