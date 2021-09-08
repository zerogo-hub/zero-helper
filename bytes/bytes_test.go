package bytes_test

import (
	"testing"

	zerobytes "github.com/zerogo-hub/zero-helper/bytes"
)

func TestInt(t *testing.T) {
	if int8(55) != zerobytes.ToInt8(zerobytes.PutInt8(int8(55))) {
		t.Error("int8 error")
	}

	if uint8(55) != zerobytes.ToUint8(zerobytes.PutUint8(uint8(55))) {
		t.Error("uint8 error")
	}

	if int16(55) != zerobytes.ToInt16(zerobytes.PutInt16(int16(55))) {
		t.Error("int16 error")
	}

	if uint16(55) != zerobytes.ToUint16(zerobytes.PutUint16(uint16(55))) {
		t.Error("uint16 error")
	}

	if int32(55) != zerobytes.ToInt32(zerobytes.PutInt32(int32(55))) {
		t.Error("int32 error")
	}

	if uint32(55) != zerobytes.ToUint32(zerobytes.PutUint32(uint32(55))) {
		t.Error("uint32 error")
	}

	if int64(55) != zerobytes.ToInt64(zerobytes.PutInt64(int64(55))) {
		t.Error("int64 error")
	}

	if uint64(55) != zerobytes.ToUint64(zerobytes.PutUint64(uint64(55))) {
		t.Error("uint64 error")
	}
}

func TestStringToBytes(t *testing.T) {
	s := "hello world"

	b1 := zerobytes.StringToBytes(s)
	b2 := []byte(s)

	if len(b1) != len(b2) {
		t.Error("TestStringToBytes error")
		return
	}

	for idx, b10 := range b1 {
		b20 := b2[idx]
		if b10 != b20 {
			t.Error("TestStringToBytes not the same")
			break
		}
	}
}

func TestCharLower(t *testing.T) {
	if zerobytes.CharLower('A') != 'a' {
		t.Error("CharLower A error")
	}

	if zerobytes.CharLower('Z') != 'z' {
		t.Error("CharLower Z error")
	}

	if zerobytes.CharLower(1) != 1 {
		t.Error("CharLower 1 error")
	}
}

func TestCharUpper(t *testing.T) {
	if zerobytes.CharUpper('a') != 'A' {
		t.Error("CharLower a error")
	}

	if zerobytes.CharUpper('z') != 'Z' {
		t.Error("CharLower z error")
	}

	if zerobytes.CharUpper(1) != 1 {
		t.Error("CharLower 1 error")
	}
}

func TestFirstLower(t *testing.T) {
	if zerobytes.FirstLower("Hello world") != "hello world" {
		t.Error("FirstLower error")
	}

	if zerobytes.FirstLower("") != "" {
		t.Error("FirstLower empty error")
	}
}

func TestFirstUpper(t *testing.T) {
	if zerobytes.FirstUpper("hello world") != "Hello world" {
		t.Error("FirstUpper error")
	}

	if zerobytes.FirstUpper("") != "" {
		t.Error("FirstUpper empty error")
	}
}

func TestSliceToInt8(t *testing.T) {
	from := []string{"1", "2", "3"}
	to := []int8{1, 2, 3}

	result, err := zerobytes.SliceStringToInt8(from)
	if err != nil {
		t.Errorf("test SliceStringToInt8 failed, err: %s", err.Error())
		return
	}

	if len(to) != len(result) {
		t.Errorf("test SliceStringToInt8 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToInt8 failed, i: %d, v: %d, r: %d, err: %s", i, v, result[i], err.Error())
		}
	}
}

func TestSliceToUint8(t *testing.T) {
	from := []string{"1", "2", "3"}
	to := []uint8{1, 2, 3}

	result, err := zerobytes.SliceStringToUint8(from)
	if err != nil {
		t.Errorf("test SliceStringToInt8 failed, err: %s", err.Error())
		return
	}

	if len(to) != len(result) {
		t.Errorf("test TestSliceToUint8 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToInt8 failed, i: %d, v: %d, r: %d, err: %s", i, v, result[i], err.Error())
		}
	}
}

func TestSliceToInt64(t *testing.T) {
	from := []string{"80068728175202304"}
	to := []int64{80068728175202304}

	result, err := zerobytes.SliceStringToInt64(from)
	if err != nil {
		t.Errorf("test TestSliceToInt64 failed, err: %s", err.Error())
		return
	}

	if len(to) != len(result) {
		t.Errorf("test TestSliceToInt64 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test TestSliceToInt64 failed, i: %d, v: %d, r: %d, err: %s", i, v, result[i], err.Error())
		}
	}
}

func TestSliceToUint64(t *testing.T) {
	from := []string{"80068728175202304"}
	to := []uint64{80068728175202304}

	result, err := zerobytes.SliceStringToUint64(from)
	if err != nil {
		t.Errorf("test TestSliceToUint64 failed, err: %s", err.Error())
		return
	}

	if len(to) != len(result) {
		t.Errorf("test TestSliceToUint64 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test TestSliceToUint64 failed, i: %d, v: %d, r: %d, err: %s", i, v, result[i], err.Error())
		}
	}
}
