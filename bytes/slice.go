package bytes

import (
	"errors"
	"reflect"
	"unsafe"
)

func sliceConvert(from interface{}, to reflect.Type) (interface{}, error) {
	sv := reflect.ValueOf(from)
	if sv.Kind() != reflect.Slice {
		return errors.New("from non-slice"), nil
	}
	if to.Kind() != reflect.Slice {
		return errors.New("to non-slice"), nil
	}
	newSlice := reflect.New(to)
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(newSlice.Pointer()))
	hdr.Cap = sv.Cap() * int(sv.Type().Elem().Size()) / int(to.Elem().Size())
	hdr.Len = sv.Len() * int(sv.Type().Elem().Size()) / int(to.Elem().Size())
	hdr.Data = uintptr(sv.Pointer())
	return newSlice.Elem().Interface(), nil
}

// SliceStringToInt8 []string 转 []int
func SliceStringToInt8(from []string) ([]int8, error) {
	var to []int8
	result, err := sliceConvert(from, reflect.TypeOf(to))
	if err != nil {
		return nil, err
	}

	return result.([]int8), nil
}

// SliceStringToUint8 []string 转 []uint8
func SliceStringToUint8(from []string) ([]uint8, error) {
	var to []uint8
	result, err := sliceConvert(from, reflect.TypeOf(to))
	if err != nil {
		return nil, err
	}

	return result.([]uint8), nil
}

// SliceStringToInt16 []string 转 []int
func SliceStringToInt16(from []string) ([]int16, error) {
	var to []int16
	result, err := sliceConvert(from, reflect.TypeOf(to))
	if err != nil {
		return nil, err
	}

	return result.([]int16), nil
}

// SliceStringToUint16 []string 转 []uint16
func SliceStringToUint16(from []string) ([]uint16, error) {
	var to []uint16
	result, err := sliceConvert(from, reflect.TypeOf(to))
	if err != nil {
		return nil, err
	}

	return result.([]uint16), nil
}

// SliceStringToInt32 []string 转 []int
func SliceStringToInt32(from []string) ([]int32, error) {
	var to []int32
	result, err := sliceConvert(from, reflect.TypeOf(to))
	if err != nil {
		return nil, err
	}

	return result.([]int32), nil
}

// SliceStringToUint32 []string 转 []uint32
func SliceStringToUint32(from []string) ([]uint32, error) {
	var to []uint32
	result, err := sliceConvert(from, reflect.TypeOf(to))
	if err != nil {
		return nil, err
	}

	return result.([]uint32), nil
}

// SliceStringToInt []string 转 []int
func SliceStringToInt(from []string) ([]int, error) {
	var to []int
	result, err := sliceConvert(from, reflect.TypeOf(to))
	if err != nil {
		return nil, err
	}

	return result.([]int), nil
}

// SliceStringToInt []string 转 []uint
func SliceStringToUint(from []string) ([]uint, error) {
	var to []uint
	result, err := sliceConvert(from, reflect.TypeOf(to))
	if err != nil {
		return nil, err
	}

	return result.([]uint), nil
}

// SliceStringToInt64 []string 转 []int64
func SliceStringToInt64(from []string) ([]int64, error) {
	var to []int64
	result, err := sliceConvert(from, reflect.TypeOf(to))
	if err != nil {
		return nil, err
	}

	return result.([]int64), nil
}

// SliceStringToUint64 []string 转 []uint64
func SliceStringToUint64(from []string) ([]uint64, error) {
	var to []uint64
	result, err := sliceConvert(from, reflect.TypeOf(to))
	if err != nil {
		return nil, err
	}

	return result.([]uint64), nil
}
