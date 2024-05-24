package bytes_test

import (
	"testing"

	zerobytes "github.com/zerogo-hub/zero-helper/bytes"
)

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
			t.Errorf("test SliceStringToInt8 failed, i: %d, v: %d, r: %d", i, v, result[i])
		}
	}
}

func TestSliceToInt8Invalid(t *testing.T) {
	if result, err := zerobytes.SliceStringToInt8(nil); len(result) > 0 || err != nil {
		t.Errorf("test SliceStringToInt8 failed when empty, err: %s", err.Error())
	}

	if _, err := zerobytes.SliceStringToInt8([]string{"a1"}); err != nil {
		t.Errorf("test SliceStringToInt8 failed when invalid string, err: %s", err.Error())
	}
}

func TestSliceToUint8(t *testing.T) {
	from := []string{"1", "2", "3"}
	to := []uint8{1, 2, 3}

	result, err := zerobytes.SliceStringToUint8(from)
	if err != nil {
		t.Errorf("test TestSliceToUint8 failed, err: %s", err.Error())
		return
	}

	if len(to) != len(result) {
		t.Errorf("test TestSliceToUint8 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test TestSliceToUint8 failed, i: %d, v: %d, r: %d", i, v, result[i])
		}
	}
}

func TestSliceToUint8Invalid(t *testing.T) {
	if result, err := zerobytes.SliceStringToUint8(nil); len(result) > 0 || err != nil {
		t.Errorf("test SliceStringToUint8 failed when empty, err: %s", err.Error())
	}

	if _, err := zerobytes.SliceStringToUint8([]string{"a1"}); err != nil {
		t.Errorf("test SliceStringToUint8 failed when invalid string, err: %s", err.Error())
	}
}

func TestSliceToInt16(t *testing.T) {
	from := []string{"29999", "100", "50"}
	to := []int16{29999, 100, 50}

	result, err := zerobytes.SliceStringToInt16(from)
	if err != nil {
		t.Errorf("test SliceStringToInt16 failed, err: %s", err.Error())
		return
	}

	if len(to) != len(result) {
		t.Errorf("test SliceStringToInt16 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToInt16 failed, i: %d, v: %d, r: %d", i, v, result[i])
		}
	}
}

func TestSliceToInt16Invalid(t *testing.T) {
	if result, err := zerobytes.SliceStringToInt16(nil); len(result) > 0 || err != nil {
		t.Errorf("test SliceStringToInt16 failed when empty, err: %s", err.Error())
	}

	if _, err := zerobytes.SliceStringToInt16([]string{"a1"}); err != nil {
		t.Errorf("test SliceStringToInt16 failed when invalid string, err: %s", err.Error())
	}
	if _, err := zerobytes.SliceStringToInt16([]string{"99999"}); err != nil {
		t.Errorf("test SliceStringToInt16 failed when invalid string, err: %s", err.Error())
	}
}

func TestSliceToUint16(t *testing.T) {
	from := []string{"29999", "100", "50"}
	to := []uint16{29999, 100, 50}

	result, err := zerobytes.SliceStringToUint16(from)
	if err != nil {
		t.Errorf("test SliceStringToUint16 failed, err: %s", err.Error())
		return
	}

	if len(to) != len(result) {
		t.Errorf("test SliceStringToUint16 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToUint16 failed, i: %d, v: %d, r: %d", i, v, result[i])
		}
	}
}

func TestSliceToUint16Invalid(t *testing.T) {
	if result, err := zerobytes.SliceStringToUint16(nil); len(result) > 0 || err != nil {
		t.Errorf("test SliceStringToUint16 failed when empty, err: %s", err.Error())
	}

	if _, err := zerobytes.SliceStringToUint16([]string{"a1"}); err != nil {
		t.Errorf("test SliceStringToUint16 failed when invalid string, err: %s", err.Error())
	}
	if _, err := zerobytes.SliceStringToUint16([]string{"99999"}); err != nil {
		t.Errorf("test SliceStringToUint16 failed when invalid string, err: %s", err.Error())
	}
	if _, err := zerobytes.SliceStringToUint16([]string{"-1"}); err != nil {
		t.Errorf("test SliceStringToUint16 failed when invalid string, err: %s", err.Error())
	}
}

func TestSliceToInt32(t *testing.T) {
	from := []string{"29999", "100", "50"}
	to := []int32{29999, 100, 50}

	result, err := zerobytes.SliceStringToInt32(from)
	if err != nil {
		t.Errorf("test SliceStringToInt32 failed, err: %s", err.Error())
		return
	}

	if len(to) != len(result) {
		t.Errorf("test SliceStringToInt32 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToInt32 failed, i: %d, v: %d, r: %d", i, v, result[i])
		}
	}
}

func TestSliceToInt32Invalid(t *testing.T) {
	if result, err := zerobytes.SliceStringToInt32(nil); len(result) > 0 || err != nil {
		t.Errorf("test SliceStringToInt32 failed when empty, err: %s", err.Error())
	}

	if _, err := zerobytes.SliceStringToInt32([]string{"a1"}); err != nil {
		t.Errorf("test SliceStringToInt32 failed when invalid string, err: %s", err.Error())
	}
	if _, err := zerobytes.SliceStringToInt32([]string{"80068728175202304"}); err != nil {
		t.Errorf("test SliceStringToInt32 failed when invalid string, err: %s", err.Error())
	}
}

func TestSliceToUint32(t *testing.T) {
	from := []string{"29999", "100", "50"}
	to := []uint32{29999, 100, 50}

	result, err := zerobytes.SliceStringToUint32(from)
	if err != nil {
		t.Errorf("test SliceStringToUint32 failed, err: %s", err.Error())
		return
	}

	if len(to) != len(result) {
		t.Errorf("test SliceStringToUint32 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToUint32 failed, i: %d, v: %d, r: %d", i, v, result[i])
		}
	}
}

func TestSliceToUint32Invalid(t *testing.T) {
	if result, err := zerobytes.SliceStringToUint32(nil); len(result) > 0 || err != nil {
		t.Errorf("test SliceStringToUint32 failed when empty, err: %s", err.Error())
	}

	if _, err := zerobytes.SliceStringToUint32([]string{"a1"}); err != nil {
		t.Errorf("test SliceStringToUint32 failed when invalid string, err: %s", err.Error())
	}
	if _, err := zerobytes.SliceStringToUint32([]string{"80068728175202304"}); err != nil {
		t.Errorf("test SliceStringToUint32 failed when invalid string, err: %s", err.Error())
	}
	if _, err := zerobytes.SliceStringToUint32([]string{"-1"}); err != nil {
		t.Errorf("test SliceStringToUint32 failed when invalid string, err: %s", err.Error())
	}
}

func TestSliceToInt(t *testing.T) {
	from := []string{"29999", "100", "50"}
	to := []int{29999, 100, 50}

	result, err := zerobytes.SliceStringToInt(from)
	if err != nil {
		t.Errorf("test SliceStringToInt failed, err: %s", err.Error())
		return
	}

	if len(to) != len(result) {
		t.Errorf("test SliceStringToInt failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToInt failed, i: %d, v: %d, r: %d", i, v, result[i])
		}
	}
}

func TestSliceToIntInvalid(t *testing.T) {
	if result, err := zerobytes.SliceStringToInt(nil); len(result) > 0 || err != nil {
		t.Errorf("test SliceStringToInt failed when empty, err: %s", err.Error())
	}

	if _, err := zerobytes.SliceStringToInt([]string{"a1"}); err != nil {
		t.Errorf("test SliceStringToInt failed when invalid string, err: %s", err.Error())
	}
	if _, err := zerobytes.SliceStringToInt([]string{"80068728175202304"}); err != nil {
		t.Errorf("test SliceStringToInt failed when invalid string, err: %s", err.Error())
	}
}

func TestSliceToUint(t *testing.T) {
	from := []string{"29999", "100", "50"}
	to := []uint{29999, 100, 50}

	result, err := zerobytes.SliceStringToUint(from)
	if err != nil {
		t.Errorf("test SliceStringToUint failed, err: %s", err.Error())
		return
	}

	if len(to) != len(result) {
		t.Errorf("test SliceStringToUint failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToUint failed, i: %d, v: %d, r: %d", i, v, result[i])
		}
	}
}

func TestSliceToUintInvalid(t *testing.T) {
	if result, err := zerobytes.SliceStringToUint(nil); len(result) > 0 || err != nil {
		t.Errorf("test SliceStringToUint failed when empty, err: %s", err.Error())
	}

	if _, err := zerobytes.SliceStringToUint([]string{"a1"}); err != nil {
		t.Errorf("test SliceStringToUint failed when invalid string, err: %s", err.Error())
	}
	if _, err := zerobytes.SliceStringToUint([]string{"80068728175202304"}); err != nil {
		t.Errorf("test SliceStringToUint failed when invalid string, err: %s", err.Error())
	}
	if _, err := zerobytes.SliceStringToUint([]string{"-1"}); err != nil {
		t.Errorf("test SliceStringToUint failed when invalid string, err: %s", err.Error())
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
			t.Errorf("test TestSliceToInt64 failed, i: %d, v: %d, r: %d", i, v, result[i])
		}
	}
}

func TestSliceToInt64Invalid(t *testing.T) {
	if result, err := zerobytes.SliceStringToInt64(nil); len(result) > 0 || err != nil {
		t.Errorf("test SliceStringToInt64 failed when empty, err: %s", err.Error())
	}

	if _, err := zerobytes.SliceStringToInt64([]string{"a1"}); err != nil {
		t.Errorf("test SliceStringToInt64 failed when invalid string, err: %s", err.Error())
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
			t.Errorf("test TestSliceToUint64 failed, i: %d, v: %d, r: %d", i, v, result[i])
		}
	}
}

func TestSliceToUint64Invalid(t *testing.T) {
	if result, err := zerobytes.SliceStringToUint64(nil); len(result) > 0 || err != nil {
		t.Errorf("test SliceStringToUint64 failed when empty, err: %s", err.Error())
	}

	if _, err := zerobytes.SliceStringToUint64([]string{"a1"}); err != nil {
		t.Errorf("test SliceStringToUint64 failed when invalid string, err: %s", err.Error())
	}
}

func TestInt8ToString(t *testing.T) {
	from := []int8{1, 2, 3}
	to := []string{"1", "2", "3"}

	result := zerobytes.SliceInt8ToString(from)
	if len(to) != len(result) {
		t.Errorf("test SliceStringToInt8 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToInt8 failed, i: %d, v: %s, r: %s", i, v, result[i])
		}
	}
}

func TestInt8ToStringInvalid(t *testing.T) {
	if result := zerobytes.SliceInt8ToString(nil); len(result) > 0 {
		t.Errorf("test SliceStringToInt8 failed when empty")
	}
}

func TestUint8ToString(t *testing.T) {
	from := []uint8{1, 2, 3}
	to := []string{"1", "2", "3"}

	result := zerobytes.SliceUint8ToString(from)
	if len(to) != len(result) {
		t.Errorf("test SliceStringToUint8 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToUint8 failed, i: %d, v: %s, r: %s", i, v, result[i])
		}
	}
}

func TestUint8ToStringInvalid(t *testing.T) {
	if result := zerobytes.SliceUint8ToString(nil); len(result) > 0 {
		t.Errorf("test SliceStringToUint8 failed when empty")
	}
}

func TestInt16ToString(t *testing.T) {
	from := []int16{1, 2, 3}
	to := []string{"1", "2", "3"}

	result := zerobytes.SliceInt16ToString(from)
	if len(to) != len(result) {
		t.Errorf("test SliceStringToInt16 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToInt16 failed, i: %d, v: %s, r: %s", i, v, result[i])
		}
	}
}

func TestInt16ToStringInvalid(t *testing.T) {
	if result := zerobytes.SliceInt16ToString(nil); len(result) > 0 {
		t.Errorf("test SliceStringToInt16 failed when empty")
	}
}

func TestUint16ToString(t *testing.T) {
	from := []uint16{1, 2, 3}
	to := []string{"1", "2", "3"}

	result := zerobytes.SliceUint16ToString(from)
	if len(to) != len(result) {
		t.Errorf("test SliceStringToUint16 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToUint16 failed, i: %d, v: %s, r: %s", i, v, result[i])
		}
	}
}

func TestUint16ToStringInvalid(t *testing.T) {
	if result := zerobytes.SliceUint16ToString(nil); len(result) > 0 {
		t.Errorf("test SliceStringToUint16 failed when empty")
	}
}

func TestInt32ToString(t *testing.T) {
	from := []int32{1, 2, 3}
	to := []string{"1", "2", "3"}

	result := zerobytes.SliceInt32ToString(from)
	if len(to) != len(result) {
		t.Errorf("test SliceStringToInt32 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToInt32 failed, i: %d, v: %s, r: %s", i, v, result[i])
		}
	}
}

func TestInt32ToStringInvalid(t *testing.T) {
	if result := zerobytes.SliceInt32ToString(nil); len(result) > 0 {
		t.Errorf("test SliceStringToInt32 failed when empty")
	}
}

func TestUint32ToString(t *testing.T) {
	from := []uint32{1, 2, 3}
	to := []string{"1", "2", "3"}

	result := zerobytes.SliceUint32ToString(from)
	if len(to) != len(result) {
		t.Errorf("test SliceStringToUint32 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToUint32 failed, i: %d, v: %s, r: %s", i, v, result[i])
		}
	}
}

func TestUint32ToStringInvalid(t *testing.T) {
	if result := zerobytes.SliceUint32ToString(nil); len(result) > 0 {
		t.Errorf("test SliceStringToUint32 failed when empty")
	}
}

func TestIntToString(t *testing.T) {
	from := []int{1, 2, 3}
	to := []string{"1", "2", "3"}

	result := zerobytes.SliceIntToString(from)
	if len(to) != len(result) {
		t.Errorf("test SliceStringToInt failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToInt failed, i: %d, v: %s, r: %s", i, v, result[i])
		}
	}
}

func TestIntToStringInvalid(t *testing.T) {
	if result := zerobytes.SliceIntToString(nil); len(result) > 0 {
		t.Errorf("test SliceStringToInt failed when empty")
	}
}

func TestUintToString(t *testing.T) {
	from := []uint{1, 2, 3}
	to := []string{"1", "2", "3"}

	result := zerobytes.SliceUintToString(from)
	if len(to) != len(result) {
		t.Errorf("test SliceStringToUint failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToUint failed, i: %d, v: %s, r: %s", i, v, result[i])
		}
	}
}

func TestUintToStringInvalid(t *testing.T) {
	if result := zerobytes.SliceUintToString(nil); len(result) > 0 {
		t.Errorf("test SliceStringToUint failed when empty")
	}
}

func TestInt64ToString(t *testing.T) {
	from := []int64{1, 2, 3}
	to := []string{"1", "2", "3"}

	result := zerobytes.SliceInt64ToString(from)
	if len(to) != len(result) {
		t.Errorf("test SliceStringToInt64 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToInt64 failed, i: %d, v: %s, r: %s", i, v, result[i])
		}
	}
}

func TestInt64ToStringInvalid(t *testing.T) {
	if result := zerobytes.SliceInt64ToString(nil); len(result) > 0 {
		t.Errorf("test SliceStringToInt64 failed when empty")
	}
}

func TestUint64ToString(t *testing.T) {
	from := []uint64{1, 2, 3}
	to := []string{"1", "2", "3"}

	result := zerobytes.SliceUint64ToString(from)
	if len(to) != len(result) {
		t.Errorf("test SliceStringToUint64 failed, invalid length")
		return
	}

	for i, v := range to {
		if v != result[i] {
			t.Errorf("test SliceStringToUint64 failed, i: %d, v: %s, r: %s", i, v, result[i])
		}
	}
}

func TestUint64ToStringInvalid(t *testing.T) {
	if result := zerobytes.SliceUint64ToString(nil); len(result) > 0 {
		t.Errorf("test SliceStringToUint64 failed when empty")
	}
}
